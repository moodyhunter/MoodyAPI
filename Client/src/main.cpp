#include "AppCore.hpp"
#include "AppSettings.hpp"

#include <QFontDatabase>
#include <QGuiApplication>
#include <QJniEnvironment>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QSslSocket>
#include <QUrl>
#include <private/qandroidextras_p.h>

#ifdef Q_OS_ANDROID
constexpr auto PlatformHoverEnabled = false;
#else
constexpr auto PlatformHoverEnabled = true;
#endif

int runQtAndroid(int argc, char *argv[]);
void startAndroidService();

int main(int argc, char *argv[])
{
    QCoreApplication::setApplicationName(u"MoodyApp"_qs);
    QCoreApplication::setOrganizationName(u"Moody"_qs);

    if (argc > 1 && strcmp(argv[1], "-service") == 0)
    {
        return runQtAndroid(argc, argv);
    }

    QGuiApplication::setApplicationDisplayName(u"Moody App"_qs);
    QGuiApplication app(argc, argv);

    qDebug() << "Device supports OpenSSL: " << QSslSocket::supportsSsl();
    qDebug() << "Qt SSL Backends: " << QSslSocket::availableBackends();
    qDebug() << "Qt SSL Active Backend: " << QSslSocket::activeBackend();

    startAndroidService();

    QQmlApplicationEngine engine;

    global_AppSettings = new AppSettings(&app);

    qmlRegisterSingletonInstance<AppCore>("client.api.mooody.me", 1, 0, "AppCore", new AppCore(&app));
    qmlRegisterSingletonInstance<AppSettings>("client.api.mooody.me", 1, 0, "AppSettings", global_AppSettings);

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
