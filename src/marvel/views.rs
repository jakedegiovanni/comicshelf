use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Date {
    #[serde(rename = "type")]
    typ: String,
    date: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Item {
    name: String,
    #[serde(rename = "resourceURI")]
    resource_uri: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Url {
    #[serde(rename = "type")]
    typ: String,
    url: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Thumbnail {
    path: String,
    extension: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Comic {
    id: i32,
    title: String,
    #[serde(rename = "resourceURI")]
    resource_uri: String,
    urls: Vec<Url>,
    modified: String,
    thumbnail: Thumbnail,
    format: String,
    #[serde(rename = "issueNumber")]
    issue_number: i32,
    series: Item,
    dates: Vec<Date>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DataContainer<T> {
    offset: i32,
    limit: i32,
    total: i32,
    count: i32,
    results: Vec<T>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DataWrapper<T> {
    code: i32,
    status: String,
    copyright: String,
    #[serde(rename = "attributionText")]
    attribution_text: String,
    #[serde(rename = "attributionHTML")]
    attribution_html: String,
    pub etag: String,
    data: DataContainer<T>,
}

#[derive(Serialize, Debug)]
pub struct MarvelUnlimited {
    endpoint: String,
    date: String,
    results: DataWrapper<Comic>,
}

impl MarvelUnlimited {
    pub fn new(endpoint: String, date: String, results: DataWrapper<Comic>) -> MarvelUnlimited {
        MarvelUnlimited {
            endpoint,
            date,
            results,
        }
    }

    pub fn new_page(endpoint: String, date: String, results: DataWrapper<Comic>) -> tera::Context {
        let mut ctx = tera::Context::new();
        ctx.insert("page", &Self::new(endpoint, date, results));
        ctx
    }
}
