use std::ops::Deref;
use std::sync::Arc;

use axum::extract::{FromRequestParts, State};
use tera::Tera;

use crate::marvel;

pub struct App {
    pub marvel_client: Box<dyn marvel::Client>,
    pub tera: Tera,
}

impl App {
    pub fn new(client: reqwest::Client, tera: Tera) -> Self {
        let marvel_client = Box::new(marvel::RealClient::new(
            client,
            include_str!("../pub.txt"),
            include_str!("../priv.txt"),
            "https://gateway.marvel.com/v1/public",
        ));

        App {
            marvel_client,
            tera,
        }
    }
}

#[allow(clippy::module_name_repetitions)]
#[derive(Clone, FromRequestParts)]
#[from_request(via(State))]
pub struct AppState(pub Arc<App>);

// deref so you can still access the inner fields easily
impl Deref for AppState {
    type Target = App;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}
