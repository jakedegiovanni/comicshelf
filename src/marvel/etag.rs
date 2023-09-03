use std::collections::HashMap;
use std::fmt::Debug;
use std::sync::{Arc, RwLock};
use std::task::{Context, Poll};

use futures_util::future::BoxFuture;
use hyper::{Body, Request, Response, StatusCode};
use tower::{Layer, Service};

use crate::marvel::template::DataWrapper;

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

impl<S> Service<Request<Body>> for EtagCacheMiddleware<S>
where
    S: Service<Request<Body>, Response = Response<Body>> + Clone + Send + 'static,
    S::Error: Debug,
    S::Future: Send,
{
    type Response = DataWrapper;
    type Error = S::Error; // todo use anyhow::Error
    type Future = BoxFuture<'static, Result<DataWrapper, S::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let key = req.uri().path_and_query().unwrap().clone();
        let mut headers = req.headers().clone();

        let (mut p, b) = req.into_parts();

        match self.cache.read().unwrap().get(key.as_str()) {
            Some(wrapper) => {
                println!(
                    "key {:?} exists in cache, using etag {:?}",
                    key, wrapper.etag
                );

                headers.insert("If-None-Match", wrapper.etag.parse().unwrap());
            }
            None => {
                println!("key {:?} does not exist in cache", key);
            }
        }

        p.headers = headers;
        let req = Request::from_parts(p, b);

        let cache = self.cache.clone();
        let key = key.to_string();

        let future = self.inner.call(req);
        Box::pin(async move {
            let response: Response<Body> = future.await?;
            if response.status() == StatusCode::NOT_MODIFIED {
                println!("using cache");
                Ok(cache.read().unwrap().get(key.as_str()).unwrap().clone())
            } else {
                let result: DataWrapper = serde_json::from_slice(
                    hyper::body::to_bytes(response.into_body())
                        .await
                        .unwrap()
                        .iter()
                        .as_slice(),
                )
                .unwrap();
                println!("storing cache");
                cache.write().unwrap().insert(key, result.clone());
                Ok(result)
            }
        })
    }
}
