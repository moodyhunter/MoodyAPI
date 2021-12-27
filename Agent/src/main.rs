use common::GlobalState;
use std::{
    sync::{atomic::AtomicBool, Arc},
    thread::sleep,
    time::Duration,
};
use tonic::transport::Channel;

mod camera_api;
mod command_listener;
mod common;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let api_host = std::env::args().nth(1).expect("no api host given");
    let api_secret = std::env::args().nth(2).expect("no api secret given");

    loop {
        let channel = Channel::from_shared(api_host.clone())?
            .connect()
            .await
            .expect("Can't create a channel");

        let state = Arc::new(GlobalState {
            api_secret: api_secret.clone(),
            camera_state: AtomicBool::new(false),
            channel,
        });

        command_listener::listen_for_state_change(&state).await;

        println!("Hey, server stopped responding, we are reconnecting.");
        sleep(Duration::from_secs(5));
    }
}
