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

bool AppCore::GetCameraStatus()
{
    return m_cameraStatus;
}

void AppCore::SetDarkModeStatus(bool status)
{
    qDebug() << "Darkmode status: " << status;
    m_darkmodeStatus = status;
    emit DarkModeChanged();
}

bool AppCore::DarkModeStatus()
{
    return m_darkmodeStatus;
}
