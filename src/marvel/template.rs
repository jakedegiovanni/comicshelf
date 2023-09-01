use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Serialize, Deserialize, Debug)]
pub struct DataContainer {
    offset: i32,
    limit: i32,
    total: i32,
    count: i32,
    results: Value,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize, Debug)]
pub struct DataWrapper {
    code: Value,
    status: String,
    copyright: String,
    attributionText: String,
    attributionHTML: String,
    etag: String,
    data: DataContainer,
}