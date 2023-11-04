use std::task::{Context, Poll};

use futures_util::future::BoxFuture;
use hyper::{http, Body, Request};
use tower::{BoxError, Layer, Service};

pub struct MiddlewareLayer {
    host: &'static str,
    scheme: http::uri::Scheme,
    path_prefix: &'static str,
}

impl MiddlewareLayer {
    pub fn new(host: &'static str, scheme: http::uri::Scheme, path_prefix: &'static str) -> Self {
        MiddlewareLayer {
            host,
            scheme,
            path_prefix,
        }
    }
}

impl<S> Layer<S> for MiddlewareLayer {
    type Service = Middleware<S>;
    fn layer(&self, inner: S) -> Self::Service {
        Middleware::new(inner, self.host, self.scheme.clone(), self.path_prefix)
    }
}

#[derive(Clone)]
pub struct Middleware<S> {
    inner: S,
    host: &'static str,
    scheme: http::uri::Scheme,
    path_prefix: &'static str,
}

impl<S> Middleware<S> {
    fn new(
        inner: S,
        host: &'static str,
        scheme: http::uri::Scheme,
        path_prefix: &'static str,
    ) -> Self {
        Middleware {
            inner,
            host,
            scheme,
            path_prefix,
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
        let scheme = self.scheme.clone();
        let host = self.host;
        let path_prefix = self.path_prefix;

        Box::pin(async move {
            let (mut p, b) = req.into_parts();

            let mut up = p.uri.clone().into_parts();
            up.authority = Some(http::uri::Authority::from_static(host));
            up.scheme = Some(scheme);

            let pq = p
                .uri
                .into_parts()
                .path_and_query
                .unwrap_or(http::uri::PathAndQuery::from_static(""));
            let path = pq.path();
            let q = pq.query().map_or(String::new(), |q| format!("?{q}"));

            let path = if path.contains(path_prefix) {
                path.to_string()
            } else {
                format!("{path_prefix}{path}")
            };

            up.path_and_query = Some(http::uri::PathAndQuery::try_from(format!("{path}{q}"))?);

            p.uri = hyper::Uri::from_parts(up)?;

            this.call(Request::from_parts(p, b))
                .await
                .map_err(Into::into)
        })
    }
}
