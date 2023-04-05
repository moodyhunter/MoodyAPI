use core::panic;

use super::{DEFAULT_ENCRYPT_KEY, DEFAULT_PHONE_KEY};
use crate::{fastcon_ble_encrypt, fastcon_ble_header_encrypt};

const BLE_CMD_RETRY_CNT: i32 = 1;
const BLE_CMD_ADVERTISE_LENGTH: i32 = 3000; // how long, in ms, to advertise for a command

#[allow(unused)]
enum SingleLightCommand {
    MaybeRGB,
    MaybeSetWhite,
    Brightness,
    OnOff,
    Some5,
    Some6,
    Some7,
    Some8,
    Some9,
}

fn generate_single_light_command(
    on_off: bool,
    brightness: i32,
    g: i32,
    b: i32,
    r: i32,
    i5: i32,
    i6: i32,
    z2: bool,
    i7: i32,
    command: SingleLightCommand,
    z3: bool,
    z4: bool,
) -> Vec<u8> {
    let i9 = i5;
    let i10 = i6;
    let maybe_r = r & 255;
    let maybe_g = g & 255;
    let maybe_b = b & 255;
    let color_normalisation = if z4 {
        255.0 / ((maybe_r + maybe_g) + maybe_b) as f32
    } else {
        1.0
    };

    match command {
        SingleLightCommand::MaybeRGB => {
            let mut arr = vec![0; 6];
            arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
            arr[1] = ((maybe_r as f32) * color_normalisation) as u32 as u8;
            arr[2] = ((maybe_g as f32) * color_normalisation) as u32 as u8;
            arr[3] = ((maybe_b as f32) * color_normalisation) as u32 as u8;
            arr[4] = 0;
            arr[5] = 0;
            arr
        }
        SingleLightCommand::MaybeSetWhite => {
            let mut arr = vec![0; 6];
            arr[0] = ((if on_off { 128 } else { 0 }) + (brightness & 127)) as u8;
            arr[1] = 0;
            arr[2] = 0;
            arr[3] = 0;
            arr[4] = i9 as u8;
            arr[5] = i10 as u8;
            arr
        }
        SingleLightCommand::Brightness => {
            let mut arr = vec![0; if z3 { 6 } else { 1 }];
            arr[0] = (if on_off { brightness & 127 } else { 0 }) as u8;
            arr
        }
        SingleLightCommand::OnOff => {
            let mut arr = vec![0; 6];
            arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
            arr[1] = 0;
            arr[2] = 0;
            arr[3] = 0;
            arr[4] = 0;
            arr[5] = 0;
            arr
        }
        SingleLightCommand::Some5 => {
            let mut arr = vec![0; 7];
            arr[0] = 0;
            arr[1] = 0;
            arr[2] = 0;
            arr[3] = 0;
            arr[4] = u8::MAX;
            arr[5] = u8::MAX;
            arr[6] = if z2 { 128 } else { 0 } as u8;
            arr
        }
        SingleLightCommand::Some6 => {
            let mut arr = vec![0; 7];
            arr[0] = 0;
            arr[1] = 0;
            arr[2] = 0;
            arr[3] = 0;
            arr[4] = u8::MAX;
            arr[5] = u8::MAX;
            arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
            arr
        }
        SingleLightCommand::Some7 => {
            let mut arr = vec![0; 7];
            arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
            arr[1] = ((maybe_r as f32) * color_normalisation) as u32 as u8;
            arr[2] = ((maybe_g as f32) * color_normalisation) as u32 as u8;
            arr[3] = ((maybe_b as f32) * color_normalisation) as u32 as u8;
            arr[4] = i9 as u8;
            arr[5] = i10 as u8;
            arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
            arr
        }
        SingleLightCommand::Some8 => {
            let mut arr = vec![0; 7];
            arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
            arr[1] = u8::MAX;
            arr[2] = u8::MAX;
            arr[3] = u8::MAX;
            arr[4] = u8::MAX;
            arr[5] = u8::MAX;
            arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
            arr
        }
        SingleLightCommand::Some9 => {
            let mut arr = vec![0; 7];
            arr[0] = u8::MAX;
            arr[1] = ((maybe_r as f32) * color_normalisation) as u32 as u8;
            arr[2] = ((maybe_g as f32) * color_normalisation) as u32 as u8;
            arr[3] = ((maybe_b as f32) * color_normalisation) as u32 as u8;
            arr[4] = u8::MAX;
            arr[5] = u8::MAX;
            arr[6] = u8::MIN;
            arr
        }
    }
}

fn generate_on_off_command(on: bool) -> Vec<u8> {
    generate_single_light_command(
        on,
        0,
        0,
        0,
        0,
        0,
        0,
        false,
        0,
        SingleLightCommand::OnOff,
        false,
        false,
    )
}
// int i, int i2, int i3, int i4, int i5, byte[] bArr, int i6, byte[] bArr2, byte[] bArr3
fn package_ble_fastcon_body(
    i: i32,
    i2: u8,
    sequence: i32,
    safe_key: u8,
    forward: bool,
    data: &[u8],
    data_size: usize,
    key: Option<&[u8]>,
) -> Vec<u8> {
    // bit 7 is forward
    // bit 6-4 is i2
    // bit 3-0 is i
    let mut body = vec![0; data.len() + 4];
    body[0] = (i & 15) as u8 + ((i2 & 7) << 4) + (u8::from(forward) << 7);
    body[1] = sequence as u8;
    body[2] = safe_key;
    body[3] = 0; // checksum
    body[4..].copy_from_slice(data);

    let mut checksum = 0;
    for i in 0..body.len() {
        if i == 3 {
            continue; // skip checksum itself
        }
        checksum += body[i];
    }

    body[3] = checksum;

    let real_key = key.unwrap_or(&DEFAULT_ENCRYPT_KEY);
    fastcon_ble_header_encrypt!(body, body, 4);
    fastcon_ble_encrypt!(body[4..], body[4..], data_size, real_key);
    body
}

