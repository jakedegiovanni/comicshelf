use std::collections::HashMap;
use std::sync::Arc;

use app::{App, AppState};

use serde_json::Value;
use tera::Tera;

mod app;
mod comicshelf;
mod controllers;
mod errors;
mod marvel;
mod middleware;
mod router;
mod views;

fn following(args: &HashMap<String, Value>) -> tera::Result<Value> {
    let _ = args.get("index").ok_or(tera::Error::msg("not found"))?; // todo use for db check
    Ok(tera::to_value(false)?)
}

#[tokio::main]
async fn main() {
    let mut tera = Tera::new("templates/**/*.html").unwrap();
    tera.register_function("following", following);

    let client = reqwest::Client::new();
    let marvel_client = marvel::Client::new(
        client,
        include_str!("../marvel/pub.txt"),
        include_str!("../marvel/priv.txt"),
        "https://gateway.marvel.com/v1/public",
    );

    let state = AppState(Arc::new(App::new(Box::new(marvel_client), tera)));

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    let router = router::build(state);

    axum::Server::bind(&"127.0.0.1:8080".parse().unwrap())
        .serve(router.into_make_service())
        .await
        .unwrap();
}
