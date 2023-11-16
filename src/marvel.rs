use std::{
    collections::HashMap,
    sync::{Arc, RwLock},
};

use anyhow::anyhow;
use async_trait::async_trait;
use axum::http::HeaderValue;
use chrono::{Datelike, Days, Months, NaiveDate, Utc, Weekday};
use hyper::{HeaderMap, StatusCode};

use crate::marvel::views::DataWrapper;

pub mod controllers;
pub mod views;

#[async_trait]
pub trait Client: Send + Sync {
    async fn weekly_comics(&self, date: NaiveDate) -> Result<DataWrapper, anyhow::Error>;
}

#[derive(Debug)]
pub struct RealClient {
    client: reqwest::Client,
    pub_key: &'static str,
    priv_key: &'static str,
    base_url: &'static str,
    cache: Arc<RwLock<HashMap<String, DataWrapper>>>,
}

impl RealClient {
    pub fn new(
        client: reqwest::Client,
        pub_key: &'static str,
        priv_key: &'static str,
        base_url: &'static str,
    ) -> Self {
        RealClient {
            client,
            pub_key,
            priv_key,
            base_url,
            cache: Arc::new(RwLock::new(HashMap::new())),
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
impl Client for RealClient {
    async fn weekly_comics(&self, date: NaiveDate) -> Result<DataWrapper, anyhow::Error> {
        let (date, date2) = RealClient::week_range(date);
        let endpoint = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={date},{date2}&hasDigitalIssue=true&orderBy=issueNumber&limit=100");
        let uri = self.uri(&endpoint);

        let mut headers = HeaderMap::new();
        if let Some(wrapper) = self
            .cache
            .read()
            .expect("could not read from the cache")
            .get(&endpoint)
        {
            println!(
                "key {:?} exists in cache using etag {:?}",
                &endpoint, wrapper.etag
            );
            headers.insert(
                "If-None-Match",
                HeaderValue::from_str(&wrapper.etag).map_err(|e| anyhow!(e))?,
            );
        };

        let response = self
            .client
            .get(uri)
            .headers(headers)
            .send()
            .await
            .map_err(|e| anyhow!(e))?
            .error_for_status()
            .map_err(|e| anyhow!(e))?;

        if response.status() == StatusCode::NOT_MODIFIED {
            let wrapper = self
                .cache
                .read()
                .expect("could not read from cache")
                .get(&endpoint)
                .ok_or(anyhow!(
                    "an item expeected to be in the cache could not be found"
                ))?
                .clone();

            println!("etag {:?} not modified, using cache", wrapper.etag);
            return Ok(wrapper);
        }

        let result = response
            .json::<DataWrapper>()
            .await
            .map_err(|e| anyhow!(e))?;

        self.cache
            .write()
            .expect("could not writ eto the cache")
            .insert(endpoint, result.clone());

        Ok(result)
    }
}
