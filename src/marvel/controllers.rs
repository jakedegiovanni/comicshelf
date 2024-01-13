use axum::{
    debug_handler,
    extract::{OriginalUri, Query},
    response::Html,
};

use crate::{app::AppState, errors::Error, middleware::Date};

use super::views::MarvelUnlimited;

#[debug_handler(state = AppState)]
pub async fn comics(
    state: AppState,
    Query(query): Query<Date>,
    OriginalUri(original_uri): OriginalUri,
) -> Result<Html<String>, Error> {
    let result = state.marvel_client.weekly_comics(query.date).await?;

    let page = MarvelUnlimited::new_page(
        original_uri.path().to_string(),
        query.date.to_string(),
        result,
    );

    Ok(Html(state.tera.render("marvel-unlimited.html", &page)?))
}
