use axum::{extract::State, http::StatusCode, response::Html, routing::get};
use std::sync::Arc;
use tera::{Context, Tera};

struct ComicShelf {
    tera: Tera,
}

async fn index(state: State<Arc<ComicShelf>>) -> Result<Html<String>, StatusCode> {
    let mut ctx = Context::new();
    ctx.insert("name", "World");

    let body = match state.tera.render("index1.html", &ctx) {
        Ok(b) => b,
        Err(e) => {
            println!("tera rendering error: {}", e);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };

    Ok(Html(body))
}

#[tokio::main]
async fn main() {
    let tera = match Tera::new("templates/**/*.html") {
        Ok(t) => t,
        Err(e) => {
            println!("Parsing error(s): {}", e);
            std::process::exit(1);
        }
    };

    let state = Arc::new(ComicShelf { tera });
    let app = axum::Router::new().route("/", get(index)).with_state(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
