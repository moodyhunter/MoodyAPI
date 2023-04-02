mod fastcon;

use std::{collections::BTreeMap, time::Duration};
use tokio::{
    io::{AsyncBufReadExt, BufReader},
    time::{sleep, Interval},
};

use bluer::{adv::Advertisement, Adapter, AdapterEvent, Address, DeviceEvent, DeviceProperty};
use fastcon::{
    broadcast_parser::{parse_ble_broadcast, BroadcastType},
    command_wrapper::DEFAULT_PHONE_KEY,
};
use futures::{stream::SelectAll, StreamExt};

use crate::fastcon::command_wrapper::single_on_off_command;

const MY_ADDRESS: Address = Address::new([0x11, 0x22, 0x33, 0x44, 0x55, 0x66]);
const MY_MANUFACTURER_DATA_KEY: u16 = 0xfff0;

async fn query_all_device_properties(adapter: &Adapter, addr: Address) -> bluer::Result<()> {
    let device = adapter.device(addr)?;
    let props = device.all_properties().await?;
    for prop in props {
        println!("    {:?}", &prop);
    }
    Ok(())
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

            // parse the data
            let Some(x) = parse_ble_broadcast(data, &DEFAULT_PHONE_KEY) else {
                println!("Invalid light bulb manufacturer data");
                return Ok(());
            };

            match x {
                BroadcastType::HeartBeat(heartbeat) => {
                    println!("Heartbeat: {:?}", heartbeat);
                    single_on_off_command(heartbeat.short_addr, true);
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

#[tokio::main(flavor = "current_thread")]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let session = bluer::Session::new().await?;
    let adapter = session.default_adapter().await?;
    println!(
        "Discovering devices using Bluetooth adapter {}",
        adapter.name()
    );
    adapter.set_powered(true).await?;

    println!(
        "Advertising on Bluetooth adapter {} with address {}",
        adapter.name(),
        adapter.address().await?
    );

    let mut my_data: BTreeMap<u16, Vec<u8>> = BTreeMap::new();
    my_data.insert(
        0xFFF0,
        vec![
            0x6d, 0xb6, 0x43, 0x68, 0x93, 0x1d, 0x5d, 0x13, 0x07, 0x05, 0xa8, 0xaa, 0x87, 0x83,
            0x91, 0xaf, 0xdd, 0x07, 0xbc, 0xf7, 0xcb, 0x6b, 0x67, 0x06,
        ],
    );

    let le_advertisement = Advertisement {
        advertisement_type: bluer::adv::Type::Peripheral,
        manufacturer_data: my_data,
        // advertisting_data: my_data,
        min_interval: Some(Duration::from_millis(100)),
        max_interval: Some(Duration::from_millis(200)),
        // duration: Some(Duration::from_secs(100)),
        discoverable: Some(true),
        // local_name: Some("SlowCon".to_string()),
        ..Default::default()
    };
    println!("{:?}", &le_advertisement);
    let handle = adapter.advertise(le_advertisement).await?;

    println!("Press enter to quit");
    let stdin = BufReader::new(tokio::io::stdin());
    let mut lines = stdin.lines();
    let _ = lines.next_line().await;

    println!("Removing advertisement");
    drop(handle);
    sleep(Duration::from_secs(1)).await;

    return Ok(());

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
                        let res = query_all_device_properties(&adapter, addr).await;
                        if let Err(err) = res {
                            println!("    Error: {}", &err);
                        }

                        let device = adapter.device(MY_ADDRESS).expect("Error getting device");
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
                break;
            }
        }
    }

    println!("Done");

    Ok(())
}
