#pragma once

#include "ServerConnection.hpp"
#include "rep_IPC_merged.h"

#ifdef Q_OS_ANDROID
#include <QJniObject>
#include <jni.h>
#include <private/qandroidextras_p.h>

extern "C"
{
    JNIEXPORT void JNICALL Java_client_api_mooody_me_QtAndroidService_doWork();
}

int runQtAndroid(int argc, char *argv[]);
void showAndroidNotification(const QString &title, const QString &message);
void startAndroidService();
#endif

void setupRemoteObject();
void runRemoteObject();

inline ServerConnection *m_worker = nullptr;
inline QRemoteObjectHost *host;

class StateSender : public JNIIPCBridgeSimpleSource
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

inline StateSender *source = nullptr;
