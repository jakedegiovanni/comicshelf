use std::{
    collections::HashMap,
    future::Future,
    sync::{Arc, RwLock},
};

use crate::comicshelf;
use crate::comicshelf::Page;
use anyhow::{anyhow, Error};
use async_trait::async_trait;
use axum::http::HeaderValue;
use chrono::{DateTime, Datelike, Days, Months, NaiveDate, Utc, Weekday};
use hyper::{HeaderMap, StatusCode};
use serde::{de::DeserializeOwned, Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Date {
    #[serde(rename = "type")]
    pub typ: String,
    pub date: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Item {
    pub name: String,
    #[serde(rename = "resourceURI")]
    pub resource_uri: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Url {
    #[serde(rename = "type")]
    pub typ: String,
    pub url: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Thumbnail {
    pub path: String,
    pub extension: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Comic {
    pub id: i64,
    pub title: String,
    #[serde(rename = "resourceURI")]
    pub resource_uri: String,
    pub urls: Vec<Url>,
    pub modified: String,
    pub thumbnail: Thumbnail,
    pub format: String,
    #[serde(rename = "issueNumber")]
    pub issue_number: i64,
    pub series: Item,
    pub dates: Vec<Date>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DataContainer<T> {
    pub offset: i64,
    pub limit: i64,
    pub total: i64,
    pub count: i64,
    pub results: Vec<T>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DataWrapper<T> {
    pub code: i64,
    pub status: String,
    pub copyright: String,
    #[serde(rename = "attributionText")]
    pub attribution_text: String,
    #[serde(rename = "attributionHTML")]
    pub attribution_html: String,
    pub etag: String,
    pub data: DataContainer<T>,
}

impl<T> Cacheable for DataWrapper<T> {
    fn get_key(&self) -> String {
        self.etag.clone()
    }
}

trait Cacheable {
    fn get_key(&self) -> String;
}

#[derive(Debug)]
struct EtagCache<T: Cacheable> {
    cache: Arc<RwLock<HashMap<String, T>>>,
}

impl<T: Cacheable + DeserializeOwned + Clone> EtagCache<T> {
    fn new() -> Self {
        EtagCache {
            cache: Arc::new(RwLock::new(HashMap::<String, T>::new())),
        }
    }

    async fn retrieve<F, Fut>(&self, endpoint: &str, func: F) -> Result<T, anyhow::Error>
    where
        Fut: Future<Output = Result<reqwest::Response, anyhow::Error>>,
        F: FnOnce(HeaderMap) -> Fut,
    {
        let mut headers = HeaderMap::new();
        if let Some(wrapper) = self
            .cache
            .read()
            .expect("could not read from the cache")
            .get(endpoint)
        {
            println!(
                "key {:?} exists in cache using etag {:?}",
                &endpoint,
                wrapper.get_key()
            );

            headers.insert(
                "If-None-Match",
                HeaderValue::from_str(wrapper.get_key().as_str()).map_err(|e| anyhow!(e))?,
            );
        };

        let response = func(headers).await?;

        if response.status() == StatusCode::NOT_MODIFIED {
            let wrapper = self
                .cache
                .read()
                .expect("could not read from cache")
                .get(endpoint)
                .ok_or(anyhow!(
                    "an item expeected to be in the cache could not be found"
                ))?
                .clone();

            println!("etag {:?} not modified, using cache", wrapper.get_key());
            return Ok(wrapper);
        }

        let result = response.json::<T>().await.map_err(|e| anyhow!(e))?;

        self.cache
            .write()
            .expect("could not writ eto the cache")
            .insert(endpoint.to_string(), result.clone());

        Ok(result)
    }
}

impl Into<comicshelf::Page<comicshelf::Comic>> for DataWrapper<Comic> {
    fn into(self) -> Page<comicshelf::Comic> {
        let mut comics = Vec::<comicshelf::Comic>::new();

        for comic in self.data.results {
            let mut urls = Vec::<comicshelf::Url>::new();
            for url in comic.urls {
                urls.push(comicshelf::Url::new(url.typ, url.url));
            }

            let mut on_sale_date: NaiveDate = Utc::now().date_naive();
            for date in comic.dates {
                if !date.typ.eq_ignore_ascii_case("onsaledate") {
                    continue;
                }

                on_sale_date = DateTime::parse_from_str(&date.date, "%Y-%m-%dT%T%z")
                    .expect("could not parse marvel time")
                    .date_naive();
            }

            let re = regex::Regex::new("/([0-9]+)/?").expect("could not compile regex");
            let Some((_, [i])) = re.captures(&comic.series.resource_uri).map(|c| c.extract())
            else {
                todo!()
            };

            let series_id: i64 = i.parse::<i64>().unwrap();

            comics.push(comicshelf::Comic::new(
                comic.id,
                comic.title,
                urls,
                format!(
                    "{}/portrait_uncanny.{}",
                    comic.thumbnail.path, comic.thumbnail.extension
                ),
                comic.format,
                comic.issue_number,
                on_sale_date,
                self.attribution_text.clone(),
                "https://marvel.com".to_string(),
                series_id,
            ));
        }

        Page::new(
            self.data.limit,
            self.data.total,
            self.data.count,
            self.data.offset,
            comics,
        )
    }
}

#[derive(Debug)]
pub struct Client {
    #[allow(clippy::struct_field_names)]
    http_client: reqwest::Client,
    pub_key: &'static str,
    priv_key: &'static str,
    base_url: &'static str,
    comic_cache: EtagCache<DataWrapper<Comic>>,
}

impl Client {
    pub fn new(
        client: reqwest::Client,
        pub_key: &'static str,
        priv_key: &'static str,
        base_url: &'static str,
    ) -> Self {
        Client {
            http_client: client,
            pub_key,
            priv_key,
            base_url,
            comic_cache: EtagCache::<DataWrapper<Comic>>::new(),
        }
    }

    fn week_range(time: NaiveDate) -> (NaiveDate, NaiveDate) {
        let mut date = time - Months::new(3);

        let one_day = Days::new(1);
        loop {
            if date.weekday() == Weekday::Sun {
                break;
            }

            date = date - one_day;
        }

        (date - Days::new(7), date - one_day)
    }

    fn auth(&self) -> String {
        let priv_key = self.priv_key;
        let pub_key = self.pub_key;

        let ts = Utc::now().timestamp_millis();
        let hash = format!("{:x}", md5::compute(format!("{ts}{priv_key}{pub_key}")));

        format!("apikey={pub_key}&ts={ts}&hash={hash}")
    }

    fn uri(&self, endpoint: &str) -> String {
        format!("{}{endpoint}&{}", self.base_url, self.auth())
    }
}

#[async_trait]
impl comicshelf::ComicClient for Client {
    async fn weekly_comics(
        &self,
        date: NaiveDate,
    ) -> Result<comicshelf::Page<comicshelf::Comic>, anyhow::Error> {
        let (date, date2) = Client::week_range(date);
        let endpoint = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={date},{date2}&hasDigitalIssue=true&orderBy=issueNumber&limit=100");
        let uri = self.uri(&endpoint);

        let result = self
            .comic_cache
            .retrieve(&endpoint, |headers| async {
                self.http_client
                    .get(uri)
                    .headers(headers)
                    .send()
                    .await
                    .map_err(|e| anyhow!(e))?
                    .error_for_status()
                    .map_err(|e| anyhow!(e))
            })
            .await?;

        Ok(result.into())
    }

    async fn get_comic(&self, id: i64) -> Result<comicshelf::Comic, Error> {
        todo!()
    }
}

#[async_trait]
impl comicshelf::SeriesClient for Client {
    async fn get_comics_within_series(
        &self,
        id: i64,
    ) -> Result<Vec<comicshelf::Comic>, anyhow::Error> {
        todo!();
    }

    async fn get_series(&self, id: i64) -> Result<comicshelf::Series, anyhow::Error> {
        todo!();
    }
}

#[async_trait]
impl comicshelf::Client for Client {}
