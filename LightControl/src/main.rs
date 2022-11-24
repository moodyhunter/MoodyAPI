mod models;

use btleplug::{
    api::{
        Central, CharPropFlags, Characteristic, Manager as _, Peripheral as _, ScanFilter,
        WriteType,
    },
    platform::{Adapter, Manager, Peripheral},
};
use core::panic;
use ini::Ini;
use models::{
    light::{light_state::Mode, LightState},
    moody_api::moody_api_service_client::MoodyApiServiceClient,
};
use platform_dirs::AppDirs;
use std::{error::Error, str::FromStr, time::Duration};
use tokio::time::sleep;
use tonic::{transport::Channel, Request};
use uuid::Uuid;

use crate::models::{common::Auth, light::SubscribeLightRequest};

const FFD5: &str = "0000ffd5-0000-1000-8000-00805f9b34fb";
const FFD9: &str = "0000ffd9-0000-1000-8000-00805f9b34fb";

async fn get_light(central: &Adapter) -> Peripheral {
    for p in central.peripherals().await.unwrap() {
        if p.properties()
            .await
            .unwrap()
            .unwrap()
            .local_name
            .iter()
            .any(|name| name.contains("Trione"))
        {
            return p;
        }
    }
    panic!("Could not find light");
}

async fn send_command(peripheral: &Peripheral, data: &Vec<u8>) -> Result<(), btleplug::Error> {
    let cc = Characteristic {
        service_uuid: Uuid::from_str(&FFD5).unwrap(),
        uuid: Uuid::from_str(&FFD9).unwrap(),
        properties: CharPropFlags::WRITE_WITHOUT_RESPONSE,
    };
    return peripheral
        .write(&cc, &data, WriteType::WithoutResponse)
        .await;
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let adapter = Manager::new()
        .await
        .unwrap()
        .adapters()
        .await
        .expect("NO ADAPTERS FOUND")
        .into_iter()
        .nth(0)
        .expect("NO ADAPTERS FOUND");

    adapter.start_scan(ScanFilter::default()).await?;
    sleep(Duration::from_secs(5)).await;

    let light = get_light(&adapter).await;
    light.connect().await?;
    light.discover_services().await?;
    light.characteristics().iter().for_each(|c| {
        println!(
            "Found characteristic: {} {} {:?}",
            c.service_uuid, c.uuid, c.properties
        )
    });

    let dirs = AppDirs::new(Some("moodyapi"), false).unwrap();
    let ini_path = dirs.config_dir.as_path().join("LightDaemon.ini");

    let conf = if let Ok(ini_file) = Ini::load_from_file(ini_path.as_path()) {
        ini_file
    } else if let Ok(ini_file) = Ini::load_from_file("/etc/moodyapi/LightDaemon.ini") {
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

    println!("Starting in notification client mode, listening for new notifications...");

    loop {
        let mut client = MoodyApiServiceClient::new(grpc_channel.clone());

        let request = Request::new(SubscribeLightRequest {
            auth: Some(Auth {
                client_uuid: client_id.clone(),
            }),
            ..Default::default()
        });

        match client.subscribe_light_state_change(request).await {
            Ok(stream) => {
                let mut resp_stream = stream.into_inner();
                loop {
                    match resp_stream.message().await {
                        Ok(Some(l)) => {
                            println!("Received LightState: {:?}", l);
                            match send_light_command(&light, l).await {
                                Ok(_) => {}
                                Err(e) => println!("Failed to send light command: {}", e),
                            }
                        }
                        e => {
                            println!("something went wrong: {:?}", e);
                            sleep(Duration::from_secs(1)).await;
                            break;
                        }
                    }
                }
            }
            e => {
                println!("something went wrong: {:?}", e);
                sleep(Duration::from_secs(1)).await;
            }
        }

        sleep(Duration::from_secs(10)).await;
    }
}

async fn send_light_command(
    light_dev: &Peripheral,
    lightstate: LightState,
) -> Result<(), btleplug::Error> {
    let on_command = vec![0xcc, 0x23, 0x33];
    let off_command = vec![0xcc, 0x24, 0x33];

    if lightstate.on {
        send_command(&light_dev, &on_command).await?;
        let brightness = lightstate.brightness as u8;
        match lightstate.mode {
            Some(m) => match m {
                Mode::Colored(color) => {
                    let command = vec![
                        0x56,
                        color.red as u8,
                        color.green as u8,
                        color.blue as u8,
                        brightness,
                        0xf0,
                        0xaa,
                    ];
                    send_command(&light_dev, &command).await?
                }
                Mode::Warmwhite(_) => {
                    let command = vec![0x56, 0x00, 0x00, 0x00, brightness, 0x0F, 0xaa];
                    send_command(&light_dev, &command).await?
                }
            },
            None => {
                // No mode set, default to warm white
                let command = vec![0x56, 0x00, 0x00, 0x00, brightness, 0x0F, 0xaa];
                send_command(&light_dev, &command).await?
            }
        }
    } else {
        send_command(&light_dev, &off_command).await?
    }

    Ok(())
}
