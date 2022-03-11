mod notification_api;

use ini::Ini;
use notify_rust::{Hint, Notification as Notify};
use platform_dirs::AppDirs;
use prost_types::Timestamp;
use std::process::exit;
use std::thread::sleep;
use std::time::{Duration, SystemTime};
use std::{env, fs};
use tonic::{transport::Channel, Request};

use crate::notification_api::{
    moody_api_service_client::MoodyApiServiceClient, Auth, Notification, SendNotificationRequest,
    SubscribeNotificationsRequest,
};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let dirs = AppDirs::new(Some("Moody"), false).unwrap();

    let ini_path = dirs.config_dir.as_path().join("NotifyClient.ini");
    let conf = if let Ok(ini_file) = Ini::load_from_file(ini_path.as_path()) {
        ini_file
    } else {
        println!("Created configuration at: {}", ini_path.display());
        fs::create_dir_all(dirs.config_dir).expect("Failed to create configuration directory.");

        let mut ini = Ini::new();
        ini.with_section(Some("APIServer"))
            .set("Address", "apiserver.example.com")
            .set("TLS", "true")
            .set("Secret", "f8555071-fb14-402a-95c8-70265ec7c965");

        ini.write_to_file(ini_path.as_path())
            .expect("Failed to write file.");
        Ini::load_from_file(ini_path).unwrap()
    };

    let section = conf.section(Some("APIServer")).unwrap();

    let api_host = section.get("Address").unwrap().to_string();
    let api_secret = section.get("Secret").unwrap().to_string();

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

        send_notification(notification_channel, title, message, &channel, &api_secret).await;
    } else {
        println!("Starting in notification client mode, listening for new notifications...");
        listen_notification(channel, api_secret).await
    }

    Ok(())
}

async fn listen_notification(channel: Channel, api_secret: String) -> ! {
    loop {
        let mut client = MoodyApiServiceClient::new(channel.clone());

        let request = Request::new(SubscribeNotificationsRequest {
            auth: Some(Auth {
                secret: api_secret.clone(),
            }),
            channel_id: 1,
        });

        match client.subscribe_notifications(request).await {
            Err(e) => println!("something went wrong: {}", e),
            Ok(stream) => {
                let mut resp_stream = stream.into_inner();
                loop {
                    match resp_stream.message().await {
                        Err(e) => println!("something went wrong: {}", &e),
                        Ok(None) => println!("expect a notification object"),
                        Ok(Some(n)) => {
                            println!("Received Notification: {:?}", n);
                            Notify::new()
                                .summary(&n.title)
                                .body(&n.message)
                                .icon(&n.icon)
                                .appname("Notify Client")
                                .hint(Hint::Resident(true))
                                .show()
                                .unwrap();
                        }
                    }
                    sleep(Duration::from_secs(2));
                }
            }
        }

        sleep(Duration::from_secs(10));
    }
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
                secret: api_secret.clone(),
            }),
            notification: Some(n),
        }))
        .await
        .expect("Failed to send notification.");
}
