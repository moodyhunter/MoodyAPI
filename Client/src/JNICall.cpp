#include "JNICall.hpp"

#include "AppSettings.hpp"
#include "ServerConnection.hpp"
#include "rep_IPC_merged.h"

#include <QJniObject>
#include <private/qandroidextras_p.h>

void showAndroidNotification(const QString &title, const QString &message)
{
    QJniObject nTitle = QJniObject::fromString(title);
    QJniObject nContent = QJniObject::fromString(message);
    QJniObject::callStaticMethod<void>("client/api/mooody/me/AndroidUtils",                                //
                                       "notify",                                                           //
                                       "(Landroid/content/Context;Ljava/lang/String;Ljava/lang/String;)V", //
                                       QNativeInterface::QAndroidApplication::context(),                   //
                                       nTitle.object<jstring>(),                                           //
                                       nContent.object<jstring>());
}

void startAndroidService()
{
    QJniObject::callStaticMethod<void>("client/api/mooody/me/AndroidUtils", //
                                       "startService",                      //
                                       "(Landroid/content/Context;)V",      //
                                       QNativeInterface::QAndroidApplication::context());
}

ServerConnection *m_worker = nullptr;

class Source : public JNIIPCBridgeSimpleSource
{
  public slots:
    virtual void StartRecord() override
    {
        m_worker->SetCameraState(true);
    }
    virtual void StopRecord() override
    {
        m_worker->SetCameraState(false);
    }
    virtual void SetServerInfo(QString host, QString secret, bool noTls) override
    {
        m_worker->SetServerInfo(host, secret, noTls);
    }
};

QRemoteObjectHost *host;
Source *source = nullptr;

int runQtAndroid(int argc, char *argv[])
{
    qDebug() << "Service starting with from the same .so file";
    QAndroidService app(argc, argv);
    host = new QRemoteObjectHost(QUrl(QStringLiteral("localabstract:replica")));
    source = new Source;
    host->enableRemoting(source);
    return app.exec();
}

void Java_client_api_mooody_me_QtAndroidService_doWork()
{
    if (m_worker)
        return;

    m_worker = new ServerConnection;
    {
        AppSettings settings;
        m_worker->SetServerInfo(settings.getApiHost(), settings.getApiSecret(), settings.getDisableTLS());
    }

    QObject::connect(m_worker, &ServerConnection::onCameraStateChanged, source, &JNIIPCBridgeSource::pushRecordingStarted);
    QObject::connect(m_worker, &ServerConnection::onConnectionStatusChanged, source, &JNIIPCBridgeSource::pushAPIServerConnected);

    static const auto title = u"Camera State Changed"_qs;
    static const auto content = u"The camera has been turned %1"_qs;
    static bool previousState = false;
    QObject::connect(m_worker, &ServerConnection::onCameraStateChanged,
                     [](bool state)
                     {
                         if (state != previousState)
                             showAndroidNotification(title, content.arg(state ? "on" : "off"));
                         previousState = state;
                     });

    m_worker->start();
}
