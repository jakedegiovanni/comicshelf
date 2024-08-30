use axum::{
    debug_handler,
    extract::{Path, Query},
    response::Html,
};

use crate::{app::AppState, errors::Error, middleware::Date, views::HtmlPage};

#[debug_handler(state = AppState)]
pub async fn comics(state: AppState, Query(query): Query<Date>) -> Result<Html<String>, Error> {
    let result = state.comic_client.weekly_comics(query.date).await?;

    let page = HtmlPage::new_page("Weekly Comics", query.date, &result);

    Ok(Html(state.tera.render("comic-page.html", &page)?))
}

#[debug_handler(state = AppState)]
pub async fn series(
    state: AppState,
    Path(series_id): Path<i64>,
    Query(query): Query<Date>,
) -> Result<Html<String>, Error> {
    let results = state
        .comic_client
        .get_comics_within_series(series_id)
        .await?;

    let page = HtmlPage::new_page("Series Issues", query.date, &results);

    Ok(Html(state.tera.render("comic-page.html", &page)?))
}
