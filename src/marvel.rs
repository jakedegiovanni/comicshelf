use std::sync::Arc;

use chrono::{DateTime, Datelike, Days, Months, Utc, Weekday};
use hyper::client::HttpConnector;
use hyper::{Body, Request, Response};
use hyper_tls::HttpsConnector;
use tokio::sync::Mutex;
use tower::util::BoxCloneService;
use tower::{Service, ServiceBuilder};

use crate::template::DataWrapper;

type HyperService = BoxCloneService<Request<Body>, Response<Body>, hyper::Error>;

pub struct Marvel {
    svc: Arc<Mutex<HyperService>>,
}

impl Marvel {
    pub fn new(client: &hyper::Client<HttpsConnector<HttpConnector>, hyper::Body>) -> Self {
        let svc = Marvel::svc(client);
        let svc = Arc::new(Mutex::new(svc));
        Marvel { svc }
    }

    fn uri_middleware(req: Request<Body>) -> Request<Body> {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        up.authority = Some(hyper::http::uri::Authority::from_static(
            "gateway.marvel.com",
        ));
        up.scheme = Some(hyper::http::uri::Scheme::HTTPS);

        p.uri = hyper::Uri::from_parts(up).unwrap();
        Request::from_parts(p, b)
    }

    fn auth_middleware(req: Request<Body>) -> Request<Body> {
        let (mut p, b) = req.into_parts();

        let mut up = p.uri.into_parts();
        let pq = up.path_and_query.unwrap();
        let path = pq.path();
        let q = pq.query().unwrap_or("");

        let path = {
            if !path.contains("/v1/public") {
                format!("/v1/public{}", path)
            } else {
                path.to_string()
            }
        };

        let pub_key = include_str!("../pub.txt"); // todo: formalize, this is janky
        let priv_key = include_str!("../priv.txt");

        let ts = Utc::now().timestamp_millis();
        let hash = format!(
            "{:x}",
            md5::compute(format!("{}{}{}", ts, priv_key, pub_key))
        );
        let query = format!("apikey={}&ts={}&hash={}", pub_key, ts, hash);
        let query = {
            if q.is_empty() {
                format!("?{}", query)
            } else {
                format!("{}&{}", q, query)
            }
        };
        let query = format!("{}?{}", path, query);
        up.path_and_query = Some(hyper::http::uri::PathAndQuery::try_from(query).unwrap());

        p.uri = hyper::Uri::from_parts(up).unwrap();
        Request::from_parts(p, b)
    }

    fn svc(client: &hyper::Client<HttpsConnector<HttpConnector>, Body>) -> HyperService {
        let svc = ServiceBuilder::new()
            .map_request(Marvel::uri_middleware)
            .map_request(Marvel::auth_middleware)
            .service(client.clone());

        BoxCloneService::new(svc)
    }

    pub async fn weekly_comics(&self) -> DataWrapper {
        let (date, date2) = self.week_range(Utc::now());
        let date = self.fmt_date(&date);
        let date2 = self.fmt_date(&date2);

        let uri = format!("/comics?format=comic&formatType=comic&noVariants=true&dateRange={},{}&hasDigitalIssue=true&orderBy=issueNumber&limit=100", date, date2);
        let req = Request::get(uri).body(hyper::Body::empty()).unwrap();

        let result = {
            let mut svc = self.svc.lock().await; // todo - future bottleneck
            svc.call(req).await.unwrap().into_body()
        };

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
