use super::command_wrapper::single_control_nodelay;

mod light_impl;

pub enum LightState {
    Off,
    WarmWhite,       // brightness
    RGB(u8, u8, u8), // brightness, r, g, b
}

pub struct BLELight {
    pub addr: u32,
    pub key: Vec<u8>,
    pub state: LightState,
    pub brightness: u8,
}

impl From<&BLELight> for light_impl::LightCommand {
    fn from(light: &BLELight) -> Self {
        match light.state {
            LightState::Off => light_impl::LightCommand::OnOff(false, 0),
            LightState::WarmWhite => {
                light_impl::LightCommand::WarmWhite(true, light.brightness, 127, 127)
            }
            LightState::RGB(r, g, b) => {
                light_impl::LightCommand::Colored(true, light.brightness, r, g, b, false)
            }
        }
    }
}

impl BLELight {
    pub fn new(addr: u32, key: &[u8]) -> Self {
        BLELight {
            addr,
            key: key.to_vec(),
            state: LightState::Off,
            brightness: 0,
        }
    }

    pub fn set_brightness(&mut self, brightness: u8) {
        self.brightness = brightness.min(0).max(127);
    }

    pub fn set_state(&mut self, state: LightState) {
        self.state = state;
    }

    fn get_command(&self) -> light_impl::LightCommand {
        light_impl::LightCommand::from(self)
    }

    pub fn get_advertisement(&self) -> Vec<u8> {
        let raw = Vec::<u8>::from(self.get_command());
        single_control_nodelay(self.addr, Some(&self.key), raw)
    }
}
