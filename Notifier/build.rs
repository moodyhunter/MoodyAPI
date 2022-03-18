fn main() {
    tonic_build::compile_protos("../MoodyAPI.proto").unwrap();
}
