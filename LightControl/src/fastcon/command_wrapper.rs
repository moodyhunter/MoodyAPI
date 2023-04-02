pub const DEFAULT_PHONE_KEY: [u8; 4] = [0xA1, 0xA2, 0xA3, 0xA4];

const BLE_CMD_RETRY_CNT: i32 = 1;
const BLE_CMD_ADVERTISE_LENGTH: i32 = 3000; // how long, in ms, to advertise for a command

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
    mode: SingleLightCommand,
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

    match mode {
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

/*
void package_device_control(char device_id, char const *src_buf, uint64_t src_len, char *result)
{
    result[0] = (2 | ((int8_t) ((0xfffffff & (src_len + 1)) << 4)));
    result[1] = device_id;
    memcpy(&result[2], src_buf, src_len);
}
*/

fn package_device_control(device_id: u8, src_buf: &[u8], src_len: usize, result: &mut [u8]) {
    result[0] = (2 | ((0xfffffff & (src_len + 1)) << 4)) as u8;
    result[1] = device_id;
    result[2..src_len + 2].copy_from_slice(src_buf);
}

/*

   public boolean sendCommand(int i, byte[] bArr, byte[] bArr2, int i2, int i3, boolean z, boolean z2, boolean z3, int i4)
   {
       return sendCommand(i, bArr, bArr2, i2, i3, z, bArr2 != null, 0, z2, z3, i4);
   }

   public boolean sendCommand(int i, byte[] bArr, byte[] bArr2, int i2, int i3, boolean z, int i4, boolean z2, boolean z3, int i5)
   {
       return sendCommand(i, bArr, bArr2, i2, i3, z, bArr2 != null, i4, z2, z3, i5);
   }
*/
fn send_command_with_no_delay(
    n: i32,
    data: &[u8],
    key: &[u8],
    i2: i32,
    i3: i32,
    z: bool,
    z2: bool,
    z3: bool,
    i4: i32,
) -> bool {
    send_command_with_delay_impl(n, data, key, i2, i3, z, key != &[], 0, z2, z3, i4)
}

fn send_command_with_delay(
    n: i32,
    data: &[u8],
    key: &[u8],
    retry_cnt: i32,
    send_time: i32,
    z: bool,
    i4: i32,
    z2: bool,
    z3: bool,
    i5: i32,
) -> bool {
    send_command_with_delay_impl(
        n,
        data,
        key,
        retry_cnt,
        send_time,
        z,
        key != &[],
        i4,
        z2,
        z3,
        i5,
    )
}

// public boolean sendCommand(int i, byte[] bArr, byte[] bArr2, int i2, int i3, boolean z, boolean z2, int i4, boolean z3, boolean z4, int i5)
fn send_command_with_delay_impl(
    n: i32,
    data: &[u8],
    key: &[u8],
    i2: i32,
    i3: i32,
    z: bool,
    has_key: bool,
    delay: i32,
    z3: bool,
    z4: bool,
    i5: i32,
) -> bool {
    if delay <= 0 {
        // return do_send_command(n, data, key, i2, i3, z, has_key, z3, z4, i5);
    }
    false
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

fn control_with_device(addr: i32, data: Vec<u8>, i2: i32) {
    /*
       byte[] parseStringToByte = data;
       byte[] bArr = new byte[12];
       BLEUtil.package_device_control(i, parseStringToByte, parseStringToByte.length, bArr);
       return sendCommand(5, bArr, this.mPhoneKey, BLE_CMD_RETRY_CNT, BLE_CMD_SEND_TIME, true, i2, true, i > 256, i / 256);
    */

    let mut result_data = vec![0; 12];
    package_device_control(addr as u8, &data, data.len(), &mut result_data);
    send_command_with_delay(
        5,
        &result_data,
        &DEFAULT_PHONE_KEY,
        BLE_CMD_RETRY_CNT,
        BLE_CMD_ADVERTISE_LENGTH,
        true,
        i2,
        true,
        addr > 256,
        addr / 256,
    );
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
