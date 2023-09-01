use hyper::{Body, Request};
use std::task::{Context, Poll};
use tower::{Layer, Service};

pub struct UriMiddlewareLayer {
    host: &'static str,
    scheme: &'static str,
}

impl UriMiddlewareLayer {
    pub fn new(host: &'static str, scheme: &'static str) -> Self {
        UriMiddlewareLayer { host, scheme }
    }
}

impl<S> Layer<S> for UriMiddlewareLayer {
    type Service = UriMiddleware<S>;
    fn layer(&self, inner: S) -> Self::Service {
        UriMiddleware::new(inner, self.host, self.scheme)
    }
}

#[derive(Clone)]
pub struct UriMiddleware<S> {
    inner: S,
    host: &'static str,
    scheme: &'static str,
}

impl<S> UriMiddleware<S> {
    fn new(inner: S, host: &'static str, scheme: &'static str) -> Self {
        UriMiddleware {
            inner,
            host,
            scheme,
        }
    }
}

impl<S> Service<Request<Body>> for UriMiddleware<S>
where
    S: Service<Request<Body>>,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = S::Future;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        up.authority = Some(hyper::http::uri::Authority::from_static(self.host));
        up.scheme = Some(self.scheme.parse().unwrap());

        p.uri = hyper::Uri::from_parts(up).unwrap();

        self.inner.call(Request::from_parts(p, b))
    }
}
