#[allow(unused)]
pub(crate) enum LightCommand {
    Colored(
        bool, /*on*/
        u8,   /*brightness*/
        u8,
        u8,
        u8,   /*RGB*/
        bool, /*absolute (do not normalize)*/
    ),
    WarmWhite(
        bool, /* on */
        u8,   /* brightness */
        u8,   /* warm val */
        u8,   /* cold val */
    ),
    Brightness(bool /* on */, u8 /* brightness */),
    OnOff(bool /* on */, u8 /* brightness */),
    Some5(bool),                                         // unknown
    Some6(bool, u8),                                     // unknown
    Some7(bool, u8, u8, u8, u8, bool, u8, u8, bool, u8), // unknown
    Some8(bool, u8, bool, u8),                           // unknown
    Some9(u8, u8, u8, bool),                             // unknown
}

impl From<LightCommand> for Vec<u8> {
    fn from(command: LightCommand) -> Self {
        match command {
            LightCommand::Colored(on, brightness, r, g, b, abs) => {
                let mut arr = vec![0; 6];
                let color_normalisation = if abs {
                    1.0
                } else {
                    255.0 / (r as u32 + g as u32 + b as u32) as f32
                };
                arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
                arr[1] = ((b as f32) * color_normalisation) as u32 as u8;
                arr[2] = ((r as f32) * color_normalisation) as u32 as u8;
                arr[3] = ((g as f32) * color_normalisation) as u32 as u8;
                arr[4] = 0;
                arr[5] = 0;
                arr
            }
            LightCommand::WarmWhite(on, brightness, i5, i6) => {
                let mut arr = vec![0; 6];
                arr[0] = ((if on { 128 } else { 0 }) + (brightness & 127)) as u8;
                arr[1] = 0;
                arr[2] = 0;
                arr[3] = 0;
                arr[4] = i5 as u8;
                arr[5] = i6 as u8;
                arr
            }
            LightCommand::Brightness(on, val) => {
                // if maybe_batch {
                //     let mut arr = vec![0; 6];
                //     arr[0] = (if on { val & 127 } else { 0 }) as u8;
                //     arr[1] = 0;
                //     arr[2] = 0;
                //     arr[3] = 0;
                //     arr[4] = 0;
                //     arr[5] = 0;
                //     return arr;
                // }
                vec![if on { val & 127 } else { 0 } as u8]
            }
            LightCommand::OnOff(on, brightness) => {
                // if maybe_batch {
                //     let mut arr = vec![0; 6];
                //     arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
                //     arr[1] = 0;
                //     arr[2] = 0;
                //     arr[3] = 0;
                //     arr[4] = 0;
                //     arr[5] = 0;
                //     return arr;
                // }

                vec![if on { 128 } else { 0 } + (brightness & 127) as u8]
            }
            LightCommand::Some5(some_z) => {
                let mut arr = vec![0; 7];
                arr[0] = 0;
                arr[1] = 0;
                arr[2] = 0;
                arr[3] = 0;
                arr[4] = u8::MAX;
                arr[5] = u8::MAX;
                arr[6] = if some_z { 128 } else { 0 } as u8;
                arr
            }
            LightCommand::Some6(some_z, some_i) => {
                let mut arr = vec![0; 7];
                arr[0] = 0;
                arr[1] = 0;
                arr[2] = 0;
                arr[3] = 0;
                arr[4] = u8::MAX;
                arr[5] = u8::MAX;
                arr[6] = if some_z { 128 } else { 0 } + (some_i & 127) as u8;
                arr
            }
            LightCommand::Some7(on, brightness, r, g, b, absolute, i5, i6, z, i7) => {
                let color_normalisation = if absolute {
                    1.0
                } else {
                    255.0 / (r as u32 + g as u32 + b as u32) as f32
                };
                let mut arr = vec![0; 7];
                arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
                arr[1] = ((r as f32) * color_normalisation) as u32 as u8;
                arr[2] = ((g as f32) * color_normalisation) as u32 as u8;
                arr[3] = ((b as f32) * color_normalisation) as u32 as u8;
                arr[4] = i5 as u8;
                arr[5] = i6 as u8;
                arr[6] = if z { 128 } else { 0 } + (i7 & 127) as u8;
                arr
            }
            LightCommand::Some8(on, brightness, z2, i7) => {
                let mut arr = vec![0; 7];
                arr[0] = (if on { 128 } else { 0 } + (brightness & 127)) as u8;
                arr[1] = u8::MAX;
                arr[2] = u8::MAX;
                arr[3] = u8::MAX;
                arr[4] = u8::MAX;
                arr[5] = u8::MAX;
                arr[6] = if z2 { 128 } else { 0 } + (i7 & 127) as u8;
                arr
            }
            LightCommand::Some9(r, g, b, absolute) => {
                let color_normalisation = if absolute {
                    1.0
                } else {
                    255.0 / (r as u32 + g as u32 + b as u32) as f32
                };
                let mut arr = vec![0; 7];
                arr[0] = u8::MAX;
                arr[1] = ((r as f32) * color_normalisation) as u32 as u8;
                arr[2] = ((g as f32) * color_normalisation) as u32 as u8;
                arr[3] = ((b as f32) * color_normalisation) as u32 as u8;
                arr[4] = u8::MAX;
                arr[5] = u8::MAX;
                arr[6] = u8::MIN;
                arr
            }
        }
    }
}
