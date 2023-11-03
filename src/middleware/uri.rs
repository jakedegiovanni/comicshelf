use std::task::{Context, Poll};

use hyper::{Body, http, Request};
use tower::{BoxError, Layer, Service};

use crate::middleware::MiddlewareFuture;

pub struct UriMiddlewareLayer {
    host: &'static str,
    scheme: http::uri::Scheme,
    path_prefix: &'static str
}

impl UriMiddlewareLayer {
    pub fn new(host: &'static str, scheme: http::uri::Scheme, path_prefix: &'static str) -> Self {
        UriMiddlewareLayer { host, scheme, path_prefix }
    }
}

impl<S> Layer<S> for UriMiddlewareLayer {
    type Service = UriMiddleware<S>;
    fn layer(&self, inner: S) -> Self::Service {
        UriMiddleware::new(inner, self.host, self.scheme.clone(), self.path_prefix)
    }
}

#[derive(Clone)]
pub struct UriMiddleware<S> {
    inner: S,
    host: &'static str,
    scheme: http::uri::Scheme,
    path_prefix: &'static str,
}

impl<S> UriMiddleware<S> {
    fn new(inner: S, host: &'static str, scheme: http::uri::Scheme, path_prefix: &'static str) -> Self {
        UriMiddleware {
            inner,
            host,
            scheme,
            path_prefix,
        }
    }
}

impl<S> Service<Request<Body>> for UriMiddleware<S>
where
    S: Service<Request<Body>>,
    S::Error: Into<BoxError>,
{
    type Response = S::Response;
    type Error = BoxError;
    type Future = MiddlewareFuture<S::Future>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.clone().into_parts();
        up.authority = Some(http::uri::Authority::from_static(self.host));
        up.scheme = Some(self.scheme.clone());

        let pq = p.uri.into_parts().path_and_query.unwrap_or(
            http::uri::PathAndQuery::from_static("")
        );
        let path = pq.path();
        let q = pq.query().map_or(
            "".to_string(),
            |q| format!("?{}", q)
        );

        let path = if !path.contains(self.path_prefix) {
            format!("{}{}", self.path_prefix, path)
        } else {
            path.to_string()
        };

        up.path_and_query = Some(http::uri::PathAndQuery::try_from(format!("{}{}", path, q)).unwrap());

        p.uri = hyper::Uri::from_parts(up).unwrap();

        MiddlewareFuture::new(self.inner.call(Request::from_parts(p, b)))
    }
}
