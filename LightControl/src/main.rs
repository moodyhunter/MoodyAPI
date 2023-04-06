mod fastcon;
mod models;

use crate::models::{common::Auth, light::SubscribeLightRequest};

use bluer::{adv::Advertisement, Adapter};
use core::panic;
use fastcon::light::{BLELight, LightState as FastConLightState};
use ini::Ini;
use models::{
    light::{light_state::Mode, LightState},
    moody_api::moody_api_service_client::MoodyApiServiceClient,
};
use platform_dirs::AppDirs;
use std::{collections::BTreeMap, error::Error, time::Duration};
use tokio::time::sleep;
use tonic::{transport::Channel, Request};

const MY_MANUFACTURER_DATA_KEY: u16 = 0xfff0;

async fn do_advertise(adapter: &Adapter, data: &Vec<u8>) -> Result<(), Box<dyn std::error::Error>> {
    let mut my_data: BTreeMap<u16, Vec<u8>> = BTreeMap::new();
    my_data.insert(MY_MANUFACTURER_DATA_KEY, data.clone());

    let le_advertisement = Advertisement {
        advertisement_type: bluer::adv::Type::Peripheral,
        manufacturer_data: my_data,
        min_interval: Some(Duration::from_millis(100)),
        max_interval: Some(Duration::from_millis(200)),
        duration: Some(Duration::from_secs(100)),
        discoverable: Some(true),
        tx_power: Some(20),
        ..Default::default()
    };

    let handle = adapter.advertise(le_advertisement).await?;
    drop(handle);
    sleep(Duration::from_millis(500)).await;

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
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

    let session = bluer::Session::new().await?;
    let adapter = session.default_adapter().await?;
    adapter.set_powered(true).await?;

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
                            match send_light_command(&adapter, l).await {
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

const PHONE_KEY: [u8; 4] = [0x38, 0x35, 0x36, 0x30];

async fn send_light_command(
    adapter: &Adapter,
    lightstate: LightState,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut light = BLELight::new(1, &PHONE_KEY);

    if lightstate.on {
        light.set_brightness(lightstate.brightness as u8);
        match lightstate.mode {
            Some(Mode::Warmwhite(_)) => {
                light.set_state(FastConLightState::WarmWhite);
            }
            Some(Mode::Colored(color)) => {
                light.set_state(FastConLightState::RGB(
                    color.red as u8,
                    color.green as u8,
                    color.blue as u8,
                ));
            }
            None => todo!(),
        }
    } else {
        light.set_state(FastConLightState::Off);
    }

    let data = light.get_advertisement();
    do_advertise(&adapter, &data).await?;
    Ok(())
}
