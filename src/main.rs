use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use reqwest::Client;
use serde_json::Value;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::SystemTime;
use tera::{Context, Tera};
use tower_http::services::ServeDir;

struct ComicShelf {
    tera: Tera,
    client: Client,
}

async fn marvel_unlimited_comics(
    State(state): State<Arc<ComicShelf>>,
) -> Result<Html<String>, StatusCode> {
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", "/marvel-unlimited/comics");
    ctx.insert("Date", "2023-08-12");

    let ts = SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .unwrap()
        .as_millis();
    let hash = "";

    let result = match state.client.get(format!("https://gateway.marvel.com/v1/public/comics?format=comic&formatType=comic&noVariants=true&dateRange=2023-01-01,2023-01-07&hasDigitalIssue=true&orderBy=issueNumber&limit=100&apikey=...&ts={}&hash={}", ts, hash)).send().await {
        Ok(r) => r,
        Err(e) => {
            println!("http client error: {}", e);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    }
    .text()
    .await;

    let result: HashMap<String, Value> = match result {
        Ok(r) => serde_json::from_str(&r).unwrap(),
        Err(e) => {
            println!("text errors: {}", e);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };

    ctx.insert("results", &result);

    let body = match state.tera.render("marvel-unlimited.html", &ctx) {
        Ok(b) => b,
        Err(e) => {
            println!("tera rendering error: {}", e);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };

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

    let client = Client::new();
    let state = Arc::new(ComicShelf { tera, client });

    let app = axum::Router::new()
        .route("/marvel-unlimited/comics", get(marvel_unlimited_comics))
        .nest_service("/static", ServeDir::new("static"))
        .with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
