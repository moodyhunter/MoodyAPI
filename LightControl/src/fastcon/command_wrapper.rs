use core::panic;
use std::{
    time::{SystemTime, UNIX_EPOCH},
    vec,
};

use super::{
    utils::{whitening_encode, whitening_init, WhiteningContext},
    DEFAULT_BLE_FASTCON_ADDRESS, DEFAULT_ENCRYPT_KEY,
};
use crate::fastcon::common::print_bytes;
use crate::fastcon::utils::{crc16, reverse_8};
use num_traits::WrappingAdd;

const BLE_CMD_RETRY_CNT: i32 = 1;
const BLE_CMD_ADVERTISE_LENGTH: i32 = 3000; // how long, in ms, to advertise for a command

#[allow(unused)]
enum LightCommand {
    RGB(
        bool, /*on*/
        u8,   /*brightness*/
        u8,
        u8,
        u8,   /*RGB*/
        bool, /*absolute (do not normalize)*/
    ),
    WarmWhite(
        bool, /* on */
        u8,   /* brightness */
        u8,   /* warm val */
        u8,   /* cold val */
    ),
    Brightness(bool /* on */, u8 /* brightness */),
    OnOff(bool /* on */, u8 /* brightness */),
    // Some5(bool /* z2 */),
    // Some6,
    // Some7,
    // Some8,
    // Some9,
}

fn single_light_command(
    // on_off: bool,
    // brightness: u8,
    // g: i32,
    // b: i32,
    // r: i32,
    // i5: i32,
    // i6: i32,
    // z2: bool,
    // i7: i32,
    command: LightCommand,
    // maybe_batch: bool,
    // relative_color: bool,
) -> Vec<u8> {
    match command {
        LightCommand::RGB(on, brightness, r, g, b, abs) => {
            let mut arr = vec![0; 6];
            let color_normalisation = if abs {
                1.0
            } else {
                255.0 / (r as u32 + g as u32 + b as u32) as f32
            };
            arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
            arr[1] = ((b as f32) * color_normalisation) as u32 as u8;
            arr[2] = ((r as f32) * color_normalisation) as u32 as u8;
            arr[3] = ((g as f32) * color_normalisation) as u32 as u8;
            arr[4] = 0;
            arr[5] = 0;
            arr
        }
        LightCommand::WarmWhite(on, brightness, i5, i6) => {
            let mut arr = vec![0; 6];
            arr[0] = ((if on { 128 } else { 0 }) + (brightness & 127)) as u8;
            arr[1] = 0;
            arr[2] = 0;
            arr[3] = 0;
            arr[4] = i5 as u8;
            arr[5] = i6 as u8;
            arr
        }
        LightCommand::Brightness(on, val) => {
            // if maybe_batch {
            //     let mut arr = vec![0; 6];
            //     arr[0] = (if on { val & 127 } else { 0 }) as u8;
            //     arr[1] = 0;
            //     arr[2] = 0;
            //     arr[3] = 0;
            //     arr[4] = 0;
            //     arr[5] = 0;
            //     return arr;
            // }
            vec![if on { val & 127 } else { 0 } as u8]
        }
        LightCommand::OnOff(on, brightness) => {
            // if maybe_batch {
            //     let mut arr = vec![0; 6];
            //     arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
            //     arr[1] = 0;
            //     arr[2] = 0;
            //     arr[3] = 0;
            //     arr[4] = 0;
            //     arr[5] = 0;
            //     return arr;
            // }

            vec![if on { 128 } else { 0 } + (brightness & 127) as u8]
        } //
          // SingleLightCommand::Some5(z) => {
          //     let mut arr = vec![0; 7];
          //     arr[0] = 0;
          //     arr[1] = 0;
          //     arr[2] = 0;
          //     arr[3] = 0;
          //     arr[4] = u8::MAX;
          //     arr[5] = u8::MAX;
          //     arr[6] = if z { 128 } else { 0 } as u8;
          //     arr
          // }
          // SingleLightCommand::Some6 => {
          //     let mut arr = vec![0; 7];
          //     arr[0] = 0;
          //     arr[1] = 0;
          //     arr[2] = 0;
          //     arr[3] = 0;
          //     arr[4] = u8::MAX;
          //     arr[5] = u8::MAX;
          //     arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
          //     arr
          // }
          // SingleLightCommand::Some7 => {
          //     let color_normalisation = if relative_color {
          //         255.0 / (r + g + b) as f32
          //     } else {
          //         1.0
          //     };
          //     let mut arr = vec![0; 7];
          //     arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
          //     arr[1] = ((r as f32) * color_normalisation) as u32 as u8;
          //     arr[2] = ((g as f32) * color_normalisation) as u32 as u8;
          //     arr[3] = ((b as f32) * color_normalisation) as u32 as u8;
          //     arr[4] = i5 as u8;
          //     arr[5] = i6 as u8;
          //     arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
          //     arr
          // }
          // SingleLightCommand::Some8 => {
          //     let mut arr = vec![0; 7];
          //     arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
          //     arr[1] = u8::MAX;
          //     arr[2] = u8::MAX;
          //     arr[3] = u8::MAX;
          //     arr[4] = u8::MAX;
          //     arr[5] = u8::MAX;
          //     arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
          //     arr
          // }
          // SingleLightCommand::Some9 => {
          //     let color_normalisation = if relative_color {
          //         255.0 / (r + g + b) as f32
          //     } else {
          //         1.0
          //     };
          //     let mut arr = vec![0; 7];
          //     arr[0] = u8::MAX;
          //     arr[1] = ((r as f32) * color_normalisation) as u32 as u8;
          //     arr[2] = ((g as f32) * color_normalisation) as u32 as u8;
          //     arr[3] = ((b as f32) * color_normalisation) as u32 as u8;
          //     arr[4] = u8::MAX;
          //     arr[5] = u8::MAX;
          //     arr[6] = u8::MIN;
          //     arr
          // }
    }
}

