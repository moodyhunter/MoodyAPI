#pragma once

#include <QObject>

class gRPCServerConnection : public QObject
{
    Q_OBJECT
  public:
    explicit gRPCServerConnection(QObject *parent);
    virtual ~gRPCServerConnection();
};
