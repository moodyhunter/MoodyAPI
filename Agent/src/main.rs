use camera_api::{
    camera_service_client::CameraServiceClient,
    {Auth, SubscribeCameraStateChangeRequest},
};

use tonic::{transport::Channel, Request};

mod camera_api {
    tonic::include_proto!("camera_api");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let api_host = std::env::args().nth(1).expect("no pattern given");
    let api_secret = std::env::args().nth(2).expect("no path given");

    let channel = Channel::from_shared(api_host)?
        .connect()
        .await
        .expect("Can't create a channel");

    let mut client = CameraServiceClient::new(channel);

    let request = Request::new(SubscribeCameraStateChangeRequest {
        auth: Some(Auth {
            secret: api_secret.clone(),
        }),
    });

    match client.subscribe_camera_state_change(request).await {
        Ok(stream) => {
            let mut resp = stream.into_inner();
            while let Some(feature) = resp.message().await? {
                println!("NOTE = {:?}", feature);
            }
        }
        Err(e) => println!("something went wrong: {:?}", e),
    }

    Ok(())
}
