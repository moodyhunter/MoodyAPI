#include "AppCore.hpp"

#include <QtDebug>

AppCore::AppCore(QObject *parent) : QObject(parent)
{
}

void AppCore::SetCameraStatus(bool status)
{
    qDebug() << "Camera status: " << status;
    m_cameraStatus = status;
    emit CameraStatusChanged();
}

bool AppCore::GetCameraStatus() const
{
    return m_cameraStatus;
}
