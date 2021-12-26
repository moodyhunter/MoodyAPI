#include "ServerConnection.hpp"

#include "CameraAPI.grpc.pb.h"

#include <QMutex>
#include <QThread>
#include <QtDebug>
#include <grpcpp/grpcpp.h>

using namespace std::chrono_literals;

ServerConnection::ServerConnection(const QString &serverAddress, const QString &secret)
{
    m_serverAddress = serverAddress.trimmed();
    m_secret = secret.trimmed();
    m_isRunning = true;
}

void ServerConnection::run()
{
    emit onConnectionStatusChanged(false);
    static QMutex m;

    if (m_serverChannel)
    {
        m_serverChannel.reset();
    }

    m.lock();

    while (m_isRunning)
    {
#ifdef QT_DEBUG
        m_serverChannel = grpc::CreateChannel(m_serverAddress.toStdString(), grpc::InsecureChannelCredentials());
#else
        m_serverChannel = grpc::CreateChannel(m_serverAddress.toStdString(), grpc::SslCredentials({}));
#endif

        while (m_isServerConnected = m_serverChannel->WaitForConnected(std::chrono::system_clock::now() + 2s), !m_isServerConnected && m_isRunning)
            qDebug() << "Server not connected, retry.";

        if (!m_isRunning)
        {
            qDebug() << "Leaving";
            m.unlock();
            return;
        }

        auto serverStub = CameraAPI::CameraService::NewStub(m_serverChannel);

        while (m_isRunning)
        {
            grpc::ClientContext m_pollingContext;
            CameraAPI::SubscribeCameraStateChangeRequest request;
            request.mutable_auth()->set_secret(m_secret.toStdString());

            auto reader = serverStub->SubscribeCameraStateChange(&m_pollingContext, request);
            CameraAPI::CameraStateChangedResponse resp;
            while (reader->Read(&resp) && m_isRunning)
            {
                emit onConnectionStatusChanged(true);
                switch (resp.values_case())
                {
                    case CameraAPI::CameraStateChangedResponse::kNewState:
                    {
                        std::cout << std::boolalpha << resp.newstate() << std::endl;
                        emit onCameraStateChanged(resp.newstate());
                        break;
                    }
                    case CameraAPI::CameraStateChangedResponse::kIp4Address:
                    {
                        std::cout << resp.ip4address() << std::endl;
                        break;
                    }
                    case CameraAPI::CameraStateChangedResponse::kIp6Address:
                    {
                        std::cout << resp.ip6address() << std::endl;
                        break;
                    }
                    case CameraAPI::CameraStateChangedResponse::kMotionEventId:
                    {
                        std::cout << resp.motioneventid() << std::endl;
                        break;
                    }
                    case CameraAPI::CameraStateChangedResponse::VALUES_NOT_SET:
                    {
                        std::cout << "???" << std::endl;
                        break;
                    }
                }
            }
            emit onConnectionStatusChanged(false);

            qDebug() << "Cannot read more responses, retry.";
            m_pollingContext.TryCancel();
            QThread::sleep(1);
        }
    }
    m.unlock();
}

void ServerConnection::StopPolling()
{
    m_isRunning = false;
}

void ServerConnection::SetCameraState(bool newState)
{
    if (!m_isServerConnected)
        return;

    grpc::ClientContext m_pollingContext;
    auto serverStub = CameraAPI::CameraService::NewStub(m_serverChannel);

    CameraAPI::SetCameraStateRequest request;
    request.mutable_auth()->set_secret(m_secret.toStdString());
    request.mutable_state()->set_newstate(newState);

    ::google::protobuf::Empty empty;
    serverStub->SetCameraState(&m_pollingContext, request, &empty);
}

ServerConnection::~ServerConnection()
{
}
