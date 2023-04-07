use std::time::{SystemTime, UNIX_EPOCH};

use crate::fastcon::{DEFAULT_BLE_FASTCON_ADDRESS, DEFAULT_ENCRYPT_KEY};

use super::utils::{crc16, reverse_8, whitening_encode, whitening_init, WhiteningContext};

fn get_rf_payload(addr: &[u8], data: &[u8]) -> Vec<u8> {
    let data_offset = 0x12;
    let inverse_offset = 0x0f;
    let result_data_size = data_offset + addr.len() + data.len(); // yes, the first 0x12 bytes are garbage
    let mut resultbuf = vec![0; result_data_size + 2]; // with 2 byte checksum

    // some hardcoded values
    resultbuf[0x0f] = 0x71;
    resultbuf[0x10] = 0x0f;
    resultbuf[0x11] = 0x55;

    // reverse copy the address
    for i in 0..addr.len() {
        resultbuf[data_offset + addr.len() - i - 1] = addr[i];
    }

    resultbuf[data_offset + addr.len()..data_offset + addr.len() + data.len()]
        .copy_from_slice(data);

    for i in inverse_offset..inverse_offset + addr.len() + 3 {
        resultbuf[i] = reverse_8(resultbuf[i]);
    }

    let crc = crc16(addr, data);
    resultbuf[result_data_size] = crc as u8;
    resultbuf[result_data_size + 1] = (crc >> 8) as u8;
    resultbuf
}

fn package_ble_fastcon_body(
    i: u8,
    i2: u8,
    sequence: u32,
    safe_key: u8,
    forward: bool,
    data: &[u8],
    key: Option<&[u8]>,
) -> Vec<u8> {
    // bit 7 is forward
    // bit 6-4 is i2
    // bit 3-0 is i
    let mut body = vec![0; data.len() + 4];
    body[0] = (i2 & 0b1111) << 0 | (i & 0b111) << 4 | u8::from(forward) << 7;
    body[1] = sequence as u8;
    body[2] = safe_key;
    body[3] = 0; // checksum

    body[4..].copy_from_slice(data);

    let mut checksum: u8 = 0;
    for i in 0..body.len() {
        if i == 3 {
            continue; // skip checksum itself
        }

        // allow overflow
        checksum = checksum.wrapping_add(body[i]);
    }

    body[3] = checksum;

    for i in 0..4 {
        body[i] = DEFAULT_ENCRYPT_KEY[i & 3] ^ body[i];
    }

    let real_key = key.unwrap_or(&DEFAULT_ENCRYPT_KEY);
    for i in 0..data.len() {
        body[4 + i] = real_key[i & 3] ^ body[4 + i];
    }

    body
}

fn get_payload_with_inner_retry(
    i: u8,
    data: &[u8],
    i2: u8,
    key: Option<&[u8]>,
    forward: bool,
    use_22_data: bool,
) -> Vec<u8> {
    static mut SEND_SEQ: u32 = 0;
    static mut SEND_COUNT: i32 = 0;

    unsafe {
        if SEND_SEQ == 0 {
            SEND_SEQ = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_millis() as u32
                % 256;
        }
    }

    let send_cnt = unsafe { SEND_COUNT };
    let some_sequence;

    unsafe {
        if send_cnt >= 5 || i2 <= 1 {
            let next = SEND_SEQ + 1;
            SEND_SEQ = next;
            if next == 0 || next == 256 {
                SEND_SEQ = 1; // reset sequence
            }
            some_sequence = SEND_SEQ;
        } else {
            some_sequence = (SEND_SEQ + 10) % 255;
        }
    }

    let safe_key: u8 = if key.is_some() { key.unwrap()[3] } else { 255 };

    if use_22_data {
        panic!("22 data not implemented");
    } else {
        package_ble_fastcon_body(i, i2, some_sequence, safe_key, forward, data, key)
    }
}

pub(crate) fn do_generate_command(
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    _retry_count: i32,
    _send_interval: i32,
    forward: bool,
    use_default_adapter: bool,
    use_22_data: bool,
    i4: u8,
) -> Vec<u8> {
    // TODO: handle retry_count and send_interval
    if use_22_data {
        todo!("what is use_22_data?")
    }

    if !use_default_adapter {
        todo!("use specific adapter")
    }

    let i4 = std::cmp::max(i4, 0);
    let mut payload = get_payload_with_inner_retry(n, data, i4, key, forward, use_22_data);

    payload = get_rf_payload(&DEFAULT_BLE_FASTCON_ADDRESS, &payload);

    let mut context = WhiteningContext::new();
    whitening_init(0x25, &mut context);
    whitening_encode(&mut payload, &mut context);
    payload[0xf..].to_vec() // drop the first 0xf bytes
}
