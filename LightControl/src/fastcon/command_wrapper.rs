use std::vec;

use super::core::protocol::do_generate_command;

const BLE_CMD_RETRY_CNT: i32 = 1;
const BLE_CMD_ADVERTISE_LENGTH: i32 = 3000; // how long, in ms, to advertise for a command

fn command_with_no_delay(
    i: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_count: i32,
    send_time: i32,
    z: bool,
    use_default_adapter: bool,
    use_22_data: bool,
    i2: u8,
) -> Vec<u8> {
    command_with_delay_impl(
        i,
        data,
        key,
        retry_count,
        send_time,
        z,
        0,
        use_default_adapter,
        use_22_data,
        i2,
    )
}

fn command_with_delay(
    i: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_cnt: i32,
    send_time: i32,
    z: bool,
    delay: i32,
    use_default_adapter: bool,
    use_22_data: bool,
    i2: u8,
) -> Vec<u8> {
    command_with_delay_impl(
        i,
        data,
        key,
        retry_cnt,
        send_time,
        z,
        delay,
        use_default_adapter,
        use_22_data,
        i2,
    )
}

fn command_with_delay_impl(
    n: u8,
    data: &[u8],
    key: Option<&[u8]>,
    retry_count: i32,
    send_interval: i32,
    z: bool,
    delay: i32,
    use_default_adapter: bool,
    use_22_data: bool,
    i2: u8,
) -> Vec<u8> {
    if delay <= 0 {
        do_generate_command(
            n,
            data,
            key,
            retry_count,
            send_interval,
            z,
            use_default_adapter,
            use_22_data,
            i2,
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

pub fn single_control(addr: u32, key: Option<&[u8]>, data: Vec<u8>, delay: i32) -> Vec<u8> {
    let mut result_data = vec![0; 12];

    result_data[0] = 2 | (((0xfffffff & (data.len() + 1)) << 4) as u8);
    result_data[1] = addr as u8;
    result_data[2..data.len() + 2].copy_from_slice(&data);

    command_with_delay(
        5, // unknown value
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

pub fn single_control_nodelay(addr: u32, key: Option<&[u8]>, data: Vec<u8>) -> Vec<u8> {
    single_control(addr, key, data, 0)
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
