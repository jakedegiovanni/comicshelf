use axum::{
    extract::{OriginalUri, Query},
    response::IntoResponse,
};
use chrono::{NaiveDate, Utc};
use hyper::Uri;

use crate::middleware;

pub const MARVEL_PATH: &str = "/marvel-unlimited";

fn today() -> NaiveDate {
    Utc::now().date_naive()
}

#[derive(serde::Deserialize)]
pub struct Date {
    #[serde(default = "today")]
    pub date: NaiveDate,
}

pub async fn enforce_date_query<B>(
    req: axum::http::Request<B>,
    next: axum::middleware::Next<B>,
) -> Result<axum::response::Response, middleware::Error> {
    if req.uri().query().is_some() {
        return Ok(next.run(req).await);
    }

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

    Ok(axum::response::Redirect::temporary(&format!("{original_uri}?date={date}")).into_response())
}
