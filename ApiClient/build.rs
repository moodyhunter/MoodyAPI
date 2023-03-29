fn main() -> Result<(), std::io::Error> {
    println!("cargo:rerun-if-changed=../proto/");
    println!("cargo:rerun-if-changed=models/generated");

    tonic_build::configure()
        .build_server(false)
        .out_dir("src/models/generated")
        .compile(&["../proto/MoodyAPI.proto"], &["../proto"])
}
