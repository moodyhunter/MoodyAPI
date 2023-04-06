mod fastcon;

use std::{collections::BTreeMap, time::Duration};

use tokio::time::sleep;

use bluer::{adv::Advertisement, Adapter, AdapterEvent, Address, DeviceEvent, DeviceProperty};
use fastcon::{
    broadcast_parser::{parse_ble_broadcast, BroadcastType},
    light::{BLELight, LightState},
};
use futures::{stream::SelectAll, StreamExt};

use crate::fastcon::DEFAULT_PHONE_KEY;

const MY_ADDRESS: Address = Address::new([0x11, 0x22, 0x33, 0x44, 0x55, 0x66]);
const MY_MANUFACTURER_DATA_KEY: u16 = 0xfff0;

fn print_bytes(name: &str, bytes: &[u8]) {
    print!("{}: ", name);
    for b in bytes {
        print!("{:02x} ", b);
    }
    println!();
}

async fn handle_light_event(
    _adapter: &Adapter,
    ev: DeviceEvent,
) -> Result<(), Box<dyn std::error::Error>> {
    let DeviceEvent::PropertyChanged(prop) = ev;

    match prop {
        DeviceProperty::Rssi(_) => Ok(()), // ignore RSSI, we don't care
        DeviceProperty::ManufacturerData(dmap) => {
            let Some(data) = dmap.get(&MY_MANUFACTURER_DATA_KEY) else {
                println!("Invalid light bulb manufacturer data");
                return Ok(());
            };

            // dump the data
            println!("    Manufacturer data: {:?}", data);

            // parse the data
            let Some(x) = parse_ble_broadcast(data, &DEFAULT_PHONE_KEY) else {
                println!("Invalid light bulb manufacturer data");
                return Ok(());
            };

            match x {
                BroadcastType::HeartBeat(heartbeat) => {
                    println!("Heartbeat: {:?}", heartbeat);
                    // single_on_off_command(heartbeat.short_addr, true);
                }
                BroadcastType::TimerUploadResponse => {
                    println!("Timer upload response");
                }
                BroadcastType::DeviceAnnouncement(device) => {
                    println!("Device announcement: {:?}", device);
                }
            }

            Ok(())
        }
        other => {
            println!("    Unhandled property: {:?}", other);
            Ok(())
        }
    }
}

async fn do_advertise(adapter: &Adapter, data: &Vec<u8>) -> Result<(), Box<dyn std::error::Error>> {
    print_bytes("Advertisement Data", data);
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
    println!("{:?}", &le_advertisement);
    let handle = adapter.advertise(le_advertisement).await?;

    println!("Removing advertisement");
    drop(handle);
    sleep(Duration::from_millis(500)).await;

    Ok(())
}

async fn do_poll_device_events(adapter: &Adapter) -> Result<(), Box<dyn std::error::Error>> {
    let mut device_events = adapter.discover_devices().await?;
    let mut light_change_events = SelectAll::new();

    loop {
        tokio::select! {
            Some(de) = device_events.next() => {
                match de {
                    AdapterEvent::DeviceAdded(addr) => {
                        if addr != MY_ADDRESS {
                            continue;
                        }

                        println!("Device added: {addr}");
                        let device = adapter.device(MY_ADDRESS).expect("Error getting device");
                        let current_data = device.manufacturer_data().await?;
                        if let Some(data) = current_data {
                            handle_light_event(&adapter, DeviceEvent::PropertyChanged(DeviceProperty::ManufacturerData(data))).await.expect("Error handling light event");
                        }

                        light_change_events.push(device.events().await?);
                    }
                    AdapterEvent::DeviceRemoved(addr) => {
                        if addr != MY_ADDRESS {
                            continue;
                        }
                        device_events = adapter.discover_devices().await?;
                        println!("Light {addr} fell off the bus");
                    }
                    _ => (),
                }
            },
            Some(le) = light_change_events.next() => handle_light_event(&adapter, le).await.expect("Error handling light event"),
            else => {
                println!("No more events");
                std::process::exit(1);
            }
        }
    }
}

#[tokio::main(flavor = "current_thread")]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let session = bluer::Session::new().await?;
    let adapter = session.default_adapter().await?;
    adapter.set_powered(true).await?;

    let phone_key = vec![0x38, 0x35, 0x36, 0x30];

    let mut light = BLELight::new(1, &phone_key);
    light.set_state(LightState::WarmWhite);
    light.set_brightness(127);
    let data = light.get_advertisement();

    do_advertise(&adapter, &data).await?;

    do_poll_device_events(&adapter).await?;

    Ok(())
}
