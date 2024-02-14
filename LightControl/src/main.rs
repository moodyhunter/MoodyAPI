mod fastcon;

use bluer::{adv::Advertisement, Adapter};
use fastcon::light::{BLELight, LightState as FastConLightState};
use futures::StreamExt;
use paho_mqtt as mqtt;
use serde_json::Value;
use std::{collections::BTreeMap, env, error::Error, time::Duration};
use tokio::time::sleep;

const MY_MANUFACTURER_DATA_KEY: u16 = 0xfff0;

const PHONE_KEY: [u8; 4] = [0x33, 0x31, 0x37, 0x33];

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
    let session = bluer::Session::new().await?;
    let adapter = session.default_adapter().await?;
    adapter.set_powered(true).await?;

    let host = env::args()
        .nth(1)
        .unwrap_or_else(|| "mqtt://localhost:1883".to_string());

    println!("Connecting to the MQTT server at '{}'...", host);

    // Create the client. Use a Client ID for a persistent session.
    // A real system should try harder to use a unique ID.
    let create_opts = mqtt::CreateOptionsBuilder::new()
        .server_uri(host)
        .client_id("rust_async_subscribe")
        .finalize();

    let mut cli = mqtt::AsyncClient::new(create_opts)?;
    let mut stream = cli.get_stream(25);

    // Create the connect options, explicitly requesting MQTT v3.x
    let conn_opts = mqtt::ConnectOptionsBuilder::new()
        .keep_alive_interval(Duration::from_secs(30))
        .clean_session(false)
        .finalize();

    // Make the connection to the broker
    cli.connect(conn_opts).await?;
    cli.subscribe("brMesh/3/set", 1).await?;

    // Just loop on incoming messages.
    println!("Waiting for messages...");

    let mut rconn_attempt: usize = 0;

    let mut light = BLELight::new(1, &PHONE_KEY);
    let mut prev_state = FastConLightState::WarmWhite;

    while let Some(msg_opt) = stream.next().await {
        match msg_opt {
            Some(msg) => {
                let payload_str = msg.payload_str();

                let data: Value = serde_json::from_str(&payload_str)?;
                let data = data.as_object().expect("Data is not an object");

                if !data.contains_key("state") {
                    println!("No state key found in the message");
                    continue;
                }

                if data["state"].as_str().unwrap() == "OFF" {
                    light.set_state(FastConLightState::Off);
                } else {
                    if let Some(brightness) = data.get("brightness") {
                        // brightness is 3..255, we scale it to 0..127
                        let brightness = brightness
                            .as_u64()
                            .map(|b| (b - 3) * (127 - 1) / (253 - 3) + 1)
                            .unwrap() as u8;
                        println!("Setting brightness to {}", brightness);
                        light.set_brightness(brightness);
                    }

                    if let Some(color) = data.get("color") {
                        let color = color.as_object().unwrap();
                        light.set_state(FastConLightState::RGB(
                            color["r"].as_u64().unwrap() as u8,
                            color["g"].as_u64().unwrap() as u8,
                            color["b"].as_u64().unwrap() as u8,
                        ));
                    } else if let Some(_) = data.get("color_temp") {
                        light.set_state(FastConLightState::WarmWhite);
                    } else {
                        light.set_state(prev_state)
                    }

                    prev_state = light.get_state();
                }

                let data = light.get_advertisement();
                do_advertise(&adapter, &data).await?;
            }
            _ => {
                // A "None" means we were disconnected. Try to reconnect...
                println!("Lost connection. Attempting reconnect...");
                while let Err(err) = cli.reconnect().await {
                    rconn_attempt += 1;
                    println!("Error reconnecting #{}: {}", rconn_attempt, err);
                    sleep(Duration::from_secs(1)).await;
                }
                println!("Reconnected.");
            }
        }
    }

    Ok(())
}
