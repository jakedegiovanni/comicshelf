use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};
use tera::{Tera, Context};

struct ComicShelf {
    tera: Tera
}

#[get("/")]
async fn index(data: web::Data<ComicShelf>) -> impl Responder {
    let mut ctx = Context::new();
    ctx.insert("name", "World");

    let body = data.tera.render("index1.html", &ctx).unwrap();

    HttpResponse::Ok().content_type("text/html; charset=utf-8").body(body)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let tera = match Tera::new("templates/**/*.html") {
        Ok(t) => t,
        Err(e) => {
            println!("Parsing error(s): {}", e);
            std::process::exit(1);
        }
    };

    let d = web::Data::new(ComicShelf{tera});

    HttpServer::new(move || {
        App::new()
            .app_data(d.clone())
            .service(index)
    })
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}
