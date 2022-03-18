#pragma once

#include <QThread>
#include <grpcpp/grpcpp.h>

class ServerConnection : public QThread
{
    Q_OBJECT
  public:
    explicit ServerConnection();
    virtual ~ServerConnection();

    void run() override;

  public slots:
    void SetCameraState(bool newState);
    void SetServerInfo(const QString &host, const QString &secret, bool noTls);

  signals:
    void onCameraStateChanged(bool newState);
    void onConnectionStatusChanged(bool newState);

  private:
    QString m_serverAddr;
    QString m_secret;
    bool m_noTls;

    std::unique_ptr<grpc::ClientContext> m_pollingContext = nullptr;
    std::shared_ptr<grpc::Channel> m_channel;

    bool m_isServerConnected = false;
    QAtomicInteger<bool> m_serverInfoChanged = false;
};
