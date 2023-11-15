use std::collections::HashMap;

use axum::{
    extract::{OriginalUri, Query},
    response::IntoResponse,
};
use chrono::{NaiveDate, Utc};
use hyper::Uri;

use crate::errors::Error;

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
) -> Result<axum::response::Response, Error> {
    if req.uri().query().is_some() && req.uri().query().unwrap().contains("date") {
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

    let mut query = req
        .extensions()
        .get::<Query<HashMap<String, String>>>()
        .map_or(HashMap::new(), |Query(q)| q.clone());

    query.insert("date".to_owned(), date.to_string());

    let query = query
        .iter()
        .map(|(k, v)| format!("{k}={v}"))
        .collect::<Vec<String>>()
        .join("&");

    Ok(axum::response::Redirect::temporary(&format!("{original_uri}?{query}")).into_response())
}
