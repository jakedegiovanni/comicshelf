use std::sync::Arc;

use axum::{
    extract::{OriginalUri, Query, State},
    response::{Html, IntoResponse},
    routing::get,
};
use chrono::{NaiveDate, Utc};
use hyper::{client::HttpConnector, Client, Uri};
use hyper_tls::HttpsConnector;
use tera::{Context, Tera};
use tower::ServiceBuilder;

use crate::middleware::{self, Error};

use super::{Marvel, WebClient};

pub const MARVEL_PATH: &str = "/marvel-unlimited";

struct MarvelState<S> {
    tera: Tera,
    marvel_client: Marvel<S>,
}

pub fn new(tera: Tera, client: &Client<HttpsConnector<HttpConnector>>) -> axum::Router {
    let marvel_client = Marvel::new_from_client(client);
    let state = Arc::new(MarvelState {
        tera,
        marvel_client,
    });

    axum::Router::new()
        .route("/comics", get(comics))
        .layer(ServiceBuilder::new().layer(axum::middleware::from_fn(enforce_date_query)))
        .with_state(state)
}

async fn comics<S>(
    State(state): State<Arc<MarvelState<S>>>,
    Query(query): Query<Date>,
    OriginalUri(original_uri): OriginalUri,
) -> Result<Html<String>, Error>
where
    S: WebClient,
{
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", original_uri.path());

    ctx.insert("Date", &query.date.to_string());

    let result = state.marvel_client.weekly_comics(query.date).await?;
    ctx.insert("results", &result);

    Ok(Html(state.tera.render("marvel-unlimited.html", &ctx)?))
}

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

    Ok(axum::response::Redirect::temporary(&format!("{original_uri}?date={date}")).into_response())
}
