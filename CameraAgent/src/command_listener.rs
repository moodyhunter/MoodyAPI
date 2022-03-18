use std::error::Error;
use std::process::Command;
use tonic::transport::Channel;
use tonic::Request;

use crate::camera_api::moody_api_service_client::MoodyApiServiceClient;
use crate::camera_api::{Auth, SubscribeCameraStateChangeRequest};

pub async fn listen_for_state_change(
    api_host: &String,
    client_id: &String,
) -> Result<(), Box<dyn Error>> {
    let channel = Channel::from_shared(api_host.clone())?.connect().await?;
    let mut client = MoodyApiServiceClient::new(channel.clone());

    let request = Request::new(SubscribeCameraStateChangeRequest {
        auth: Some(Auth {
            client_id: client_id.clone(),
        }),
    });

    match client.subscribe_camera_state_change(request).await {
        Ok(stream) => {
            let mut resp = stream.into_inner();
            loop {
                match resp.message().await {
                    Ok(None) => println!("Received an empty message."),
                    Ok(Some(s)) => {
                        update_camera_status(s.state());
                    }
                    Err(e) => {
                        println!("What? {:?}", e.message());
                        break;
                    }
                }
            }
        }
        Err(e) => println!("something went wrong: {:?}", e.message()),
    };

    Ok(())
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
