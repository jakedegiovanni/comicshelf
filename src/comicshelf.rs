use async_trait::async_trait;
use chrono::NaiveDate;
use serde::Serialize;

#[derive(Serialize, Debug, Clone)]
pub struct Page<T> {
    pub limit: i64,
    pub total: i64,
    pub count: i64,
    pub offset: i64,
    pub results: Vec<T>,
}

impl<T> Page<T> {
    pub fn new(limit: i64, total: i64, count: i64, offset: i64, results: Vec<T>) -> Self {
        Page {
            limit,
            total,
            count,
            offset,
            results,
        }
    }
}

#[derive(Serialize, Debug, Clone)]
pub struct Url {
    pub r#type: String,
    pub url: String,
}

impl Url {
    pub fn new(r#type: String, url: String) -> Self {
        Url { r#type, url }
    }
}

#[derive(Serialize, Debug, Clone)]
pub struct Comic {
    pub id: i64,
    pub title: String,
    pub urls: Vec<Url>,
    pub thumbnail: String,
    pub format: String,
    pub issue_number: i64,
    pub on_sale_date: NaiveDate,
    pub attribution: String,
    pub attribution_link: String,
    pub series_id: i64,
}

impl Comic {
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        id: i64,
        title: String,
        urls: Vec<Url>,
        thumbnail: String,
        format: String,
        issue_number: i64,
        on_sale_date: NaiveDate,
        attribution: String,
        attribution_link: String,
        series_id: i64,
    ) -> Self {
        Comic {
            id,
            title,
            urls,
            thumbnail,
            format,
            issue_number,
            on_sale_date,
            attribution,
            attribution_link,
            series_id,
        }
    }
}

#[derive(Serialize, Debug, Clone)]
pub struct Series {
    comics: Vec<Comic>,
    id: i64,
    title: String,
    urls: Vec<Url>,
    thumbnail: String,
}

#[async_trait]
pub trait ComicClient: Send + Sync {
    async fn weekly_comics(&self, date: NaiveDate) -> Result<Page<Comic>, anyhow::Error>;
}

#[async_trait]
pub trait SeriesClient: Send + Sync {
    async fn get_comics_within_series(&self, series_id: i64) -> Result<Page<Comic>, anyhow::Error>;
}

#[async_trait]
pub trait Client: ComicClient + SeriesClient {}
