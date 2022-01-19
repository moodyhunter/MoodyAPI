#include "JNICall.hpp"

#include "AppSettings.hpp"
#include "rep_IPC_merged.h"

void setupRemoteObject()
{
    host = new QRemoteObjectHost(QUrl(QStringLiteral("localabstract:replica")));
    source = new StateSender;
    host->enableRemoting(source);
}

void runRemoteObject()
{
    if (m_worker)
        return;

    m_worker = new ServerConnection;
    m_worker->SetServerInfo(global_AppSettings->getApiHost(), global_AppSettings->getApiSecret(), global_AppSettings->getDisableTLS());

    QObject::connect(m_worker, &ServerConnection::onCameraStateChanged, source, &JNIIPCBridgeSource::pushRecordingStarted);
    QObject::connect(m_worker, &ServerConnection::onConnectionStatusChanged, source, &JNIIPCBridgeSource::pushAPIServerConnected);
    m_worker->start();
}

#ifdef Q_OS_ANDROID
void Java_client_api_mooody_me_QtAndroidService_doWork()
{
    runRemoteObject();
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
}

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

int runQtAndroid(int argc, char *argv[])
{
    qDebug() << "Service starting with from the same .so file";
    QAndroidService app(argc, argv);
    global_AppSettings = new AppSettings(&app);
    setupRemoteObject();
    return app.exec();
}
#endif
