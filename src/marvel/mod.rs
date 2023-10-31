use anyhow::anyhow;
use chrono::{Datelike, DateTime, Days, Months, Utc, Weekday};
use hyper::{Body, Request};
use tower::Service;

use self::template::DataWrapper;

pub mod auth;
pub mod etag;
pub mod template;

pub trait MarvelService: Service<
    Request<Body>,
    Response = DataWrapper,
    Error = anyhow::Error,
> + Send + Sync + Clone {}

impl<S> MarvelService for S
where
S: Service<
    Request<Body>,
    Response = DataWrapper,
    Error = anyhow::Error,
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

    pub async fn weekly_comics(&self, date: DateTime<Utc>) -> anyhow::Result<DataWrapper> {
        let (date, date2) = self.week_range(date).ok_or(anyhow!("bad date"))?;
        let date = self.fmt_date(&date);
        let date2 = self.fmt_date(&date2);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={},{}&hasDigitalIssue=true&orderBy=issueNumber&limit=100", date, date2);
        let req = Request::get(uri).body(Body::empty())?;

        let mut svc = self.svc.clone();
        svc.call(req).await
    }

    fn week_range(&self, time: DateTime<Utc>) -> Option<(DateTime<Utc>, DateTime<Utc>)> {
        let mut date = time.checked_sub_months(Months::new(3))?;
        loop {
            if date.weekday() == Weekday::Sun {
                break;
            }

            date = date.checked_sub_days(Days::new(1))?;
        }

        date = date.checked_sub_days(Days::new(7))?;
        let date2 = date.checked_add_days(Days::new(6))?;

        Some((date, date2))
    }

    fn fmt_date(&self, date: &DateTime<Utc>) -> String {
        date.format("%Y-%m-%d").to_string()
    }
}
