use axum::{http::uri::InvalidUri, response::IntoResponse, BoxError};
use hyper::StatusCode;
use thiserror::Error;

pub mod uri;

#[derive(Error, Debug)]
pub enum Error {
    #[error("internal error")]
    Anyhow(#[from] anyhow::Error),
    #[error("rendering error")]
    Tera(#[from] tera::Error),
    #[error("box error")]
    Box(#[from] BoxError),
    #[error("hyper error")]
    HyperError(#[from] hyper::http::Error),
    #[error("uri error")]
    UriError(#[from] InvalidUri),
}

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("something weent wrong: {self}"),
        )
            .into_response()
    }
}
