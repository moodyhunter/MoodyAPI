// This is a reverse-engineered implementation of the BLE FastCon protocol
// used by some light bulbs.

#[derive(Debug, Clone)]
pub struct HeartBeat {
    pub version: String,
    pub addr: u32,
    pub group_addr: u8,
}

pub enum BroadcastType {
    HeartBeat(HeartBeat),
    TimerUploadResponse,
    DeviceAnnouncement,
}

pub fn bl_ble_fastcon_header_encrty(src: &[u8], dst: &mut [u8], arg3: usize) {
    if arg3 == 0 {
        return;
    }

    let k = b"^6{\0"; // mysterious key reversed from the binary

    for i in 0..arg3 {
        dst[i] = k[i & 3] ^ src[i];
    }
}

pub fn bl_ble_fastcon_encrty(src: &[u8], dst: &mut [u8], arg3: usize, key: &[u8]) {
    if arg3 == 0 {
        return;
    }

    for i in 0..arg3 {
        dst[i] = key[i & 3] ^ src[i];
    }
}

pub fn parse_ble_broadcast(source: &[u8], phone_key: &[u8; 4]) -> Option<BroadcastType> {
    let mut header = source[0..4].to_vec();
    bl_ble_fastcon_header_encrty(&source, &mut header, 4);

    // high 3 bits
    match (header[0] >> 4) & 7 {
        3 => {
            let mut content = source[4..].to_vec(); // skip 4 bytes of header
            bl_ble_fastcon_encrty(&source[4..], &mut content, source.len() - 4, phone_key);

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

                        // From the binary:
                        // ldrb    w2, [x22, #0x7]
                        content[3] as u32,
                        // ldrb    w3, [x22, #0x8]
                        content[4] as u32,
                        // ldrh    w4, [x22, #0xe]
                        (content[10] as u32) | (content[11] as u32) << 8,
                        // ldur    w5, [x22, #0xa]
                        (content[6] as u32) | (content[7] as u32) << 8,
                        // ldrb    w6, [x22, #0x9]
                        content[5] as u32
                    );

                    // const int addr = (uint32_t) data_buf[5] | (*data_buf & 0xf) << 8;
                    // const int group_addr = data_buf[6];

                    let addr = (content[1] as u32) | ((header[0] & 0xf) as u32) << 8;
                    let group_addr = content[2];

                    Some(BroadcastType::HeartBeat(HeartBeat {
                        version,
                        addr,
                        group_addr,
                    }))
                }
                content_type => {
                    println!("Unknown content type: {}", content_type);
                    None
                }
            }
        }
        1 => {
            // Device Announcement
            print!("todo: Device Announcement");
            Some(BroadcastType::DeviceAnnouncement)
        }
        other => {
            println!("Unknown header type: {}", other);
            None
        }
    }
}
