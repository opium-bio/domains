#[macro_use]
extern crate rocket;

mod routes;

use rocket::{Build, Rocket};

pub async fn serve() -> Rocket<Build> {
    rocket::build()
        .mount("/add", routes![routes::add::add::add_route])  
}

#[launch]
async fn rocket() -> Rocket<Build> {
    serve().await
}
