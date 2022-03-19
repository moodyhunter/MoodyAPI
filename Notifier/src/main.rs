mod notification_api;

use ini::Ini;
use notify_rust::{Hint, Notification as Notify};
use platform_dirs::AppDirs;
use prost_types::Timestamp;
use std::{
    env,
    process::exit,
    time::{Duration, SystemTime},
};
use tokio::time::sleep;
use tonic::{transport::Channel, Request};

use crate::notification_api::{
    moody_api_service_client::MoodyApiServiceClient, Auth, Notification, SendNotificationRequest,
    SubscribeNotificationsRequest,
};

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

    let channel = Channel::from_shared(api_host.clone())?
        .connect()
        .await
        .expect("Can't create a channel");

    if env::args().len() >= 2 {
        if env::args().len() < 3 {
            println!("You should provide 2 or 3 parameters:");
            println!("  title message [channel]");
            exit(1);
        }

        println!("Sending notifications...");
        let title = env::args().nth(1).unwrap().to_string();
        let message = env::args().nth(2).unwrap().to_string();
        let notification_channel = 1;

        send_notification(notification_channel, title, message, &channel, &client_id).await;
    } else {
        println!("Starting in notification client mode, listening for new notifications...");
        listen_notification(channel, client_id).await
    }

    Ok(())
}

async fn listen_notification(channel: Channel, api_secret: String) -> ! {
    loop {
        let mut client = MoodyApiServiceClient::new(channel.clone());

        let request = Request::new(SubscribeNotificationsRequest {
            auth: Some(Auth {
                client_id: api_secret.clone(),
            }),
            channel_id: 1,
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
        .body(&n.message)
        .icon(&n.icon)
        .appname("Notify Client")
        .hint(Hint::Resident(true))
        .show()
        .unwrap();
}

async fn send_notification(
    notification_channel: i32,
    title: String,
    message: String,
    channel: &Channel,
    api_secret: &String,
) {
    let n = Notification {
        time: Some(Timestamp::from(SystemTime::now())),
        channel_id: notification_channel,
        title,
        message,
        icon: "invalid".to_string(),
    };
    let mut client = MoodyApiServiceClient::new(channel.clone());
    client
        .send_notification(Request::new(SendNotificationRequest {
            auth: Some(Auth {
                client_id: api_secret.clone(),
            }),
            notification: Some(n),
        }))
        .await
        .expect("Failed to send notification.");
}
