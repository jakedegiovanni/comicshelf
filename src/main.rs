use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use chrono::{Datelike, Days, Months, Utc, Weekday};
use hyper::client::HttpConnector;
use hyper::Request;
use hyper_tls::HttpsConnector;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;
use std::sync::Arc;
use tera::{Context, Tera};
use tower::{Service, ServiceBuilder};
use tower_http::services::ServeDir;

#[derive(Serialize, Deserialize, Debug)]
struct DataContainer {
    offset: i32,
    limit: i32,
    total: i32,
    count: i32,
    results: Value,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize, Debug)]
struct DataWrapper {
    code: Value,
    status: String,
    copyright: String,
    attributionText: String,
    attributionHTML: String,
    etag: String,
    data: DataContainer,
}

type HyperService = dyn Service<
        Request<hyper::Body>,
        Error = hyper::Error,
        Future = hyper::client::ResponseFuture,
        Response = hyper::Response<hyper::Body>,
    > + Send
    + Sync;

struct ComicShelf {
    tera: Tera,
    client: hyper::Client<HttpsConnector<HttpConnector>, hyper::Body>,
}

impl ComicShelf {
    fn new(tera: Tera, client: hyper::Client<HttpsConnector<HttpConnector>, hyper::Body>) -> Self {
        ComicShelf { tera, client }
    }

    fn marvel_client(&self) -> Marvel {
        Marvel::new(&self.client)
    }
}

struct Marvel {
    svc: Box<HyperService>,
}

impl Marvel {
    fn new(client: &hyper::Client<HttpsConnector<HttpConnector>, hyper::Body>) -> Self {
        let svc = ServiceBuilder::new()
            .map_request(Marvel::uri_middleware)
            .map_request(Marvel::auth_middleware)
            .service(client.clone());

        let svc = Box::new(svc);

        Marvel { svc }
    }

    fn uri_middleware(req: Request<hyper::Body>) -> Request<hyper::Body> {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        up.authority = Some(hyper::http::uri::Authority::from_static(
            "gateway.marvel.com",
        ));
        up.scheme = Some(hyper::http::uri::Scheme::HTTPS);

        p.uri = hyper::Uri::from_parts(up).unwrap();
        Request::from_parts(p, b)
    }

    fn auth_middleware(req: Request<hyper::Body>) -> Request<hyper::Body> {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        let pq = up.path_and_query.unwrap();
        let path = pq.path();
        let q = pq.query().unwrap_or("");

        let path = {
            if !path.contains("/v1/public") {
                format!("/v1/public{}", path)
            } else {
                path.to_string()
            }
        };

        let pub_key = include_str!("../pub.txt"); // todo: formalize, this is janky
        let priv_key = include_str!("../priv.txt");

        let ts = Utc::now().timestamp_millis();
        let hash = format!(
            "{:x}",
            md5::compute(format!("{}{}{}", ts, priv_key, pub_key))
        );
        let query = format!("apikey={}&ts={}&hash={}", pub_key, ts, hash);
        let query = {
            if q.is_empty() {
                format!("?{}", query)
            } else {
                format!("{}&{}", q, query)
            }
        };
        let query = format!("{}?{}", path, query);
        up.path_and_query = Some(hyper::http::uri::PathAndQuery::try_from(query).unwrap());

        p.uri = hyper::Uri::from_parts(up).unwrap();
        Request::from_parts(p, b)
    }

    async fn weekly_comics(self) -> DataWrapper {
        let mut svc = self.svc;

        let date = Utc::now();
        let mut date = date.checked_sub_months(Months::new(3)).unwrap();
        loop {
            if date.weekday() == Weekday::Sun {
                break;
            }

            date = date.checked_sub_days(Days::new(1)).unwrap();
        }

        date = date.checked_sub_days(Days::new(7)).unwrap();
        let date2 = date.checked_add_days(Days::new(6)).unwrap();

        let date = date.format("%Y-%m-%d").to_string();
        let date2 = date2.format("%Y-%m-%d").to_string();

        let req = Request::get(format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={},{}&hasDigitalIssue=true&orderBy=issueNumber&limit=100", date, date2)).body(hyper::Body::empty()).unwrap();
        let result = svc.call(req).await.unwrap().into_body();
        let result: DataWrapper = serde_json::from_slice(
            hyper::body::to_bytes(result)
                .await
                .unwrap()
                .iter()
                .as_slice(),
        )
        .unwrap();
        result
    }
}

async fn marvel_unlimited_comics(
    State(state): State<Arc<ComicShelf>>,
) -> Result<Html<String>, StatusCode> {
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", "/marvel-unlimited/comics");
    ctx.insert("Date", "2023-08-12");

    let client = state.marvel_client();
    let result = client.weekly_comics().await;

    ctx.insert("results", &result);

    let body = state.tera.render("marvel-unlimited.html", &ctx).unwrap();
    Ok(Html(body))
}

fn following(args: &HashMap<String, Value>) -> tera::Result<Value> {
    let _ = args.get("index").unwrap(); // todo use for db check
    Ok(tera::to_value(false).unwrap())
}

#[tokio::main]
async fn main() {
    let mut tera = match Tera::new("templates/**/*.html") {
        Ok(t) => t,
        Err(e) => {
            println!("Parsing error(s): {}", e);
            std::process::exit(1);
        }
    };
    tera.register_function("following", following);

    let https = HttpsConnector::new();
    let client = hyper::Client::builder().build::<_, hyper::Body>(https);

    let state = Arc::new(ComicShelf::new(tera, client));

    let app = axum::Router::new()
        .route("/marvel-unlimited/comics", get(marvel_unlimited_comics))
        .nest_service("/static", ServeDir::new("static"))
        .with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
