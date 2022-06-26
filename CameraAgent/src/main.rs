mod command_listener;
mod models;

use ini::Ini;
use std::error::Error;
use tonic::transport::Channel;

use crate::command_listener::keep_alive;
use crate::command_listener::listen_for_state_change;
use crate::command_listener::report_camera_status;
use crate::command_listener::CHANNEL;
use crate::command_listener::CLIENT_ID;

#[tokio::main(flavor = "multi_thread")]
async fn main() -> Result<(), Box<dyn Error>> {
    let conf = Ini::load_from_file("/etc/moodyapi/CameraAgent.ini")?;

    let api_host = conf.general_section().get("Server").unwrap().to_string();
    CLIENT_ID
        .set(conf.general_section().get("ClientID").unwrap().to_string())
        .expect("Failed to set CLIENT_ID");

    CHANNEL
        .set(Channel::from_shared(api_host.clone())?.connect().await?)
        .expect("Failed to connect to server");

    tokio::select! {
        _ = keep_alive() => {},
        _ = listen_for_state_change() => {},
        _ = report_camera_status() => {},
    }

    unreachable!();
}
