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
use std::str::FromStr;
use std::sync::Arc;

use anyhow::anyhow;
use axum::extract::{OriginalUri, Query};
use axum::http::uri::InvalidUri;
use axum::response::{IntoResponse, Response};
use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use chrono::{NaiveDate, ParseError, Utc};
use hyper::Uri;
use hyper_tls::HttpsConnector;
use serde_json::Value;
use tera::{Context, Tera};
use thiserror::Error;
use tower::{BoxError, ServiceBuilder};
use tower_http::services::ServeDir;

use crate::marvel::{auth, etag, Marvel};
use crate::middleware::uri;

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
    #[error("parse error")]
    ParseError(#[from] ParseError),
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

fn today() -> NaiveDate {
    Utc::now().date_naive()
}

#[derive(serde::Deserialize)]
struct Date {
    #[serde(default = "today")]
    date: NaiveDate,
}

async fn marvel_unlimited_comics<S>(
    State(state): State<Arc<ComicShelf<S>>>,
    Query(query): Query<Date>,
    OriginalUri(original_uri): OriginalUri,
) -> Result<Html<String>, AppError>
where
    S: marvel::Client,
{
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", original_uri.path());

    ctx.insert("Date", &query.date.to_string());

    let result = state.marvel_client.weekly_comics(query.date).await?;
    ctx.insert("results", &result);

    Ok(Html(state.tera.render("marvel-unlimited.html", &ctx)?))
}

fn following(args: &HashMap<String, Value>) -> tera::Result<Value> {
    let _ = args.get("index").ok_or(tera::Error::msg("not found"))?; // todo use for db check
    Ok(tera::to_value(false)?)
}

async fn enforce_date_query<B>(
    req: axum::http::Request<B>,
    next: axum::middleware::Next<B>,
) -> Result<axum::response::Response, AppError> {
    if req.uri().query().is_none() {
        let default_uri = OriginalUri(Uri::from_static("/"));
        let original_uri = req
            .extensions()
            .get::<OriginalUri>()
            .unwrap_or(&default_uri)
            .path();

        let date = req
            .extensions()
            .get::<Query<Date>>()
            .unwrap_or(&Query(Date { date: today() }))
            .date;

        return Ok(
            axum::response::Redirect::temporary(&format!("{original_uri}?date={date}"))
                .into_response(),
        );
    }

    Ok(next.run(req).await)
}

#[tokio::main]
async fn main() {
    let mut tera = Tera::new("templates/**/*.html").unwrap();
    tera.register_function("following", following);

    let https = HttpsConnector::new();
    let client = hyper::Client::builder().build::<_, hyper::Body>(https);
    let svc = ServiceBuilder::new()
        .layer(etag::CacheMiddlewareLayer::new(etag::new_etag_cache()))
        .layer(uri::MiddlewareLayer::new(
            "gateway.marvel.com",
            hyper::http::uri::Scheme::HTTPS,
            "/v1/public",
        ))
        .layer(auth::MiddlewareLayer::new(
            include_str!("../pub.txt"), // todo: formalize, this is janky
            include_str!("../priv.txt"),
        ))
        .service(client.clone());
    let marvel_client = Marvel::new(svc);

    let state = Arc::new(ComicShelf::new(tera, marvel_client));

    let app = axum::Router::new()
        .nest(
            "/marvel-unlimited",
            axum::Router::new()
                .route("/comics", get(marvel_unlimited_comics))
                .layer(ServiceBuilder::new().layer(axum::middleware::from_fn(enforce_date_query))),
        )
        .nest_service("/static", ServeDir::new("static"))
        .with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
