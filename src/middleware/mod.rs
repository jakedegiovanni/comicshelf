use std::future::Future;
use std::pin::Pin;
use std::task::{Context, Poll};
use pin_project::pin_project;
use tower::BoxError;

pub mod uri;

#[pin_project]
pub struct MiddlewareFuture<F> {
    #[pin]
    future: F
}

impl<F> MiddlewareFuture<F> {
    pub fn new(future: F) -> Self {
        MiddlewareFuture { future }
    }
}

impl<F, R, E> Future for MiddlewareFuture<F>
where
    F: Future<Output = Result<R, E>>,
    E: Into<BoxError>,
{
    type Output = Result<R, BoxError>;

    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output> {
        let this = self.project();
        match this.future.poll(cx) {
            Poll::Ready(result) => Poll::Ready(result.map_err(Into::into)),
            Poll::Pending => Poll::Pending
        }
    }
}