fn get_payload_with_inner_retry(
    i: i32,
    data: &[u8],
    i2: u8,
    key: Option<&[u8]>,
    maybe_forward: bool,
    use_22_data: bool,
) -> Vec<u8> {
    println!(
        "get_payload_with_inner_retry: payload={:?}, key={:?}",
        data, key
    );

    static mut SEND_COUNT: i32 = 0;
    static mut SEND_SEQ: i32 = 0; // static variables so they persist between calls

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

    //
    println!(
        "send_cnt={}, send_seq={}, seq={}",
        send_cnt,
        unsafe { SEND_SEQ },
        some_sequence
    );

    // let result_len = if use_22_data { 22 } else { 16 };
    let safe_key: u8 = if key.is_some() { key.unwrap()[3] } else { 255 };

    if use_22_data {
        panic!("22 data not implemented");
    } else {
        package_ble_fastcon_body(
            i,
            i2,
            some_sequence,
            safe_key,
            maybe_forward,
            data,
            data.len(),
            key,
        )
    }
}

fn do_send_command(
    i: i32,
    data: &[u8],
    key: Option<&[u8]>,
    retry_count: i32,
    send_interval: i32,
    maybe_forward: bool,
    use_specific_adapter: bool,
    use_22_data: bool,
    i4: u8,
) -> Vec<u8> {
    let payload = get_payload_with_inner_retry(
        i,
        data,
        std::cmp::max(i4, 0),
        key,
        maybe_forward,
        use_22_data,
    );

    payload
}

fn package_device_control(device_id: u8, src_buf: &[u8], src_len: usize, result: &mut [u8]) {
    result[0] = (2 | ((0xfffffff & (src_len + 1)) << 4)) as u8;
    result[1] = device_id;
    result[2..src_len + 2].copy_from_slice(src_buf);
}

fn command_with_no_delay(
    n: i32,
    data: &[u8],
    key: Option<&[u8]>,
    i2: i32,
    i3: i32,
    z: bool,
    z2: bool,
    use_22_data: bool,
    i4: u8,
) -> Vec<u8> {
    command_with_delay_impl(n, data, key, i2, i3, z, 0, z2, use_22_data, i4)
}

fn command_with_delay(
    n: i32,
    data: &[u8],
    key: Option<&[u8]>,
    retry_cnt: i32,
    send_time: i32,
    z: bool,
    i4: i32,
    z2: bool,
    use_22_data: bool,
    i5: u8,
) -> Vec<u8> {
    command_with_delay_impl(
        n,
        data,
        key,
        retry_cnt,
        send_time,
        z,
        i4,
        z2,
        use_22_data,
        i5,
    )
}

// public boolean sendCommand(int i, byte[] bArr, byte[] bArr2, int i2, int i3, boolean z, boolean z2, int i4, boolean z3, boolean z4, int i5)
fn command_with_delay_impl(
    n: i32,
    data: &[u8],
    key: Option<&[u8]>,
    i2: i32,
    i3: i32,
    z: bool,
    delay: i32,
    z3: bool,
    use_22_data: bool,
    i5: u8,
) -> Vec<u8> {
    if delay <= 0 {
        do_send_command(n, data, key, i2, i3, z, z3, use_22_data, i5)
    } else {
        panic!("delay not implemented");
    }
    /*
    if (this.mEnable || System.currentTimeMillis() - this.mEnableWatchDog > 2000)
    {
        this.mEnable = false;
        if (!doSendCommand(n, data, key, i2, i3, z, has_key, z3, z4, i5))
        {
            return false;
        }
        new Timer().schedule(new TimerTask() {
            public void run()
            {
                boolean unused = BLEFastconHelper.this.mEnable = true;
                long unused2 = BLEFastconHelper.this.mEnableWatchDog = System.currentTimeMillis();
            }
        }, (long) i4);
        return true;
    }
    ELogUtils.m22w("jyq_music", "sendCommand throttle throw out.");
    return false;
    */
}

fn control_with_device(addr: i32, data: Vec<u8>, i2: i32) -> Vec<u8> {
    /*
       byte[] parseStringToByte = data;
       byte[] bArr = new byte[12];
       BLEUtil.package_device_control(i, parseStringToByte, parseStringToByte.length, bArr);
       return sendCommand(5, bArr, this.mPhoneKey, BLE_CMD_RETRY_CNT, BLE_CMD_SEND_TIME, true, i2, true, i > 256, i / 256);
    */

    let mut result_data = vec![0; 12];
    package_device_control(addr as u8, &data, data.len(), &mut result_data);
    command_with_delay(
        5,
        &result_data,
        Some(&DEFAULT_PHONE_KEY),
        BLE_CMD_RETRY_CNT,
        BLE_CMD_ADVERTISE_LENGTH,
        true,
        i2,
        true,
        addr > 256,
        (addr / 256).try_into().expect("addr / 256 beyond u8"),
    )
}

fn send_single_control(addr: i32, data: Vec<u8>) {
    control_with_device(addr, data, 0);
}

pub fn single_on_off_command(short_addr: i32, on: bool) {
    println!(
        "single_on_off_command: short_addr: {:04x}, on: {}",
        short_addr, on
    );

    send_single_control(short_addr, generate_on_off_command(on));
}

pub fn command_start_scan() -> Vec<u8> {
    let packet = vec![0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
    command_with_no_delay(
        0,
        &packet,
        None,
        BLE_CMD_RETRY_CNT,
        -1,
        false,
        false,
        false,
        0,
    )
}
