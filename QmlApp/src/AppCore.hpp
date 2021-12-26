#pragma once

#include <QObject>
#include <QThread>

class ServerConnection;

class AppCore : public QObject
{
    Q_OBJECT

  public:
    explicit AppCore(QObject *parent = nullptr);
    ~AppCore();

    Q_PROPERTY(bool IsRecording READ GetCameraStatus WRITE SetCameraState NOTIFY CameraStatusChanged)
    Q_PROPERTY(bool ServerConnected MEMBER m_connectionStatus NOTIFY ConnectionStatusChanged)

    bool GetCameraStatus() const;
    Q_INVOKABLE void connectToServer(const QString &serverAddress, const QString &secret);

  private slots:
    void m_HandleCameraStateChanged(bool newState);

  signals:
    // To ServerConnection
    void SetCameraState(bool status);

    // QML interop
    void ConnectionStatusChanged(bool connected);
    void CameraStatusChanged();

  private:
    bool m_connectionStatus = false;
    bool m_cameraState = false;
    ServerConnection *m_worker = nullptr;
};
