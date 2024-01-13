use axum::{http::uri::InvalidUri, response::IntoResponse};
use hyper::StatusCode;
use thiserror::Error;
use tower::BoxError;

#[derive(Error, Debug)]
pub enum Error {
    #[error("internal error")]
    Anyhow(#[from] anyhow::Error),
    #[error("rendering error: Error: {0}")]
    Tera(#[from] tera::Error),
    #[error("box error")]
    Box(#[from] BoxError),
    #[error("hyper error")]
    Hyper(#[from] hyper::http::Error),
    #[error("uri error")]
    Uri(#[from] InvalidUri),
}

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        let msg = format!("something went wrong: {self}");
        println!("{msg}");

        (StatusCode::INTERNAL_SERVER_ERROR, msg).into_response()
    }
}
