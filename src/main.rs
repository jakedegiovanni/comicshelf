use std::collections::HashMap;
use std::sync::Arc;

use axum::{debug_handler, extract::State, http::StatusCode, response::Html, routing::get};
use hyper_tls::HttpsConnector;
use serde_json::Value;
use tera::{Context, Tera};
use tower_http::services::ServeDir;

use crate::marvel::Marvel;

mod marvel;
mod middleware;

struct ComicShelf {
    tera: Tera,
    marvel_client: Marvel,
}

impl ComicShelf {
    fn new(tera: Tera, marvel_client: Marvel) -> Self {
        ComicShelf {
            tera,
            marvel_client,
        }
    }
}

#[debug_handler]
async fn marvel_unlimited_comics(
    State(state): State<Arc<ComicShelf>>,
) -> Result<Html<String>, StatusCode> {
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
    let marvel_client = Marvel::new(&client);

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
