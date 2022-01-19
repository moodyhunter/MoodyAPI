#include "AppCore.hpp"
#include "AppSettings.hpp"
#include "JNICall.hpp"

#include <QFontDatabase>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QUrl>

#ifdef Q_OS_ANDROID
constexpr auto PlatformHoverEnabled = false;
#else
constexpr auto PlatformHoverEnabled = true;
#endif

int main(int argc, char *argv[])
{
    QCoreApplication::setApplicationName(u"MoodyApp"_qs);
    QCoreApplication::setOrganizationName(u"Moody"_qs);

#ifdef Q_OS_ANDROID
    if (argc > 1 && strcmp(argv[1], "-service") == 0)
        return runQtAndroid(argc, argv);
#endif

    QGuiApplication::setApplicationDisplayName(u"Moody App"_qs);
    QGuiApplication app(argc, argv);
    global_AppSettings = new AppSettings(&app);

#ifdef Q_OS_ANDROID
    startAndroidService();
#else
    setupRemoteObject();
    runRemoteObject();
#endif

    QQmlApplicationEngine engine;

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
