#include "ServerConnection.hpp"

#include "CameraAPI.grpc.pb.h"

#include <grpcpp/grpcpp.h>
#include <iostream>

gRPCServerConnection::gRPCServerConnection(QObject *parent) : QObject(parent)
{
    return;

    {
        const auto chan = grpc::CreateChannel("localhost:1920", grpc::InsecureChannelCredentials());
        auto stub = CameraAPI::CameraService::NewStub(chan);
        grpc::ClientContext ctx;
        CameraAPI::SubscribeCameraStateChangeRequest request;
        request.mutable_auth()->set_secret("3ebd788d-6187-4f9c-8db7-5cda7d428b91");
        auto reader = stub->SubscribeCameraStateChange(&ctx, request);
        CameraAPI::CameraStateChangedResponse resp;
        while (reader->Read(&resp))
        {
            switch (resp.values_case())
            {
                case CameraAPI::CameraStateChangedResponse::kNewState: std::cout << std::boolalpha << resp.newstate() << std::endl; break;
                case CameraAPI::CameraStateChangedResponse::kIp4Address: std::cout << resp.ip4address() << std::endl; break;
                case CameraAPI::CameraStateChangedResponse::kIp6Address: std::cout << resp.ip6address() << std::endl; break;
                case CameraAPI::CameraStateChangedResponse::kMotionEventId: std::cout << resp.motioneventid() << std::endl; break;
                case CameraAPI::CameraStateChangedResponse::VALUES_NOT_SET: std::cout << "???" << std::endl; break;
            }
        }
    }
}

gRPCServerConnection::~gRPCServerConnection()
{
}
