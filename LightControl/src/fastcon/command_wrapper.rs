use core::panic;
use std::vec;

use num_traits::WrappingAdd;

use super::{DEFAULT_BLE_FASTCON_ADDRESS, DEFAULT_ENCRYPT_KEY, DEFAULT_PHONE_KEY};
use crate::{fastcon::common::print_bytes, fastcon_ble_encrypt, fastcon_ble_header_encrypt};

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
            if !z3 {
                vec![if on_off { 128 } else { 0 } as u8]
            } else {
                let mut arr = vec![0; 6];
                arr[0] = (if on_off { 128 } else { 0 } + (brightness & 127)) as u8;
                arr[1] = 0;
                arr[2] = 0;
                arr[3] = 0;
                arr[4] = 0;
                arr[5] = 0;
                arr
            }
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

fn package_ble_fastcon_body(
    i: u8,
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
    maybe_forward: bool,
    use_22_data: bool,
) -> Vec<u8> {
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

fn reverse_8(d: u8) -> u8 {
    let mut result = 0;
    for i in 0..8 {
        result |= ((d >> i) & 1) << (7 - i);
    }
    result
}

fn reverse_16(d: u16) -> u16 {
    let mut result = 0;
    for i in 0..16 {
        result |= ((d >> i) & 1) << (15 - i);
    }
    result
}

fn crc16(addr: &[u8], data: &[u8]) -> u16 {
    let mut crc = 0xffff;

    // iterate over address in reverse
    for i in addr.iter().rev() {
        crc ^= (*i as u16) << 8;
        for _ in 0..4 {
            let mut tmp = crc << 1;

            if crc & 0x8000 != 0 {
                tmp ^= 0x1021;
            }

            crc = tmp << 1;
            if tmp & 0x8000 != 0 {
                crc ^= 0x1021;
            }
        }
    }

    for i in 0..data.len() {
        crc ^= (reverse_8(data[i]) as u16) << 8;
        for _ in 0..4 {
            let mut tmp = crc << 1;

            if crc & 0x8000 != 0 {
                tmp ^= 0x1021;
            }

            crc = tmp << 1;
            if tmp & 0x8000 != 0 {
                crc ^= 0x1021;
            }
        }
    }

    crc = !reverse_16(crc);
    crc
}

fn get_rf_payload(addr: &[u8], data: &[u8]) -> Vec<u8> {
    print_bytes("GetRFPayload Address", addr);
    print_bytes("GetRFPayload Data", data);

    let data_offset = 0x12;
    let inverse_offset = 0x0f;
    let result_data_size = data_offset + addr.len() + data.len(); // yes, the first 0x12 bytes are garbage
    let mut resultbuf = vec![0; result_data_size + 2]; // with 2 byte checksum

    // stub data begin
    resultbuf[0x0] = 0x4c;
    resultbuf[0x1] = 0x63;
    resultbuf[0x2] = 0x6e;
    resultbuf[0x3] = 0x2f;
    resultbuf[0x4] = 0x63;
    resultbuf[0x5] = 0x6f;
    resultbuf[0x6] = 0x6d;
    resultbuf[0x7] = 0x2f;
    resultbuf[0x8] = 0x62;
    resultbuf[0x9] = 0x72;
    resultbuf[0xa] = 0x6f;
    resultbuf[0xb] = 0x61;
    resultbuf[0xc] = 0x64;
    resultbuf[0xd] = 0x6c;
    resultbuf[0xe] = 0x69;
    // stub data end

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

    print_bytes("Result", &resultbuf);

    resultbuf
}

#[derive(Debug, Clone, Copy, Default)]
struct WhiteningContext {
    f_0x0: u32,
    f_0x4: u32,
    f_0x8: u32,
    f_0xc: u32,
    f_0x10: u32,
    f_0x14: u32,
    f_0x18: u32,
}

fn whitening_init(val: u32, ctx: &mut WhiteningContext) {
    let v0 = [(val >> 5), (val >> 4), (val >> 3), (val >> 2)];

    ctx.f_0x0 = 1;
    ctx.f_0x4 = v0[0] & 1;
    ctx.f_0x8 = v0[1] & 1;
    ctx.f_0xc = v0[2] & 1;
    ctx.f_0x10 = v0[3] & 1;
    ctx.f_0x14 = (val >> 1) & 1;
    ctx.f_0x18 = val & 1;
}

fn whitening_encode(data: &mut Vec<u8>, ctx: &mut WhiteningContext) {
    print_bytes("whiten-in:", &data);
    for i in 0..data.len() {
        /*  const uint32_t varC = ctx->unkC;
        const uint32_t var14 = ctx->unk14;
        const uint32_t var18 = ctx->unk18;
        const uint32_t var10 = ctx->unk10;
        const uint32_t var8 = var14 ^ ctx->unk8;
        const uint32_t var4 = var10 ^ ctx->unk4;
        const uint32_t _var = var18 ^ varC;
        const uint32_t var0 = _var ^ ctx->unk0;

        const uint8_t c = data[i];
        data[i] = ((c & 0x80) ^ (uint8_t) ((var8 ^ var18) << 7)) + //
                  ((c & 0x40) ^ (uint8_t) (var0 << 6)) +           //
                  ((c & 0x20) ^ (uint8_t) (var4 << 5)) +           //
                  ((c & 0x10) ^ (uint8_t) (var8 << 4)) +           //
                  ((c & 0x08) ^ (uint8_t) (_var << 3)) +           //
                  ((c & 0x04) ^ (uint8_t) (var10 << 2)) +          //
                  ((c & 0x02) ^ (uint8_t) (var14 << 1)) +          //
                  ((c & 0x01) ^ (uint8_t) (var18 << 0));

        ctx->unk8 = var4;
        ctx->unkC = var8;
        ctx->unk10 = var8 ^ varC;
        ctx->unk14 = var0 ^ var10;
        ctx->unk18 = var4 ^ var14;
        ctx->unk0 = var8 ^ var18;
        ctx->unk4 = var0; */

        let varC = ctx.f_0xc;
        let var14 = ctx.f_0x14;
        let var18 = ctx.f_0x18;
        let var10 = ctx.f_0x10;
        let var8 = var14 ^ ctx.f_0x8;
        let var4 = var10 ^ ctx.f_0x4;
        let _var = var18 ^ varC;
        let var0 = _var ^ ctx.f_0x0;

        let c = data[i];
        data[i] = ((c & 0x80) ^ ((var8 ^ var18) << 7) as u8)
            + ((c & 0x40) ^ (var0 << 6) as u8)
            + ((c & 0x20) ^ (var4 << 5) as u8)
            + ((c & 0x10) ^ (var8 << 4) as u8)
            + ((c & 0x08) ^ (_var << 3) as u8)
            + ((c & 0x04) ^ (var10 << 2) as u8)
            + ((c & 0x02) ^ (var14 << 1) as u8)
            + ((c & 0x01) ^ (var18 << 0) as u8);

        ctx.f_0x8 = var4;
        ctx.f_0xc = var8;
        ctx.f_0x10 = var8 ^ varC;
        ctx.f_0x14 = var0 ^ var10;
        ctx.f_0x18 = var4 ^ var14;
        ctx.f_0x0 = var8 ^ var18;
        ctx.f_0x4 = var0;
    }

    print_bytes("whiten-out:", &data);
}

fn x_get_real_payload(
    data: Vec<u8>,
    _send_time: u32, // how long does the advertising packet last?
    use_default_adapter: bool,
    something: bool,
    _originally_i4: u8,
) -> Vec<u8> {
    if something {
        todo!("what is something?")
    } else {
        let mut rf_payload = get_rf_payload(&DEFAULT_BLE_FASTCON_ADDRESS, &data);
        print_bytes("rf payload", &rf_payload);

        if !use_default_adapter {
            todo!("use specific adapter")
        }

        let mut ctx = WhiteningContext {
            ..Default::default()
        };

        whitening_init(0x25, &mut ctx);

        whitening_encode(&mut rf_payload, &mut ctx);

        rf_payload[0xf..].to_vec()
    }
}

fn do_generate_command(
    i: u8,
    data: &[u8],
    key: Option<&[u8]>,
    _retry_count: i32,
    _send_interval: i32,
    maybe_forward: bool,
    use_default_adapter: bool,
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

    // TODO: handle retry_count and send_interval
    x_get_real_payload(payload, 0, use_default_adapter, use_22_data, i4)
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
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_cnt: i32,
    send_time: i32,
    z: bool,
    i4: i32,
    use_default_adapter: bool,
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
        use_default_adapter,
        use_22_data,
        i5,
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

fn control_with_device(addr: i32, data: Vec<u8>, key: Option<&[u8]>, i2: i32) -> Vec<u8> {
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
        i2,
        true,
        addr > 256,
        (addr / 256).try_into().expect("addr / 256 beyond u8"),
    )
}

fn send_single_control(addr: i32, data: Vec<u8>, key: Option<&[u8]>) -> Vec<u8> {
    control_with_device(addr, data, key, 0)
}

pub fn single_on_off_command(key: Option<&[u8]>, short_addr: i32, on: bool) -> Vec<u8> {
    println!(
        "single_on_off_command: short_addr: {:04x}, on: {}",
        short_addr, on
    );

    send_single_control(short_addr, generate_on_off_command(on), key)
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
