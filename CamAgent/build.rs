fn main() {
    tonic_build::compile_protos("../CameraAPI.proto").expect("Failed to run protoc");
}
