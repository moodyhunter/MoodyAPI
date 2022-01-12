#include "AppCore.hpp"

#include "AppSettings.hpp"

AppCore::AppCore(QObject *parent) : QObject(parent)
{
    m_remoteNode.connectToNode(QUrl(QStringLiteral("localabstract:replica")));
    m_source.reset(m_remoteNode.acquire<JNIIPCBridgeReplica>());
    bool res = m_source->waitForSource();
    Q_ASSERT(res);
    connect(m_source.get(), &JNIIPCBridgeReplica::RecordingStartedChanged, this, &AppCore::m_setCameraState);
    connect(m_source.get(), &JNIIPCBridgeReplica::APIServerConnectedChanged, this, &AppCore::m_setConnectionStatus);
}

AppCore::~AppCore()
{
}

void AppCore::m_setCameraState(bool newState)
{
    m_cameraState = newState;
    emit CameraStatusChanged();
}

void AppCore::m_setConnectionStatus(bool newState)
{
    m_connectionStatus = newState;
    emit ConnectionStatusChanged();
}

void AppCore::connectToServer(const QString &serverAddress, const QString &secret, bool noTls)
{
    m_source->SetServerInfo(serverAddress, secret, noTls);
}

void AppCore::startRecording()
{
    m_source->StartRecord();
}

void AppCore::stopRecording()
{
    m_source->StopRecord();
}
