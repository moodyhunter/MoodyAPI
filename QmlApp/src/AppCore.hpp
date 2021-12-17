#pragma once

#include <QObject>

class AppCore : public QObject
{
    Q_OBJECT

  public:
    explicit AppCore(QObject *parent = nullptr);

    Q_PROPERTY(bool CameraStatus READ GetCameraStatus WRITE SetCameraStatus NOTIFY CameraStatusChanged)
    Q_PROPERTY(bool DarkMode READ DarkModeStatus WRITE SetDarkModeStatus NOTIFY DarkModeChanged)

  private slots:
    void SetCameraStatus(bool status);
    bool GetCameraStatus();

    void SetDarkModeStatus(bool status);
    bool DarkModeStatus();

  signals:
    void CameraStatusChanged();
    void DarkModeChanged();

  private:
    bool m_cameraStatus = false;
    bool m_darkmodeStatus = false;
};
