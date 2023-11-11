use chrono::{Datelike, Days, Months, NaiveDate, Weekday};
use futures_util::future::BoxFuture;
use hyper::{client::HttpConnector, Body, Request};
use hyper_tls::HttpsConnector;
use tower::{BoxError, Service, ServiceBuilder};

use crate::middleware::uri;

use self::template::DataWrapper;

pub mod auth;
pub mod etag;
pub mod template;

pub trait WebClient:
    Service<
        Request<Body>,
        Response = DataWrapper,
        Error = BoxError,
        Future = BoxFuture<'static, Result<DataWrapper, BoxError>>,
    > + Send
    + Sync
    + Clone
{
}

impl<S> WebClient for S where
    S: Service<
            Request<Body>,
            Response = DataWrapper,
            Error = BoxError,
            Future = BoxFuture<'static, Result<DataWrapper, BoxError>>,
        > + Send
        + Sync
        + Clone
{
}

fn new_web_client(client: &hyper::Client<HttpsConnector<HttpConnector>>) -> impl WebClient {
    ServiceBuilder::new()
        .layer(etag::CacheMiddlewareLayer::new(etag::new_etag_cache()))
        .layer(uri::MiddlewareLayer::new(
            "gateway.marvel.com",
            hyper::http::uri::Scheme::HTTPS,
            "/v1/public",
        ))
        .layer(auth::MiddlewareLayer::new(
            include_str!("../../pub.txt"), // todo: formalize, this is janky
            include_str!("../../priv.txt"),
        ))
        .service(client.clone())
}

pub fn new_marvel_service(
    client: &hyper::Client<HttpsConnector<HttpConnector>>,
) -> Marvel<impl WebClient> {
    Marvel::new(new_web_client(client))
}

pub struct Marvel<S> {
    svc: S,
}

impl<S> Marvel<S> {
    pub fn new(svc: S) -> Self {
        Marvel { svc }
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
}

impl<S> Marvel<S>
where
    S: WebClient,
{
    pub async fn weekly_comics(&self, date: NaiveDate) -> Result<DataWrapper, BoxError> {
        let (date, date2) = Marvel::<S>::week_range(date);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={date},{date2}&hasDigitalIssue=true&orderBy=issueNumber&limit=100");
        let req = Request::get(uri).body(Body::empty())?;

        let mut svc = self.svc.clone();
        svc.call(req).await
    }
}
