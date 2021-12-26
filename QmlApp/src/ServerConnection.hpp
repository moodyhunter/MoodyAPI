#pragma once

#include <QThread>
#include <grpcpp/grpcpp.h>

class ServerConnection : public QThread
{
    Q_OBJECT
  public:
    explicit ServerConnection(const QString &serverAddress, const QString &secret);
    virtual ~ServerConnection();

    void run() override;

  public slots:
    void StopPolling();
    void SetCameraState(bool newState);

  signals:
    void onCameraStateChanged(bool newState);
    void onConnectionStatusChanged(bool newState);

  private:
    QString m_serverAddress;
    QString m_secret;
    std::shared_ptr<grpc::Channel> m_serverChannel;

    bool m_isRunning = false;
    bool m_isServerConnected = false;
};
