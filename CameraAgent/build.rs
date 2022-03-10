fn main() {
    tonic_build::compile_protos("../MoodyAPI.proto").expect("Failed to run protoc.");
}
