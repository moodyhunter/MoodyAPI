use num_traits::FromPrimitive;

use crate::{fastcon_ble_encrypt, fastcon_ble_header_encrypt};

use super::DeviceType;

#[derive(Debug, Clone)]
pub struct HeartBeat {
    pub version: String,
    pub short_addr: i32,
    pub group_addr: u8,
}

#[derive(Debug, Clone)]
pub struct DeviceInfo {
    pub did: Vec<u8>,
    pub key: Vec<u8>,
    pub device_type: DeviceType,
    pub high: u8,
    pub cnt: u8,
}

pub enum BroadcastType {
    HeartBeat(HeartBeat),
    TimerUploadResponse,
    DeviceAnnouncement(DeviceInfo),
}

pub fn parse_ble_broadcast(source: &[u8], phone_key: &[u8; 4]) -> Option<BroadcastType> {
    let mut header = source[0..4].to_vec();
    fastcon_ble_header_encrypt!(source, header, 4);

    let high = header[0] & 0xf; // some strange high bits

    // high 3 bits
    match (header[0] >> 4) & 7 {
        3 => {
            let mut content = source[4..].to_vec(); // skip 4 bytes of header
            fastcon_ble_encrypt!(&source[4..], content, source.len() - 4, phone_key);

            match content[0] & 0xf {
                0xb => {
                    println!("todo: timer upload response");
                    Some(BroadcastType::TimerUploadResponse)
                }
                0x4 => {
                    let version = format!(
                        "{}.{}.{}.{}.{}",
                        // note: our 'content' array does not include the first 4-byte header,
                        // so subtract 4 from the offsets
                        // From the disassembly:
                        // ldrb    w2, [x22, #0x7]
                        content[3] as u32,
                        // ldrb    w3, [x22, #0x8]
                        content[4] as u32,
                        // ldrh    w4, [x22, #0xe]
                        content[10] as u32 | (content[11] as u32) << 8,
                        // ldur    w5, [x22, #0xa]
                        content[6] as u32 | (content[7] as u32) << 8,
                        // ldrb    w6, [x22, #0x9]
                        content[5] as u32
                    );

                    // const int addr = (uint32_t) data_buf[5] | (*data_buf & 0xf) << 8;
                    // const int group_addr = data_buf[6];

                    let short_addr = (content[1] as i32) | (high as i32) << 8;
                    let group_addr = content[2];

                    Some(BroadcastType::HeartBeat(HeartBeat {
                        version,
                        short_addr,
                        group_addr,
                    }))
                }
                unknown => {
                    println!("Unknown content type: {}", unknown);
                    None
                }
            }
        }
        1 => {
            // 4eb17a50ec0bf10f00e9a1a85e367bc4
            // did: ec0bf10f00e9
            // key: 5e367bc4
            // type: 43169 (0xa8a1)

            let did = source[4..10].to_vec(); // 6 bytes
            let type_buffer = source[10..12].to_vec(); // 2 bytes
            let key = source[12..16].to_vec(); // 4 bytes
            let dev_type = type_buffer[0] as u16 | (type_buffer[1] as u16) << 8;

            Some(BroadcastType::DeviceAnnouncement(DeviceInfo {
                cnt: 1, // seems to be hardcoded to 1
                key,
                did,
                device_type: (FromPrimitive::from_u16(dev_type) as Option<DeviceType>)
                    .unwrap_or(DeviceType::Unknown),
                high,
            }))
        }
        other => {
            println!("Unknown header type: {}", other);
            None
        }
    }
}
