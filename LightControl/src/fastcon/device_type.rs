use num_derive::{FromPrimitive, ToPrimitive};

#[derive(Debug, Clone, FromPrimitive, ToPrimitive)]
#[allow(unused, non_camel_case_types)] // just lazy to rename
pub enum DeviceType {
    DeviceType_Unknown = 0, // unknown device type

    DeviceType_CURTAIN = 43499,
    DeviceType_FAN = 43531,
    DeviceType_GATEWAY = 43500,
    DeviceType_GATEWAY_AC = 43756,
    DeviceType_GATEWAY_IHG = 10058,
    DeviceType_LIGHT_BURDEN_CW = 43754,
    DeviceType_LIGHT_BURDEN_W = 43759,
    DeviceType_LIGHT_CCT = 43051,
    DeviceType_LIGHT_COMPOSE = 43709,
    DeviceType_LIGHT_PWR = 43049,
    DeviceType_LIGHT_RGB = 43168,
    DeviceType_LIGHT_RGBCW = 43050,
    DeviceType_LIGHT_RGBW = 43169,
    DeviceType_LIGHT_W_CW = 43745,
    DeviceType_META_PAD = 43518,
    DeviceType_META_PAD_2 = 43974,
    DeviceType_PANEL_3 = 43463,
    DeviceType_PANEL_3_WIRELESS = 43462,
    DeviceType_PANEL_4 = 43473,
    DeviceType_PANEL_4_WIRELESS = 43472,
    DeviceType_PANEL_6 = 43461,
    DeviceType_PANEL_6_WIRELESS = 43459,
    DeviceType_PANEL_8 = 43733,
    DeviceType_PANEL_8_2 = 43734,
    DeviceType_RELAY_1 = 43525,
    DeviceType_RELAY_2 = 43474,
    DeviceType_RELAY_4 = 43680,
    DeviceType_SENSOR_DOOR = 43505,
    DeviceType_SENSOR_IR = 43516,
    DeviceType_SENSOR_RADAR = 43808,
    DeviceType_SENSOR_WATER = 43791,
    DeviceType_THERMOSTAT = 43919,
}
