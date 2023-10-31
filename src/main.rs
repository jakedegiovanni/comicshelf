use std::collections::HashMap;
use std::sync::Arc;

use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use axum::response::{IntoResponse, Response};
use chrono::Utc;
use hyper_tls::HttpsConnector;
use serde_json::Value;
use tera::{Context, Tera};
use thiserror::Error;
use tower::ServiceBuilder;
use tower_http::services::ServeDir;

use crate::marvel::{Marvel, MarvelService};
use crate::marvel::auth::AuthMiddlewareLayer;
use crate::marvel::etag::{EtagMiddlewareLayer, new_etag_cache};
use crate::middleware::uri::UriMiddlewareLayer;

mod marvel;
mod middleware;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("internal error")]
    Anyhow(#[from] anyhow::Error),
    #[error("rendering error")]
    Tera(#[from] tera::Error)
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("something went wrong: {}", self)
        ).into_response()
    }
}

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
) -> Result<Html<String>, AppError>
where
    S: MarvelService
{
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", "/marvel-unlimited/comics");

    let date = Utc::now();
    ctx.insert("Date", &date.format("%Y-%m-%d").to_string());

    let result = state.marvel_client.weekly_comics(date).await?;
    ctx.insert("results", &result);

    Ok(Html(state.tera.render("marvel-unlimited.html", &ctx)?))
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
    let svc = ServiceBuilder::new()
        .layer(EtagMiddlewareLayer::new(new_etag_cache()))
        .layer(UriMiddlewareLayer::new(
            "gateway.marvel.com",
            hyper::http::uri::Scheme::HTTPS,
            "/v1/public"
        ))
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
