use axum::response::IntoResponse;
use hyper::StatusCode;
use thiserror::Error;

pub mod uri;

#[derive(Error, Debug)]
pub enum Error {}

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("something weent wrong: {self}"),
        )
            .into_response()
    }
}
