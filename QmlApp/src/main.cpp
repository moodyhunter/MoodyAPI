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

int main(int argc, char *argv[])
{
    QGuiApplication::setApplicationDisplayName(u"Moody App"_qs);
    QGuiApplication::setApplicationName(u"Moody App"_qs);

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
