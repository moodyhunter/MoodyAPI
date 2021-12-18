#pragma once

#include <QObject>

class AppCore : public QObject
{
    Q_OBJECT

  public:
    explicit AppCore(QObject *parent = nullptr);

    Q_PROPERTY(bool CameraStatus READ GetCameraStatus WRITE SetCameraStatus NOTIFY CameraStatusChanged)

  private:
    bool GetCameraStatus() const;

  private slots:
    void SetCameraStatus(bool status);

  signals:
    void CameraStatusChanged();

  private:
    bool m_cameraStatus = false;
};
