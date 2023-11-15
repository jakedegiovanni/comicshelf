use axum::{routing::get, Router};
use tower::ServiceBuilder;
use tower_http::{services::ServeDir, trace::TraceLayer};

use crate::{app::AppState, marvel, middleware};

pub fn build(state: AppState) -> Router {
    Router::new()
        .nest(
            "/marvel-unlimited",
            Router::new()
                .route("/comics", get(marvel::controllers::comics))
                .layer(
                    ServiceBuilder::new()
                        .layer(axum::middleware::from_fn(middleware::enforce_date_query)),
                ),
        )
        .nest_service("/static", ServeDir::new("static"))
        .layer(ServiceBuilder::new().layer(TraceLayer::new_for_http()))
        .with_state(state)
}
