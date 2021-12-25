#include "AppCore.hpp"
#include "AppSettings.hpp"

#include <QFontDatabase>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QSslSocket>
#include <QUrl>

#ifdef Q_OS_ANDROID
constexpr auto PlatformHoverEnabled = false;
#else
constexpr auto PlatformHoverEnabled = true;
#endif

#include "MoodyApi.grpc.pb.h"

#include <grpcpp/grpcpp.h>

int main(int argc, char *argv[])
{
    QGuiApplication::setApplicationDisplayName(u"Moody App"_qs);
    QGuiApplication::setApplicationName(u"Moody App"_qs);

    {
        const auto chan = grpc::CreateChannel("localhost:1920", grpc::InsecureChannelCredentials());
        auto stub = MoodyAPI::CameraService::NewStub(chan);
        grpc::ClientContext ctx;
        auto reader = stub->SubscribeCameraStateChange(&ctx, {});
        MoodyAPI::CameraStateChangedResponses resp;
        while (reader->Read(&resp))
        {
            qDebug() << resp.states_size();
            for (auto i = 0; i < resp.states_size(); i++)
            {
                const auto state = resp.states(i);
                switch (state.values_case())
                {
                    case MoodyAPI::CameraState::kNewState: qDebug() << state.newstate(); break;
                    case MoodyAPI::CameraState::kIp4Address: qDebug() << QString::fromStdString(state.ip4address()); break;
                    case MoodyAPI::CameraState::kIp6Address: qDebug() << QString::fromStdString(state.ip6address()); break;
                    case MoodyAPI::CameraState::kMotionEventId: qDebug() << QString::fromStdString(state.motioneventid()); break;
                    case MoodyAPI::CameraState::VALUES_NOT_SET: qDebug() << "???"; break;
                }
            }
        }
    }

    QGuiApplication app(argc, argv);

    qDebug() << "Device supports OpenSSL: " << QSslSocket::supportsSsl();
    qDebug() << "Qt SSL Backends: " << QSslSocket::availableBackends();
    qDebug() << "Qt SSL Active Backend: " << QSslSocket::activeBackend();

    QQmlApplicationEngine engine;

    MoodyAppSettings = new AppSettings(&app);

    qmlRegisterSingletonInstance<AppCore>("client.api.mooody.me", 1, 0, "AppCore", new AppCore(&app));
    qmlRegisterSingletonInstance<AppSettings>("client.api.mooody.me", 1, 0, "AppSettings", MoodyAppSettings);

    engine.rootContext()->setContextProperty(u"fixedFont"_qs, QFontDatabase::systemFont(QFontDatabase::FixedFont));
    engine.rootContext()->setContextProperty(u"PlatformHoverEnabled"_qs, PlatformHoverEnabled);
    engine.addImportPath(app.applicationDirPath());

    const QUrl url(u"qrc:/client/api/mooody/me/qml/main.qml"_qs);
    const auto callback = [url](const QObject *obj, const QUrl &objUrl)
    {
        if (!obj && url == objUrl)
            QCoreApplication::exit(-1);
    };

    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated, &app, callback, Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
