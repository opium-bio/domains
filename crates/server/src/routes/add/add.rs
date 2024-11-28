//remove these later
#![allow(unused_imports)] 
#![allow(private_interfaces)] 
//
use rocket::serde::{json::Json, Deserialize, Serialize};
use core::config::config; 

#[derive(Debug, Deserialize)]
struct Input<'r> {
    domain: &'r str,
}
#[derive(Serialize)]
struct Response {
    message: String,
}

#[post("/", format = "json", data = "<user_input>")]
pub fn add_route(user_input: Json<Input>) -> Json<Response> {
    Json(Response {
        message: format!("print test {}", user_input.domain),
    })
}
