#![warn(
    clippy::all,
    clippy::correctness,
    clippy::suspicious,
    clippy::style,
    clippy::complexity,
    clippy::perf,
    clippy::pedantic
)]

use std::collections::HashMap;

use axum::http::uri::InvalidUri;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use hyper_tls::HttpsConnector;
use marvel::router::MARVEL_PATH;

use serde_json::Value;
use tera::Tera;
use thiserror::Error;
use tower::BoxError;
use tower_http::services::ServeDir;

mod marvel;
mod middleware;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("internal error")]
    Anyhow(#[from] anyhow::Error),
    #[error("rendering error")]
    Tera(#[from] tera::Error),
    #[error("box error")]
    Box(#[from] BoxError),
    #[error("hyper error")]
    HyperError(#[from] hyper::http::Error),
    #[error("uri error")]
    UriError(#[from] InvalidUri),
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("something went wrong: {self}"),
        )
            .into_response()
    }
}

fn following(args: &HashMap<String, Value>) -> tera::Result<Value> {
    let _ = args.get("index").ok_or(tera::Error::msg("not found"))?; // todo use for db check
    Ok(tera::to_value(false)?)
}

#[tokio::main]
async fn main() {
    let mut tera = Tera::new("templates/**/*.html").unwrap();
    tera.register_function("following", following);

    let https = HttpsConnector::new();
    let client = hyper::Client::builder().build::<_, hyper::Body>(https);

    let app = axum::Router::new()
        .nest(MARVEL_PATH, marvel::router::new(tera, &client))
        .nest_service("/static", ServeDir::new("static"));

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
