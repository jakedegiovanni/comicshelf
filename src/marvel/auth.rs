use std::task::{Context, Poll};

use chrono::Utc;
use hyper::{Body, Request};
use tower::{BoxError, Layer, Service};

use crate::middleware::MiddlewareFuture;

pub struct AuthMiddlewareLayer {
    pub_key: &'static str,
    priv_key: &'static str,
}

impl AuthMiddlewareLayer {
    pub fn new(pub_key: &'static str, priv_key: &'static str) -> Self {
        AuthMiddlewareLayer { pub_key, priv_key }
    }
}

impl<S> Layer<S> for AuthMiddlewareLayer {
    type Service = AuthMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        AuthMiddleware::new(inner, self.pub_key, self.priv_key)
    }
}

#[derive(Clone)]
pub struct AuthMiddleware<S> {
    inner: S,
    pub_key: &'static str,
    priv_key: &'static str,
}

impl<S> AuthMiddleware<S> {
    fn new(inner: S, pub_key: &'static str, priv_key: &'static str) -> Self {
        AuthMiddleware {
            inner,
            pub_key,
            priv_key,
        }
    }
}

impl<S> Service<Request<Body>> for AuthMiddleware<S>
where
    S: Service<Request<Body>>,
    S::Error: Into<BoxError>
{
    type Response = S::Response;
    type Error = BoxError;
    type Future = MiddlewareFuture<S::Future>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        let pq = up.path_and_query.unwrap();
        let path = pq.path();
        let q = pq.query().unwrap_or("");

        let ts = Utc::now().timestamp_millis();
        let hash = format!(
            "{:x}",
            md5::compute(format!("{}{}{}", ts, self.priv_key, self.pub_key))
        );
        let query = format!("apikey={}&ts={}&hash={}", self.pub_key, ts, hash);
        let query = if q.is_empty() {
            format!("?{}", query)
        } else {
            format!("{}&{}", q, query)
        };
        let query = format!("{}?{}", path, query);
        up.path_and_query = Some(hyper::http::uri::PathAndQuery::try_from(query).unwrap());

        p.uri = hyper::Uri::from_parts(up).unwrap();
        MiddlewareFuture::new(self.inner.call(Request::from_parts(p, b)))
    }
}
