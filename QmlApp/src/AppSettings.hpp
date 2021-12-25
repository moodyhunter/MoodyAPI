#pragma once

#include <QObject>
#include <QSettings>

#define MoodyApp_Q_PROPERTY_decl(type, name, Name, conv)                                                                                                                 \
    Q_PROPERTY(type name READ get##Name WRITE set##Name NOTIFY on##Name##Changed) public : type get##Name() const                                                        \
    {                                                                                                                                                                    \
        return m_settings->value(QStringLiteral(#Name)).conv;                                                                                                            \
    }                                                                                                                                                                    \
    void set##Name(const type &newValue)                                                                                                                                 \
    {                                                                                                                                                                    \
        m_settings->setValue(QStringLiteral(#Name), newValue);                                                                                                           \
        Q_EMIT on##Name##Changed(newValue);                                                                                                                              \
    }                                                                                                                                                                    \
    Q_SIGNAL void on##Name##Changed(const type &v)

class AppSettings : public QObject
{
    Q_OBJECT

  public:
    AppSettings(QObject *parent = nullptr);
    virtual ~AppSettings();

    MoodyApp_Q_PROPERTY_decl(bool, darkMode, DarkMode, toBool());
    MoodyApp_Q_PROPERTY_decl(QString, apiHost, ApiHost, toString());
    MoodyApp_Q_PROPERTY_decl(QString, apiSecret, ApiSecret, toString());

  private:
    static inline QSettings *m_settings = nullptr;
};

inline AppSettings *global_AppSettings;
