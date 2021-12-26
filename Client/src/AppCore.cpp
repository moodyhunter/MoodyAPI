#include "AppCore.hpp"

#include "AppSettings.hpp"
#include "ServerConnection.hpp"

#include <QtDebug>

AppCore::AppCore(QObject *parent) : QObject(parent)
{
}

AppCore::~AppCore()
{
    if (m_worker)
    {
        m_worker->StopPolling();
        m_worker->quit();
        m_worker->wait();
        delete m_worker;
    }
}

void AppCore::m_HandleCameraStateChanged(bool newState)
{
    m_cameraState = newState;
    emit CameraStatusChanged();
}

void AppCore::connectToServer(const QString &serverAddress, const QString &secret)
{
    if (m_worker)
    {
        m_worker->StopPolling();
        m_worker->quit();
        m_worker->wait();
        delete m_worker;
    }

    m_worker = new ServerConnection(serverAddress, secret);

    connect(m_worker, &ServerConnection::onCameraStateChanged, this, &AppCore::m_HandleCameraStateChanged);
    connect(m_worker, &ServerConnection::onConnectionStatusChanged, this, [this](bool s) { m_connectionStatus = s, emit ConnectionStatusChanged(s); });

    m_worker->start();
}

void AppCore::SetCameraState(bool status)
{
    m_worker->SetCameraState(status);
}

bool AppCore::GetCameraStatus() const
{
    return m_cameraState;
}
