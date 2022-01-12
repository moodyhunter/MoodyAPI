#pragma once

#include "rep_IPC_merged.h"

#include <QObject>
#include <QRemoteObjectNode>

class ServerConnection;

class AppCore : public QObject
{
    Q_OBJECT

  public:
    explicit AppCore(QObject *parent = nullptr);
    ~AppCore();

    Q_PROPERTY(bool IsRecording MEMBER m_cameraState NOTIFY CameraStatusChanged)
    Q_PROPERTY(bool ServerConnected MEMBER m_connectionStatus NOTIFY ConnectionStatusChanged)

    Q_INVOKABLE void connectToServer(const QString &serverAddress, const QString &secret, bool noTls);
    Q_INVOKABLE void startRecording();
    Q_INVOKABLE void stopRecording();

  signals:
    // QML interop
    void CameraStatusChanged();
    void ConnectionStatusChanged();

  private:
    void m_setCameraState(bool newState);
    void m_setConnectionStatus(bool newState);

  private:
    bool m_connectionStatus = false;
    bool m_cameraState = false;

    QSharedPointer<JNIIPCBridgeReplica> m_source;
    QRemoteObjectNode m_remoteNode;
};
