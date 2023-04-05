// This is a reverse-engineered implementation of the BLE FastCon protocol
// used by some light bulbs.

pub const DEFAULT_PHONE_KEY: [u8; 4] = [0xA1, 0xA2, 0xA3, 0xA4];
pub const DEFAULT_ENCRYPT_KEY: [u8; 4] = [0x5e, 0x36, 0x7b, 0xc4];

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

// common types
pub mod common;

// parse the broadcast packets
pub mod broadcast_parser;

// wrapper around some of the commands
pub mod command_wrapper;

// device type enum
pub mod device_type;