fn package_ble_fastcon_body(
    i: u8,
    i2: u8,
    sequence: u32,
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
    body[0] = ((i2 & 0b1111) << 0) | ((i & 0b111) << 4) | (u8::from(forward) << 7);
    body[1] = sequence as u8;
    body[2] = safe_key;
    body[3] = 0; // checksum

    body[4..].copy_from_slice(data);

    let mut checksum = 0;
    for i in 0..body.len() {
        if i == 3 {
            continue; // skip checksum itself
        }

        // allow overflow
        checksum = checksum.wrapping_add(&body[i]);
    }

    body[3] = checksum;

    let real_key = key.unwrap_or(&DEFAULT_ENCRYPT_KEY);

    for i in 0..4 {
        body[i] = super::DEFAULT_ENCRYPT_KEY[i & 3] ^ body[i];
    }

    for i in 0..data_size {
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

    println!("send seq: {}, send count: {}", some_sequence, send_cnt);

    let safe_key: u8 = if key.is_some() { key.unwrap()[3] } else { 255 };

    if use_22_data {
        panic!("22 data not implemented");
    } else {
        package_ble_fastcon_body(
            i,
            i2,
            some_sequence,
            safe_key,
            forward,
            data,
            data.len(),
            key,
        )
    }
}

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

fn do_generate_command(
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
        todo!("what is something?")
    }

    if !use_default_adapter {
        todo!("use specific adapter")
    }

    let mut payload =
        get_payload_with_inner_retry(n, data, std::cmp::max(i4, 0), key, forward, use_22_data);

    payload = get_rf_payload(&DEFAULT_BLE_FASTCON_ADDRESS, &payload);

    let mut context = WhiteningContext::new();
    whitening_init(0x25, &mut context);
    whitening_encode(&mut payload, &mut context);
    payload[0xf..].to_vec() // drop the first 0xf bytes
}

