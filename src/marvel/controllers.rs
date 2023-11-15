use axum::{
    extract::{OriginalUri, Query},
    response::Html,
};
use tera::Context;

use crate::{app::AppState, errors::Error, middleware::Date};

pub async fn comics(
    state: AppState,
    Query(query): Query<Date>,
    OriginalUri(original_uri): OriginalUri,
) -> Result<Html<String>, Error> {
    let mut ctx = Context::new();
    ctx.insert("PageEndpoint", original_uri.path());

    ctx.insert("Date", &query.date.to_string());

    let result = state.marvel_client.weekly_comics(query.date).await?;
    ctx.insert("results", &result);

    Ok(Html(state.tera.render("marvel-unlimited.html", &ctx)?))
}
