pub(crate) fn reverse_8(d: u8) -> u8 {
    let mut result = 0;
    for i in 0..8 {
        result |= ((d >> i) & 1) << (7 - i);
    }
    result
}

pub(crate) fn reverse_16(d: u16) -> u16 {
    let mut result = 0;
    for i in 0..16 {
        result |= ((d >> i) & 1) << (15 - i);
    }
    result
}

pub(crate) fn crc16(addr: &[u8], data: &[u8]) -> u16 {
    let mut crc = 0xffff;

    // iterate over address in reverse
    for i in addr.iter().rev() {
        crc ^= (*i as u16) << 8;
        for _ in 0..4 {
            let mut tmp = crc << 1;

            if crc & 0x8000 != 0 {
                tmp ^= 0x1021;
            }

            crc = tmp << 1;
            if tmp & 0x8000 != 0 {
                crc ^= 0x1021;
            }
        }
    }

    for i in 0..data.len() {
        crc ^= (reverse_8(data[i]) as u16) << 8;
        for _ in 0..4 {
            let mut tmp = crc << 1;

            if crc & 0x8000 != 0 {
                tmp ^= 0x1021;
            }

            crc = tmp << 1;
            if tmp & 0x8000 != 0 {
                crc ^= 0x1021;
            }
        }
    }

    crc = !reverse_16(crc);
    crc
}

#[derive(Debug, Clone, Copy, Default)]
pub(crate) struct WhiteningContext {
    f_0x0: u32,
    f_0x4: u32,
    f_0x8: u32,
    f_0xc: u32,
    f_0x10: u32,
    f_0x14: u32,
    f_0x18: u32,
}

impl WhiteningContext {
    pub(crate) fn new() -> Self {
        Self::default()
    }
}

pub(crate) fn whitening_init(val: u32, ctx: &mut WhiteningContext) {
    let v0 = [(val >> 5), (val >> 4), (val >> 3), (val >> 2)];

    ctx.f_0x0 = 1;
    ctx.f_0x4 = v0[0] & 1;
    ctx.f_0x8 = v0[1] & 1;
    ctx.f_0xc = v0[2] & 1;
    ctx.f_0x10 = v0[3] & 1;
    ctx.f_0x14 = (val >> 1) & 1;
    ctx.f_0x18 = val & 1;
}

pub(crate) fn whitening_encode(data: &mut Vec<u8>, ctx: &mut WhiteningContext) {
    for i in 0..data.len() {
        let varC = ctx.f_0xc;
        let var14 = ctx.f_0x14;
        let var18 = ctx.f_0x18;
        let var10 = ctx.f_0x10;
        let var8 = var14 ^ ctx.f_0x8;
        let var4 = var10 ^ ctx.f_0x4;
        let _var = var18 ^ varC;
        let var0 = _var ^ ctx.f_0x0;

        let c = data[i];
        data[i] = ((c & 0x80) ^ ((var8 ^ var18) << 7) as u8)
            + ((c & 0x40) ^ (var0 << 6) as u8)
            + ((c & 0x20) ^ (var4 << 5) as u8)
            + ((c & 0x10) ^ (var8 << 4) as u8)
            + ((c & 0x08) ^ (_var << 3) as u8)
            + ((c & 0x04) ^ (var10 << 2) as u8)
            + ((c & 0x02) ^ (var14 << 1) as u8)
            + ((c & 0x01) ^ (var18 << 0) as u8);

        ctx.f_0x8 = var4;
        ctx.f_0xc = var8;
        ctx.f_0x10 = var8 ^ varC;
        ctx.f_0x14 = var0 ^ var10;
        ctx.f_0x18 = var4 ^ var14;
        ctx.f_0x0 = var8 ^ var18;
        ctx.f_0x4 = var0;
    }
}
