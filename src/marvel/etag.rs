use std::collections::HashMap;
use std::sync::{Arc, RwLock};
use std::task::{Context, Poll};

use anyhow::anyhow;
use futures_util::future::BoxFuture;
use hyper::{Body, Request, Response, StatusCode};
use tower::{BoxError, Layer, Service};

use super::template::DataWrapper;

pub type EtagCache = Arc<RwLock<HashMap<String, DataWrapper>>>;

pub fn new_etag_cache() -> EtagCache {
    Arc::new(RwLock::new(HashMap::new()))
}

pub struct EtagMiddlewareLayer {
    cache: EtagCache,
}

impl EtagMiddlewareLayer {
    pub fn new(cache: EtagCache) -> Self {
        EtagMiddlewareLayer { cache }
    }
}

impl<S> Layer<S> for EtagMiddlewareLayer {
    type Service = EtagCacheMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        EtagCacheMiddleware::new(inner, self.cache.clone())
    }
}

#[derive(Clone)]
pub struct EtagCacheMiddleware<S> {
    inner: S,
    cache: EtagCache,
}

impl<S> EtagCacheMiddleware<S> {
    fn new(inner: S, cache: EtagCache) -> Self {
        EtagCacheMiddleware { inner, cache }
    }
}

// todo really want this to follow the S::Error: Into<BoxError> with no Error constraint on the Service trait itself
impl<S> Service<Request<Body>> for EtagCacheMiddleware<S>
where
    S: Service<Request<Body>, Response = Response<Body>, Error = BoxError> + Clone + Send + 'static,
    S::Future: Send,
{
    type Response = DataWrapper;
    type Error = BoxError; // todo can a more specific error be returned here?
    type Future = BoxFuture<'static, Result<DataWrapper, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let this = self.inner.clone();
        let mut this = std::mem::replace(&mut self.inner, this);

        let req = req;
        let cache = self.cache.clone();

        Box::pin(async move {
            let key = req
                .uri()
                .path_and_query()
                .ok_or(anyhow!("no path or query"))?
                .clone();
            let mut headers = req.headers().clone();

            let (mut p, b) = req.into_parts();

            // todo how to properly handle poisoned locks?
            match cache
                .read()
                .expect("could not read from the cache")
                .get(key.as_str())
            {
                Some(wrapper) => {
                    println!(
                        "key {:?} exists in cache, using etag {:?}",
                        key, wrapper.etag
                    );

                    headers.insert("If-None-Match", wrapper.etag.parse()?);
                }
                None => {
                    println!("key {:?} does not exist in cache", key);
                }
            }

            p.headers = headers;
            let req = Request::from_parts(p, b);

            let key = key.to_string();

            let response = this.call(req).await?;
            if response.status() == StatusCode::NOT_MODIFIED {
                println!("using cache");
                Ok(cache
                    .read()
                    .expect("could not read from the cache")
                    .get(key.as_str())
                    .ok_or(anyhow!(
                        "an item expected to be in the cache could not be found"
                    ))?
                    .clone())
            } else {
                let result: DataWrapper = serde_json::from_slice(
                    hyper::body::to_bytes(response.into_body())
                        .await?
                        .iter()
                        .as_slice(),
                )?;
                println!("storing cache");
                cache
                    .write()
                    .expect("could not write to the cache")
                    .insert(key, result.clone());
                Ok(result)
            }
        })
    }
}
