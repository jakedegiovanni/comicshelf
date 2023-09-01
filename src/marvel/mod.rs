use chrono::{DateTime, Datelike, Days, Months, Utc, Weekday};
use hyper::client::HttpConnector;
use hyper::{Body, Client, Request, Response};
use hyper_tls::HttpsConnector;
use tokio::sync::Mutex;
use tower::util::BoxCloneService;
use tower::{Service, ServiceBuilder};

use crate::middleware::uri::UriMiddlewareLayer;
use auth::AuthMiddlewareLayer;
use template::DataWrapper;

mod auth;
mod template;

type HyperService = BoxCloneService<Request<Body>, Response<Body>, hyper::Error>;

pub struct Marvel {
    svc: Mutex<HyperService>,
}

impl Marvel {
    pub fn new(client: &Client<HttpsConnector<HttpConnector>, Body>) -> Self {
        let svc = Marvel::svc(client);
        let svc = Mutex::new(svc);
        Marvel { svc }
    }

    fn svc(client: &Client<HttpsConnector<HttpConnector>, Body>) -> HyperService {
        let svc = ServiceBuilder::new()
            .layer(UriMiddlewareLayer::new("gateway.marvel.com", "https"))
            .layer(AuthMiddlewareLayer::new(
                include_str!("../../pub.txt"), // todo: formalize, this is janky
                include_str!("../../priv.txt"),
            ))
            .service(client.clone());

        BoxCloneService::new(svc)
    }

    pub async fn weekly_comics(&self) -> DataWrapper {
        let (date, date2) = self.week_range(Utc::now());
        let date = self.fmt_date(&date);
        let date2 = self.fmt_date(&date2);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={},{}&hasDigitalIssue=true&orderBy=issueNumber&limit=100", date, date2);
        let req = Request::get(uri).body(Body::empty()).unwrap();

        // todo - anyway to avoid the lock here whilst still keeping the Marvel struct shared?
        let mut svc = { self.svc.lock().await.clone() };
        let result = svc.call(req).await.unwrap().into_body();

        let result: DataWrapper = serde_json::from_slice(
            hyper::body::to_bytes(result)
                .await
                .unwrap()
                .iter()
                .as_slice(),
        )
        .unwrap();

        result
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
