use std::env;

use notify_rust::{Hint, Notification};
use platform_dirs::AppDirs;

fn main() {
    if env::args().len() > 1 {
        print!("Has an argument, sending messages...\n")
    }

    dbg!(AppDirs::new(Some("MyApp"), false));

    Notification::new()
        .summary("Hello, World")
        .body("Just to tell you that I'm up.")
        .icon("flag-green")
        .appname("Notify Client")
        .hint(Hint::Resident(true))
        .show()
        .unwrap();
}
