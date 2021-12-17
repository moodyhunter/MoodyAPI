#pragma once

#include <QObject>

class ApiProtocol : public QObject
{
    Q_OBJECT

  public:
    explicit ApiProtocol(QObject *parent = nullptr);

    Q_PROPERTY(bool CameraStatus READ GetCameraStatus WRITE SetCameraStatus NOTIFY CameraStatusChanged)

  private slots:
    void SetCameraStatus(bool status);
    bool GetCameraStatus();

  signals:
    void CameraStatusChanged();

  private:
    bool m_cameraStatus = false;
};
