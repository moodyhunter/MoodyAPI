#include "ServerConnection.hpp"

#include "AppSettings.hpp"
#include "CameraAPI.grpc.pb.h"

#include <QMutex>
#include <QThread>
#include <QtDebug>
#include <grpcpp/grpcpp.h>

using namespace std::chrono_literals;

ServerConnection::ServerConnection(const QString &serverAddress, const QString &secret)
{
    m_serverAddr = serverAddress.trimmed();
    m_secret = secret.trimmed();
    m_isRunning = true;
}

ServerConnection::~ServerConnection()
{
}

void ServerConnection::run()
{
    emit onConnectionStatusChanged(false);
    static QMutex m;

    if (m_pollingContext)
        m_pollingContext->TryCancel();

    if (m_channel)
        m_channel.reset();

    m.lock();

    while (m_isRunning)
    {
        m_channel = grpc::CreateChannel(m_serverAddr.toStdString(), global_AppSettings->getDisableTLS() ? grpc::InsecureChannelCredentials() : grpc::SslCredentials({}));

        while (m_isServerConnected = m_channel->WaitForConnected(std::chrono::system_clock::now() + 1s), m_isRunning && !m_isServerConnected)
            qDebug() << "Server not connected, retry.";

        if (!m_isRunning)
        {
            qDebug() << "Leaving";
            m.unlock();
            return;
        }

        auto serverStub = CameraAPI::CameraService::NewStub(m_channel);

        while (m_isRunning)
        {
            m_pollingContext.reset(new grpc::ClientContext);

            CameraAPI::SubscribeCameraStateChangeRequest request;
            request.mutable_auth()->set_secret(m_secret.toStdString());

            auto reader = serverStub->SubscribeCameraStateChange(m_pollingContext.get(), request);

            CameraAPI::CameraState resp;
            while (reader->Read(&resp) && m_isRunning)
            {
                emit onConnectionStatusChanged(true);
                if (resp.has_newstate())
                {
                    std::cout << std::boolalpha << resp.newstate() << std::endl;
                    emit onCameraStateChanged(resp.newstate());
                }

                if (resp.has_ip4address())
                    std::cout << resp.ip4address() << std::endl;
                if (resp.has_ip6address())
                    std::cout << resp.ip6address() << std::endl;
                if (resp.has_motioneventid())
                    std::cout << resp.motioneventid() << std::endl;
                if (resp.has_imagepng())
                    emit onNewMotionDetected(QByteArray::fromStdString(resp.imagepng()));
            }
            emit onConnectionStatusChanged(false);

            qDebug() << "Cannot read more responses, retry.";
            m_pollingContext->TryCancel();
            QThread::sleep(1);
        }
    }
    m.unlock();
}

void ServerConnection::StopPolling()
{
    m_isRunning = false;
    if (m_pollingContext)
        m_pollingContext->TryCancel();
}

void ServerConnection::SetCameraState(bool newState)
{
    if (!m_isServerConnected)
        return;

    grpc::ClientContext m_pollingContext;
    auto serverStub = CameraAPI::CameraService::NewStub(m_channel);

    CameraAPI::SetCameraStateRequest request;
    request.mutable_auth()->set_secret(m_secret.toStdString());
    request.mutable_state()->set_newstate(newState);

    ::google::protobuf::Empty empty;
    serverStub->SetCameraState(&m_pollingContext, request, &empty);
}
