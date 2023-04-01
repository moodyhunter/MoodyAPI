mod fastcon;

use bluer::{Adapter, AdapterEvent, Address, DeviceEvent, DeviceProperty};
use futures::{stream::SelectAll, StreamExt};

const MY_ADDRESS: Address = Address::new([0x11, 0x22, 0x33, 0x44, 0x55, 0x66]);
const MY_MANUFACTURER_DATA_KEY: u16 = 0xfff0;
const DEFAULT_PHONE_KEY: [u8; 4] = [0xA1, 0xA2, 0xA3, 0xA4];

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
            let Some(x) = fastcon::parse_ble_broadcast(data, &DEFAULT_PHONE_KEY) else {
                println!("Invalid light bulb manufacturer data");
                return Ok(());
            };

            match x {
                fastcon::BroadcastType::HeartBeat(heartbeat) => {
                    println!("Heartbeat: {:?}", heartbeat);
                }
                fastcon::BroadcastType::TimerUploadResponse => {
                    println!("Timer upload response");
                }
                fastcon::BroadcastType::DeviceAnnouncement(device) => {
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
