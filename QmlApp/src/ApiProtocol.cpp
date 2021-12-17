#include "ApiProtocol.hpp"

#include <QtDebug>

ApiProtocol::ApiProtocol(QObject *parent) : QObject(parent)
{
}

void ApiProtocol::SetCameraStatus(bool status)
{
    qDebug() << "Camera status: " << status;
    m_cameraStatus = status;
    emit CameraStatusChanged();
}

bool ApiProtocol::GetCameraStatus()
{
    return m_cameraStatus;
}
