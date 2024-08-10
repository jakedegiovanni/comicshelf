use std::ops::Deref;
use std::sync::Arc;

use axum::extract::{FromRequestParts, State};
use tera::Tera;

use crate::comicshelf;

pub struct App {
    pub comic_client: Box<dyn comicshelf::Client>,
    pub tera: Tera,
}

impl App {
    pub fn new(client: Box<dyn comicshelf::Client>, tera: Tera) -> Self {
        App {
            comic_client: client,
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
