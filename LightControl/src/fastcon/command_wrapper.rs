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

fn send_single_control(addr: u32, data: Vec<u8>) {
    todo!("send_single_control")
}

pub fn single_on_off_command(short_addr: u32, on: bool) {
    println!(
        "single_on_off_command: short_addr: {:04x}, on: {}",
        short_addr, on
    );

    send_single_control(short_addr, generate_on_off_command(on));
}
