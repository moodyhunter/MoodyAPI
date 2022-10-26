pub mod common {
    include!("generated/common.rs");
}

pub mod privileged {
    include!("generated/privileged.rs");
}

pub mod notifications {
    include!("generated/notifications.rs");
}

pub mod dns {
    include!("generated/dns.rs");
}

pub mod light {
    include!("generated/light.rs");
}

pub mod moody_api {
    include!("generated/moody_api.rs");
}
