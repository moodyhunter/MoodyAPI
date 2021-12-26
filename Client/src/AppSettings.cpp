#include "AppSettings.hpp"

AppSettings::AppSettings(QObject *parent) : QObject(parent)
{
    if (m_settings)
        throw "This should not happen";
    m_settings = new QSettings;
}

AppSettings::~AppSettings()
{
}
