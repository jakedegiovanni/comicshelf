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
pub struct DataContainer {
    offset: i32,
    limit: i32,
    total: i32,
    count: i32,
    results: Vec<Comic>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DataWrapper {
    code: i32,
    status: String,
    copyright: String,
    #[serde(rename = "attributionText")]
    attribution_text: String,
    #[serde(rename = "attributionHTML")]
    attribution_html: String,
    pub etag: String,
    data: DataContainer,
}
