use std::sync::atomic::AtomicBool;

use tonic::transport::Channel;

pub struct GlobalState {
    pub api_secret: String,
    pub camera_state: AtomicBool,
    pub channel: Channel,
}
