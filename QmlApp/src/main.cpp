#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QSslSocket>
#include <QUrl>

int main(int argc, char *argv[])
{
    QGuiApplication::setApplicationDisplayName(u"Moody App"_qs);
    QGuiApplication::setApplicationName(u"Moody App"_qs);

    QGuiApplication app(argc, argv);

    qDebug() << "Device supports OpenSSL: " << QSslSocket::supportsSsl();
    qDebug() << "Qt SSL Backends: " << QSslSocket::availableBackends();
    qDebug() << "Qt SSL Active Backend: " << QSslSocket::activeBackend();

    QQmlApplicationEngine engine;
    engine.addImportPath(app.applicationDirPath());

    const QUrl url(u"qrc:/client/api/mooody/me/qml/main.qml"_qs);
    QObject::connect(
        &engine, &QQmlApplicationEngine::objectCreated, &app,
        [url](const QObject *obj, const QUrl &objUrl)
        {
            if (!obj && url == objUrl)
                QCoreApplication::exit(-1);
        },
        Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
