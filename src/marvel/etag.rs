use std::collections::HashMap;
use std::sync::{Arc, RwLock};
use std::task::{Context, Poll};

use anyhow::anyhow;
use futures_util::future::BoxFuture;
use hyper::{Body, Request, Response, StatusCode};
use tower::{BoxError, Layer, Service};

use super::template::DataWrapper;

pub type Cache = Arc<RwLock<HashMap<String, DataWrapper>>>;

pub fn new_etag_cache() -> Cache {
    Arc::new(RwLock::new(HashMap::new()))
}

pub struct CacheMiddlewareLayer {
    cache: Cache,
}

impl CacheMiddlewareLayer {
    pub fn new(cache: Cache) -> Self {
        CacheMiddlewareLayer { cache }
    }
}

impl<S> Layer<S> for CacheMiddlewareLayer {
    type Service = CacheMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        CacheMiddleware::new(inner, self.cache.clone())
    }
}

#[derive(Clone)]
pub struct CacheMiddleware<S> {
    inner: S,
    cache: Cache,
}

impl<S> CacheMiddleware<S> {
    fn new(inner: S, cache: Cache) -> Self {
        CacheMiddleware { inner, cache }
    }
}

impl<S> Service<Request<Body>> for CacheMiddleware<S>
where
    S: Service<Request<Body>, Response = Response<Body>> + Clone + Send + 'static,
    S::Error: Into<BoxError>,
    S::Future: Send,
{
    type Response = DataWrapper;
    type Error = BoxError; // todo can a more specific error be returned here?
    type Future = BoxFuture<'static, Result<DataWrapper, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
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
            let key = key.as_str();
            let mut headers = req.headers().clone();

            let (mut p, b) = req.into_parts();

            // todo how to properly handle poisoned locks?
            match cache
                .read()
                .expect("could not read from the cache")
                .get(key)
            {
                Some(wrapper) => {
                    println!(
                        "key {:?} exists in cache, using etag {:?}",
                        key, wrapper.etag
                    );

                    headers.insert("If-None-Match", wrapper.etag.parse()?);
                }
                None => {
                    println!("key {key:?} does not exist in cache");
                }
            }

            p.headers = headers;
            let req = Request::from_parts(p, b);

            let response = this.call(req).await.map_err(Into::into)?;
            if response.status() == StatusCode::NOT_MODIFIED {
                println!("using cache");
                Ok(cache
                    .read()
                    .expect("could not read from the cache")
                    .get(key)
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
                    .insert(key.to_owned(), result.clone());
                Ok(result)
            }
        })
    }
}
