use once_cell::sync::OnceCell;
use std::process::Command;
use std::time::Duration;
use tokio::time::sleep;
use tonic::{transport::Channel, Request};

use crate::models::{
    common::Auth,
    moody_api::{
        moody_api_service_client::MoodyApiServiceClient, CameraState, KeepAliveRequest,
        SubscribeCameraStateChangeRequest, UpdateCameraStateRequest,
    },
    notifications::{Notification, SendRequest},
};

pub static CLIENT_ID: OnceCell<String> = OnceCell::new();

fn get_client_id() -> String {
    CLIENT_ID.get().unwrap().clone()
}

pub async fn keep_alive(mut client: MoodyApiServiceClient<Channel>) {
    loop {
        let request = Request::new(KeepAliveRequest {
            auth: Some(Auth {
                client_uuid: get_client_id(),
            }),
        });

        match client.keep_alive(request).await {
            Ok(stream) => {
                let mut resp = stream.into_inner();
                loop {
                    match resp.message().await {
                        Ok(None) => println!("Received an empty message."),
                        Ok(Some(keep_alive_resp)) => {
                            println!("Received a keep alive response: {:?}", keep_alive_resp);
                        }
                        Err(e) => {
                            println!("Keepalive inner error: {:?}", e.message());
                            sleep(Duration::from_secs(5)).await;
                            break;
                        }
                    }
                }
            }
            Err(e) => println!("Keepalive error: {}", e.message()),
        };
        sleep(Duration::from_secs(5)).await;
    }
}

pub async fn listen_for_state_change(mut client: MoodyApiServiceClient<Channel>) -> ! {
    let mut error_message_sent: bool;
    loop {
        error_message_sent = false;
        let request = Request::new(SubscribeCameraStateChangeRequest {
            auth: Some(Auth {
                client_uuid: get_client_id(),
            }),
        });

        match client.subscribe_camera_control_signal(request).await {
            Ok(stream) => {
                let mut resp = stream.into_inner();
                loop {
                    match resp.message().await {
                        Ok(None) => println!("Received an empty message."),
                        Ok(Some(s)) => {
                            let control_status = start_stop_camera_service(s.state);
                            if !control_status && !error_message_sent {
                                send_notification(
                                    client.clone(),
                                    8,
                                    "Failed to start/stop service".to_string(),
                                    "Why?".to_string(),
                                )
                                .await;
                                error_message_sent = true;
                            }
                            report_camera_status_internal(client.clone(), get_camera_status())
                                .await;
                        }
                        Err(e) => {
                            println!("Listener inner error? {:?}", e.message());
                            sleep(Duration::from_secs(5)).await;
                            break;
                        }
                    }
                }
            }
            Err(e) => println!("Listener error: {:?}", e.message()),
        };

        sleep(Duration::from_secs(20)).await;
    }
}

pub async fn report_camera_status(mut _client: MoodyApiServiceClient<Channel>) {
    loop {
        sleep(Duration::from_secs(5)).await;
    }
}

async fn report_camera_status_internal(mut client: MoodyApiServiceClient<Channel>, started: bool) {
    let request = Request::new(UpdateCameraStateRequest {
        auth: Some(Auth {
            client_uuid: get_client_id(),
        }),
        state: Some(CameraState { state: started }),
    });

    client
        .report_camera_state(request)
        .await
        .expect("Failed to send notification.");
}

fn start_stop_camera_service(new_status: bool) -> bool {
    println!("Updating camera status: {:?}", new_status);

    if let Ok(e) = Command::new("sudo")
        .arg("/usr/bin/systemctl")
        .arg(if new_status { "start" } else { "stop" })
        .arg("motion.service")
        .status()
    {
        return e.success();
    }

    false
}

fn get_camera_status() -> bool {
    if let Ok(e) = Command::new("/usr/bin/systemctl")
        .arg("status")
        .arg("motion.service")
        .output()
    {
        return e.status.success() && String::from_utf8_lossy(&e.stdout).contains("Active: active");
    }
    false
}

async fn send_notification(
    mut client: MoodyApiServiceClient<Channel>,
    n_channel: i64,
    n_title: String,
    n_content: String,
) -> () {
    let n = Notification {
        title: n_title,
        content: n_content,
        channel_id: n_channel,
        ..Default::default()
    };

    client
        .send_notification(Request::new(SendRequest {
            auth: Some(Auth {
                client_uuid: get_client_id(),
            }),
            notification: Some(n),
        }))
        .await
        .expect("Failed to send notification.");
}
