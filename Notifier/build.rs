fn main() -> Result<(), std::io::Error> {
    println!("cargo:rerun-if-changed=../proto/");

    tonic_build::configure()
        .build_server(false)
        .out_dir("src/models/generated")
        .compile(&["../proto/MoodyAPI.proto"], &["../proto"])
}
