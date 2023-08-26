mod marvel;
mod template;

use crate::marvel::Marvel;
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
