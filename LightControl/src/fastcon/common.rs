pub fn bytes_to_string(bytes: &[u8]) -> String {
    let mut s = String::new();
    for b in bytes {
        s.push_str(&format!("{:02x}", b));
    }
    s
}

pub fn print_bytes(name: &str, bytes: &[u8]) {
    println!("{}: {}", name, bytes_to_string(bytes));
}
