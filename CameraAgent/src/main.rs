mod command_listener;
mod models;

use ini::Ini;
use std::error::Error;

use crate::command_listener::keep_alive;
use crate::command_listener::listen_for_state_change;
use crate::command_listener::report_camera_status;
use crate::command_listener::CLIENT_ID;
use crate::models::moody_api::moody_api_service_client::MoodyApiServiceClient;

static CONFIG_FILE_PATH: &str = "/etc/moodyapi/CameraAgent.ini";

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let conf = Ini::load_from_file(CONFIG_FILE_PATH)?;

    let api_host = conf.general_section().get("Server").unwrap().to_string();
    CLIENT_ID
        .set(conf.general_section().get("ClientID").unwrap().to_string())
        .expect("Failed to set CLIENT_ID");

    let client = MoodyApiServiceClient::connect(api_host.clone()).await?;

    tokio::select! {
        _ = keep_alive(client.clone()) => {},
        _ = listen_for_state_change(client.clone()) => {},
        _ = report_camera_status(client.clone()) => {},
    }

    unreachable!();
}
