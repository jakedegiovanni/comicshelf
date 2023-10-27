use chrono::{Datelike, DateTime, Days, Months, Utc, Weekday};
use futures_util::future::BoxFuture;
use hyper::{Body, Request};
use tower::Service;

use self::template::DataWrapper;

pub mod auth;
pub mod etag;
pub mod template;

pub trait MarvelService: Service<
    Request<Body>,
    Response = DataWrapper,
    Error = hyper::Error,
    Future = BoxFuture<'static, Result<DataWrapper, hyper::Error>>
> + Send + Sync + Clone {}

impl<S> MarvelService for S
where
S: Service<
    Request<Body>,
    Response = DataWrapper,
    Error = hyper::Error,
    Future = BoxFuture<'static, Result<DataWrapper, hyper::Error>>
> + Send + Sync + Clone
{}

pub struct Marvel<S>
{
    svc: S
}

impl<S> Marvel<S> {
    pub fn new(svc: S) -> Self {
        Marvel { svc }
    }
}

impl<S> Marvel<S>
where
    S: MarvelService
{

    pub async fn weekly_comics(&self) -> DataWrapper {
        let (date, date2) = self.week_range(Utc::now());
        let date = self.fmt_date(&date);
        let date2 = self.fmt_date(&date2);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={},{}&hasDigitalIssue=true&orderBy=issueNumber&limit=100", date, date2);
        let req = Request::get(uri).body(Body::empty()).unwrap();

        let mut svc = self.svc.clone();
        svc.call(req).await.unwrap()
    }

    fn week_range(&self, time: DateTime<Utc>) -> (DateTime<Utc>, DateTime<Utc>) {
        let mut date = time.checked_sub_months(Months::new(3)).unwrap();
        loop {
            if date.weekday() == Weekday::Sun {
                break;
            }

            date = date.checked_sub_days(Days::new(1)).unwrap();
        }

        date = date.checked_sub_days(Days::new(7)).unwrap();
        let date2 = date.checked_add_days(Days::new(6)).unwrap();

        (date, date2)
    }

    fn fmt_date(&self, date: &DateTime<Utc>) -> String {
        date.format("%Y-%m-%d").to_string()
    }
}
