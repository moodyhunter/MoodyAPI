use std::{process::Command, sync::atomic::Ordering};

use tonic::Request;

use crate::camera_api::moody_api_service_client::MoodyApiServiceClient;
use crate::camera_api::{Auth, SubscribeCameraStateChangeRequest};
use crate::common::GlobalState;

pub async fn listen_for_state_change(state: &Box<GlobalState>) {
    let mut client = MoodyApiServiceClient::new(state.channel.clone());

    let request = Request::new(SubscribeCameraStateChangeRequest {
        auth: Some(Auth {
            secret: state.api_secret.clone(),
        }),
    });

    match client.subscribe_camera_state_change(request).await {
        Ok(stream) => {
            let mut resp = stream.into_inner();
            loop {
                match resp.message().await {
                    Ok(None) => println!("Received an empty message."),
                    Ok(Some(s)) => {
                        update_camera_status(s.new_state());
                        state.camera_state.store(s.new_state(), Ordering::Relaxed);
                    }
                    Err(e) => {
                        println!("What? {:?}", e.message());
                        break;
                    }
                }
            }
        }
        Err(e) => println!("something went wrong: {:?}", e.message()),
    }
}

fn update_camera_status(new_status: bool) {
    println!("Updating camera status: {:?}", new_status);

    // TODO: Report exceptions.
    Command::new("sudo")
        .arg("/usr/bin/systemctl")
        .arg(if new_status { "start" } else { "stop" })
        .arg("motion.service")
        .output()
        .expect("failed to execute process");
}