fn package_device_control(device_id: u8, src_buf: &[u8], src_len: usize, result: &mut [u8]) {
    result[0] = 2 | (((0xfffffff & (src_len + 1)) << 4) as u8);
    result[1] = device_id;
    result[2..src_len + 2].copy_from_slice(src_buf);
}

fn command_with_no_delay(
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_count: i32,
    send_time: i32,
    z: bool,
    use_default_adapter: bool,
    use_22_data: bool,
    some_i: u8,
) -> Vec<u8> {
    command_with_delay_impl(
        n,
        data,
        key,
        retry_count,
        send_time,
        z,
        0,
        use_default_adapter,
        use_22_data,
        some_i,
    )
}

fn command_with_delay(
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_cnt: i32,
    send_time: i32,
    z: bool,
    delay: i32,
    use_default_adapter: bool,
    use_22_data: bool,
    some_i: u8,
) -> Vec<u8> {
    command_with_delay_impl(
        n,
        data,
        key,
        retry_cnt,
        send_time,
        z,
        delay,
        use_default_adapter,
        use_22_data,
        some_i,
    )
}

fn command_with_delay_impl(
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    i2: i32,
    i3: i32,
    z: bool,
    delay: i32,
    use_default_adapter: bool,
    use_22_data: bool,
    i5: u8,
) -> Vec<u8> {
    if delay <= 0 {
        do_generate_command(
            n,
            data,
            key,
            i2,
            i3,
            z,
            use_default_adapter,
            use_22_data,
            i5,
        )
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

fn control_device_with_delay(addr: i32, data: Vec<u8>, key: Option<&[u8]>, delay: i32) -> Vec<u8> {
    print_bytes("Control with device Data", &data);
    let mut result_data = vec![0; 12];
    package_device_control(addr as u8, &data, data.len(), &mut result_data);
    command_with_delay(
        5,
        &result_data,
        key,
        BLE_CMD_RETRY_CNT,
        BLE_CMD_ADVERTISE_LENGTH,
        true,
        delay,
        true,
        addr > 256,
        (addr / 256).try_into().expect("addr / 256 beyond u8"),
    )
}

fn single_control(addr: i32, data: Vec<u8>, key: Option<&[u8]>) -> Vec<u8> {
    control_device_with_delay(addr, data, key, 0)
}

pub fn single_on_off_command(key: Option<&[u8]>, short_addr: i32, on: bool) -> Vec<u8> {
    println!(
        "single_on_off_command: short_addr: {:04x}, on: {}",
        short_addr, on
    );

    let command = single_light_command(LightCommand::OnOff(on, if on { 255 } else { 0 }));

    single_control(short_addr, command, key)
}

pub fn single_brightness_command(key: Option<&[u8]>, short_addr: i32, brightness: u8) -> Vec<u8> {
    println!(
        "single_brightness_command: short_addr: {:04x}, brightness: {}",
        short_addr, brightness
    );

    let command = single_light_command(LightCommand::Brightness(brightness > 0, brightness));

    single_control(short_addr, command, key)
}

pub fn single_warmwhite(key: Option<&[u8]>, short_addr: i32, on: bool, brightness: u8) -> Vec<u8> {
    let command = single_light_command(LightCommand::WarmWhite(on, brightness, 127, 127));
    single_control(short_addr, command, key)
}

pub fn single_rgb_command(
    key: Option<&[u8]>,
    short_addr: i32,
    on: bool,
    brightness: u8,
    r: u8,
    g: u8,
    b: u8,
    absoulte: bool,
) -> Vec<u8> {
    println!(
        "single_rgb_command: short_addr: {:04x}, r: {}, g: {}, b: {}",
        short_addr, r, g, b
    );

    let command = single_light_command(LightCommand::RGB(on, brightness, r, g, b, absoulte));

    single_control(short_addr, command, key)
}

pub fn command_start_scan() -> Vec<u8> {
    command_with_no_delay(
        0,
        &vec![0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        None,
        BLE_CMD_RETRY_CNT,
        -1,
        false,
        false,
        false,
        0,
    )
}
