use chrono::{Datelike, Days, Months, NaiveDate, Weekday};
use hyper::{Body, Request};
use tower::{BoxError, Service};

use self::template::DataWrapper;

pub mod auth;
pub mod etag;
pub mod template;

pub trait Client:
    Service<Request<Body>, Response = DataWrapper, Error = BoxError> + Send + Sync + Clone
{
}

impl<S> Client for S where
    S: Service<Request<Body>, Response = DataWrapper, Error = BoxError> + Send + Sync + Clone
{
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
    S: Client,
{
    pub async fn weekly_comics(&self, date: NaiveDate) -> Result<DataWrapper, BoxError> {
        let (date, date2) = Marvel::<S>::week_range(date);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={date},{date2}&hasDigitalIssue=true&orderBy=issueNumber&limit=100");
        let req = Request::get(uri).body(Body::empty())?;

        let mut svc = self.svc.clone();
        svc.call(req).await
    }
}
