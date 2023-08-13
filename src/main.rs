use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use hyper::Request;
use hyper_tls::HttpsConnector;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::sync::Arc;
use std::time::SystemTime;
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

struct ComicShelf {
    tera: Tera,
}

impl ComicShelf {
    fn new(tera: Tera) -> Self {
        ComicShelf { tera }
    }
}

async fn marvel_unlimited_comics(
    State(state): State<Arc<ComicShelf>>,
) -> Result<Html<String>, StatusCode> {
    let https = HttpsConnector::new();
    let c = hyper::Client::builder().build::<_, hyper::Body>(https);
    let mut svc = ServiceBuilder::new()
        .map_request(|req: Request<hyper::Body>| {
            let (mut p, b) = req.into_parts();

            let mut up = p.uri.into_parts();
            up.authority = Some(hyper::http::uri::Authority::from_static(
                "gateway.marvel.com",
            ));
            up.scheme = Some(hyper::http::uri::Scheme::HTTPS);

            p.uri = hyper::Uri::from_parts(up).unwrap();
            Request::from_parts(p, b)
        })
        .service(c);

    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", "/marvel-unlimited/comics");
    ctx.insert("Date", "2023-08-12");

    let ts = SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .unwrap()
        .as_millis();
    let hash = format!("{:x}", md5::compute(format!("{}privKeypubKey", ts)));

    let req = Request::get(format!("/v1/public/comics?format=comic&formatType=comic&noVariants=true&dateRange=2023-01-01,2023-01-07&hasDigitalIssue=true&orderBy=issueNumber&limit=100&apikey=pubKey&ts={}&hash={}", ts, hash)).body(hyper::Body::empty()).unwrap();
    let result = svc.call(req).await.unwrap();
    let result: DataWrapper = serde_json::from_slice(
        hyper::body::to_bytes(result.into_body())
            .await
            .unwrap()
            .iter()
            .as_slice(),
    )
    .unwrap();

    ctx.insert("results", &result);

    let body = state.tera.render("marvel-unlimited.html", &ctx).unwrap();
    Ok(Html(body))
}

#[tokio::main]
async fn main() {
    let tera = match Tera::new("templates/**/*.html") {
        Ok(t) => t,
        Err(e) => {
            println!("Parsing error(s): {}", e);
            std::process::exit(1);
        }
    };

    let state = Arc::new(ComicShelf::new(tera));

    let app = axum::Router::new()
        .route("/marvel-unlimited/comics", get(marvel_unlimited_comics))
        .nest_service("/static", ServeDir::new("static"))
        .with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
