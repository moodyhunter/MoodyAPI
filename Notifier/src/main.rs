mod models;
use models::{
    common::Auth,
    moody_api::moody_api_service_client::MoodyApiServiceClient,
    notifications::{Notification, SendRequest, SubscribeRequest},
};

use ini::Ini;
use notify_rust::{Hint, Notification as Notify};
use platform_dirs::AppDirs;
use std::{env, process::exit, time::Duration};
use tokio::time::sleep;
use tonic::{transport::Channel, Request};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let dirs = AppDirs::new(Some("moodyapi"), false).unwrap();
    let ini_path = dirs.config_dir.as_path().join("Notifier.ini");

    let conf = if let Ok(ini_file) = Ini::load_from_file(ini_path.as_path()) {
        ini_file
    } else if let Ok(ini_file) = Ini::load_from_file("/etc/moodyapi/Notifier.ini") {
        ini_file
    } else {
        panic!("Failed to locate configurations.");
    };

    let api_host = conf.general_section().get("Server").unwrap().to_string();
    let client_id = conf.general_section().get("ClientID").unwrap().to_string();

    let grpc_channel = Channel::from_shared(api_host.clone())?
        .connect()
        .await
        .expect("Can't create a channel");

    if env::args().len() >= 2 {
        if env::args().len() < 3 {
            println!("You should provide 3 parameters:");
            println!("  ChannelID Title Message");
            exit(1);
        }

        println!("Sending notifications...");

        // notification channel is an integer
        let channel_str = env::args().nth(1).unwrap().to_string();
        let n_title = env::args().nth(2).unwrap().to_string();
        let n_content = env::args().nth(3).unwrap().to_string();

        let n_channel: i64;
        let is_private: bool = channel_str.ends_with("p");
        if channel_str.ends_with("p") {
            n_channel = channel_str[..channel_str.len() - 1].parse().unwrap();
        } else {
            n_channel = channel_str.parse().unwrap();
        }

        send_notification(
            n_channel,
            n_title,
            n_content,
            &grpc_channel,
            &client_id,
            is_private,
        )
        .await;
    } else {
        println!("Starting in notification client mode, listening for new notifications...");
        listen_notification(grpc_channel, client_id).await
    }

    Ok(())
}

async fn listen_notification(channel: Channel, api_secret: String) -> ! {
    loop {
        let mut client = MoodyApiServiceClient::new(channel.clone());

        let request = Request::new(SubscribeRequest {
            auth: Some(Auth {
                client_uuid: api_secret.clone(),
            }),
            ..Default::default()
        });

        match client.subscribe_notifications(request).await {
            Err(e) => println!("something went wrong: {}", e),
            Ok(stream) => {
                let mut resp_stream = stream.into_inner();
                loop {
                    match resp_stream.message().await {
                        Ok(None) => println!("expect a notification object"),
                        Ok(Some(n)) => {
                            println!("Received Notification: {:?}", n);
                            display_notification(n);
                        }
                        Err(e) => {
                            println!("something went wrong: {}", &e);
                            break;
                        }
                    }
                    sleep(Duration::from_secs(2)).await;
                }
            }
        }

        sleep(Duration::from_secs(10)).await;
    }
}

fn display_notification(n: Notification) {
    Notify::new()
        .summary(&n.title)
        .body(&n.content)
        .icon(&n.icon)
        .appname("Notify Client")
        .hint(Hint::Resident(true))
        .show()
        .unwrap();
}

async fn send_notification(
    n_channel: i64,
    n_title: String,
    n_content: String,
    channel: &Channel,
    api_secret: &String,
    is_private: bool,
) {
    let n = Notification {
        title: n_title,
        content: n_content,
        channel_id: n_channel,
        private: is_private,
        ..Default::default()
    };
    let mut client = MoodyApiServiceClient::new(channel.clone());
    client
        .send_notification(Request::new(SendRequest {
            auth: Some(Auth {
                client_uuid: api_secret.clone(),
            }),
            notification: Some(n),
        }))
        .await
        .expect("Failed to send notification.");
}
