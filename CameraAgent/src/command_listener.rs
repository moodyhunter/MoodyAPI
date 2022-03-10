use std::{process::Command, sync::atomic::Ordering};

use tonic::Request;

use crate::{
    camera_api::{
        camera_service_client::CameraServiceClient, Auth, SubscribeCameraStateChangeRequest,
    },
    common::GlobalState,
};

pub async fn listen_for_state_change(state: &Box<GlobalState>) {
    let mut client = CameraServiceClient::new(state.channel.clone());

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
                    Ok(msg) => match msg {
                        Some(s) => {
                            update_camera_status(s.new_state());
                            state.camera_state.store(s.new_state(), Ordering::Relaxed);
                        }
                        None => println!("Received an empty message."),
                    },
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
