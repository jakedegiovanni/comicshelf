use chrono::NaiveDate;
use serde::Serialize;

#[derive(Debug)]
pub struct HtmlPage;

impl HtmlPage {
    pub fn new_page<T: Serialize>(title: &str, date: NaiveDate, results: &T) -> tera::Context {
        let mut ctx = tera::Context::new();

        ctx.insert("page", results);
        ctx.insert("title", title);
        ctx.insert("date", &date.to_string());

        ctx
    }
}
