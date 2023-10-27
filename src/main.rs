use std::collections::HashMap;
use std::sync::Arc;

use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use hyper_tls::HttpsConnector;
use serde_json::Value;
use tera::{Context, Tera};
use tower::ServiceBuilder;
use tower_http::services::ServeDir;

use crate::marvel::{MarvelService, Marvel};
use crate::marvel::auth::AuthMiddlewareLayer;
use crate::marvel::etag::{EtagMiddlewareLayer, new_etag_cache};
use crate::middleware::uri::UriMiddlewareLayer;

mod marvel;
mod middleware;

struct ComicShelf<S> {
    tera: Tera,
    marvel_client: Marvel<S>,
}

impl<S> ComicShelf<S> {
    fn new(tera: Tera, marvel_client: Marvel<S>) -> Self {
        ComicShelf {
            tera,
            marvel_client,
        }
    }
}

async fn marvel_unlimited_comics<S>(
    State(state): State<Arc<ComicShelf<S>>>,
) -> Result<Html<String>, StatusCode>
where
    S: MarvelService
{
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", "/marvel-unlimited/comics");
    ctx.insert("Date", "2023-08-12");

    let result = state.marvel_client.weekly_comics().await;

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
    let svc = ServiceBuilder::new()
        .layer(EtagMiddlewareLayer::new(new_etag_cache()))
        .layer(UriMiddlewareLayer::new("gateway.marvel.com", "https"))
        .layer(AuthMiddlewareLayer::new(
            include_str!("../pub.txt"), // todo: formalize, this is janky
            include_str!("../priv.txt"),
        ))
        .service(client.clone());
    let marvel_client = Marvel::new(svc);

    let state = Arc::new(ComicShelf::new(tera, marvel_client));

    let app = axum::Router::new()
        .route("/marvel-unlimited/comics", get(marvel_unlimited_comics))
        .nest_service("/static", ServeDir::new("static"))
        .with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
