fn main() {
    println!("cargo:rerun-if-changed=../MoodyAPI.proto");
    tonic_build::compile_protos("../MoodyAPI.proto").unwrap();
}
