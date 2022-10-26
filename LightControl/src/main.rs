use btleplug::{
    api::{
        Central, CharPropFlags, Characteristic, Manager as _, Peripheral as _, ScanFilter,
        WriteType,
    },
    platform::{Adapter, Manager, Peripheral},
};
use core::panic;
use std::{error::Error, str::FromStr, time::Duration};
use tokio::time::sleep;
use uuid::Uuid;

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

async fn send_command(
    peripheral: &Peripheral,
    service_uuid: &str,
    characteristic_uuid: &str,
    data: &Vec<u8>,
) -> Result<(), btleplug::Error> {
    let service_uuid = format!("0000{}-0000-1000-8000-00805f9b34fb", service_uuid);
    let characteristic_uuid = format!("0000{}-0000-1000-8000-00805f9b34fb", characteristic_uuid);
    let cc = Characteristic {
        service_uuid: Uuid::from_str(&service_uuid).unwrap(),
        uuid: Uuid::from_str(&characteristic_uuid).unwrap(),
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

    let on_command = vec![0xcc, 0x23, 0x33];
    let off_command = vec![0xcc, 0x24, 0x33];
    send_command(&light, "ffd5", "ffd9", &on_command).await?;
    sleep(Duration::from_secs(1)).await;

    send_command(&light, "ffd5", "ffd9", &off_command).await?;
    sleep(Duration::from_secs(1)).await;

    let red_command = vec![0x56, 0xff, 0x00, 0x00, 0x00, 0xf0, 0xaa];
    let green_command = vec![0x56, 0x00, 0xff, 0x00, 0x00, 0xf0, 0xaa];
    let blue_command = vec![0x56, 0x00, 0x00, 0xff, 0x00, 0xf0, 0xaa];
    let nice_command = vec![0x56, 0x00, 0x00, 0x00, 0xff, 15, 0xaa];

    send_command(&light, "ffd5", "ffd9", &on_command).await?;
    sleep(Duration::from_secs(1)).await;

    loop {
        for command in vec![&red_command, &green_command, &blue_command, &nice_command] {
            send_command(&light, "ffd5", "ffd9", command).await?;
            sleep(Duration::from_secs(1)).await;
        }
    }
}
