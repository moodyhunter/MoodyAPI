// This is a reverse-engineered implementation of the BLE FastCon protocol
// used by some light bulbs.

pub const DEFAULT_PHONE_KEY: [u8; 4] = [0xA1, 0xA2, 0xA3, 0xA4];
pub const DEFAULT_ENCRYPT_KEY: [u8; 4] = [0x5e, 0x36, 0x7b, 0xc4];
pub const DEFAULT_BLE_FASTCON_ADDRESS: [u8; 3] = [0xC1, 0xC2, 0xC3];
use num_derive::{FromPrimitive, ToPrimitive};

#[derive(Debug, Clone, FromPrimitive, ToPrimitive)]
#[allow(unused, non_camel_case_types)] // just lazy to rename
pub enum DeviceType {
    Unknown = 0, // unknown device type

    Curtain = 43499,
    Fan = 43531,
    Gateway = 43500,
    Gateway_AC = 43756,
    Gateway_IHG = 10058,
    Light_BURDEN_CW = 43754,
    Light_BURDEN_W = 43759,
    Light_CCT = 43051,
    Light_COMPOSE = 43709,
    Light_PWR = 43049,
    Light_RGB = 43168,
    Light_RGBCW = 43050,
    Light_RGBW = 43169,
    Light_W_CW = 43745,
    Meta_PAD = 43518,
    Meta_PAD_2 = 43974,
    Panel_3 = 43463,
    Panel_3_Wireless = 43462,
    Panel_4 = 43473,
    Panel_4_Wireless = 43472,
    Panel_6 = 43461,
    Panel_6_Wireless = 43459,
    Panel_8 = 43733,
    Panel_8_2 = 43734,
    Relay_1 = 43525,
    Relay_2 = 43474,
    Relay_4 = 43680,
    Sensor_Door = 43505,
    Sensor_IR = 43516,
    Sensor_Radar = 43808,
    Sensor_Water = 43791,
    ThermoStat = 43919,
}

#[macro_export]
macro_rules! fastcon_ble_header_encrypt {
    ($src:expr, $dst:expr, $size:expr) => {
        for i in 0..$size {
            $dst[i] = super::DEFAULT_ENCRYPT_KEY[i & 3] ^ $src[i];
        }
    };
}

#[macro_export]
macro_rules! fastcon_ble_encrypt {
    ($src:expr, $dst:expr, $size:expr, $key:expr) => {
        for i in 0..$size {
            $dst[i] = $key[i & 3] ^ $src[i];
        }
    };
}

mod core; // core protocol

// parse the broadcast packets
pub mod broadcast_parser;

// wrapper around some of the commands
pub mod command_wrapper;

// BLE Devices
pub mod light;
