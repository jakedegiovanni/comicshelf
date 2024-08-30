use axum::{routing::get, Router};
use tower::ServiceBuilder;
use tower_http::{services::ServeDir, trace::TraceLayer};

use crate::{app::AppState, controllers, middleware};

pub fn build(state: AppState) -> Router {
    Router::new()
        .nest(
            "/comics",
            Router::new().route("/", get(controllers::comics)).layer(
                ServiceBuilder::new()
                    .layer(axum::middleware::from_fn(middleware::enforce_date_query)),
            ),
        )
        .nest(
            "/series",
            Router::new().route("/:series_id", get(controllers::series)),
        )
        .nest_service("/static", ServeDir::new("internal/server/static"))
        .layer(ServiceBuilder::new().layer(TraceLayer::new_for_http()))
        .with_state(state)
}
