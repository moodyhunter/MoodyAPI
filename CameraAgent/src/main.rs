mod command_listener;
mod models;

use ini::Ini;
use std::{error::Error, time::Duration};
use tokio::time::sleep;
use tonic::transport::Channel;

use crate::command_listener::keep_alive;

#[tokio::main(flavor = "multi_thread")]
async fn main() -> Result<(), Box<dyn Error>> {
    let conf = Ini::load_from_file("/etc/moodyapi/CameraAgent.ini")?;

    let api_host = conf.general_section().get("Server").unwrap().to_string();
    let client_id = conf.general_section().get("ClientID").unwrap().to_string();

    let channel = Channel::from_shared(api_host.clone())?.connect().await?;

    tokio::select! {
        _ = tokio::spawn(keep_alive(channel.clone(), client_id.clone())) => unreachable!(),
        _ = tokio::spawn(process_command(channel.clone(), client_id.clone())) => unreachable!(),
        _ = tokio::spawn(report_status(channel.clone(), client_id.clone())) => unreachable!()
    }
}

async fn process_command(channel: Channel, client_id: String) {
    loop {
        match command_listener::listen_for_state_change(&channel, &client_id).await {
            Ok(()) => println!("Aborting requests, waiting for retransmission."),
            Err(e) => println!("Error: {}", e),
        };
        sleep(Duration::from_secs(5)).await;
    }
}

async fn report_status(_channel: Channel, _client_id: String) {
    loop {
        sleep(Duration::from_secs(5)).await;
    }
}
