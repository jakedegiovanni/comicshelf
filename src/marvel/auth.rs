use std::task::{Context, Poll};

use axum::http::uri::PathAndQuery;
use chrono::Utc;
use futures_util::future::BoxFuture;
use hyper::{Body, Request};
use tower::{BoxError, Layer, Service};

pub struct MiddlewareLayer {
    pub_key: &'static str,
    priv_key: &'static str,
}

impl MiddlewareLayer {
    pub fn new(pub_key: &'static str, priv_key: &'static str) -> Self {
        MiddlewareLayer { pub_key, priv_key }
    }
}

impl<S> Layer<S> for MiddlewareLayer {
    type Service = Middleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        Middleware::new(inner, self.pub_key, self.priv_key)
    }
}

#[derive(Clone)]
pub struct Middleware<S> {
    inner: S,
    pub_key: &'static str,
    priv_key: &'static str,
}

impl<S> Middleware<S> {
    fn new(inner: S, pub_key: &'static str, priv_key: &'static str) -> Self {
        Middleware {
            inner,
            pub_key,
            priv_key,
        }
    }
}

impl<S> Service<Request<Body>> for Middleware<S>
where
    S: Service<Request<Body>> + Clone + Send + 'static,
    S::Error: Into<BoxError>,
    S::Future: Send,
{
    type Response = S::Response;
    type Error = BoxError;
    type Future = BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let this = self.inner.clone();
        let mut this = std::mem::replace(&mut self.inner, this);

        let req = req;
        let priv_key = self.priv_key;
        let pub_key = self.pub_key;

        Box::pin(async move {
            let (mut p, b) = req.into_parts();

            let mut up = p.uri.into_parts();
            let pq = up.path_and_query.unwrap_or(PathAndQuery::from_static("")); // todo test creating empty path and query to understand behaviour
            let path = pq.path();
            let q = pq.query().unwrap_or("");

            let ts = Utc::now().timestamp_millis();
            let hash = format!("{:x}", md5::compute(format!("{ts}{priv_key}{pub_key}")));
            let query = format!("apikey={pub_key}&ts={ts}&hash={hash}");
            let query = if q.is_empty() {
                format!("?{query}")
            } else {
                format!("{q}&{query}")
            };
            let query = format!("{path}?{query}");
            up.path_and_query = Some(hyper::http::uri::PathAndQuery::try_from(query)?);

            p.uri = hyper::Uri::from_parts(up)?;

            this.call(Request::from_parts(p, b))
                .await
                .map_err(Into::into)
        })
    }
}
