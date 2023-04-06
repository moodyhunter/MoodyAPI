#include "source.hpp"

#include <cstdlib>
#include <iostream>
#include <sstream>
#include <systemd/sd-bus.h>

std::string GetServicePath(sd_bus *sdBus, const std::string &serviceName)
{
    if (serviceName == "")
        return "";

    const char *servicePath = nullptr;
    sd_bus_error sdBusError = SD_BUS_ERROR_NULL;
    sd_bus_message *message = nullptr;
    const auto errorCodeCall = sd_bus_call_method( //
        sdBus,                                     // The bus connection.
        "org.freedesktop.systemd1",                // The service to call.
        "/org/freedesktop/systemd1",               // The object path.
        "org.freedesktop.systemd1.Manager",        // The interface.
        "LoadUnit",                                // The method to call.
        &sdBusError,                               // The error message.
        &message,                                  // The message.
        "s",                                       // The signature of the parameters.
        serviceName.c_str()                        // The parameters.
    );

    if (errorCodeCall < 0)
    {
        std::cerr << "Failed to call method, error: " << strerror(-errorCodeCall);
        return "";
    }

    const auto errorCodeReadMsg = sd_bus_message_read(message, "o", &servicePath);
    if (errorCodeReadMsg < 0)
    {
        std::cerr << "Failed to read message, error: " << strerror(-errorCodeReadMsg);
        return "";
    }

    sd_bus_error_free(&sdBusError);
    sd_bus_message_unref(message);
    return servicePath;
}

std::string GetServiceActiveState(sd_bus *sdBus, const std::string &servicePath)
{
    if (servicePath.empty())
        return "";

    sd_bus_error sdBusError = SD_BUS_ERROR_NULL;
    char *status = nullptr;
    const auto errorCodeGetProp = sd_bus_get_property_string( //
        sdBus,                                                // The bus connection.
        "org.freedesktop.systemd1",                           // The service to call.
        servicePath.c_str(),                                  // The object path.
        "org.freedesktop.systemd1.Unit",                      // The interface.
        "ActiveState",                                        // The property.
        &sdBusError,                                          // The error message.
        &status                                               // The property value.
    );

    if (errorCodeGetProp < 0)
    {
        std::cerr << "Failed to get property, error: " << strerror(-errorCodeGetProp) << status << servicePath;
        return "";
    }

    sd_bus_error_free(&sdBusError);
    std::string ret{ status };
    free(status);
    return ret;
}

namespace PiScreen::datasource
{
    SystemdServiceDataSource::SystemdServiceDataSource()
    {
        const auto errorCode = sd_bus_default_system(&sdBus);
        if (errorCode < 0)
        {
            std::cerr << "Failed to connect to system bus, error: " << strerror(-errorCode);
            std::exit(1);
        }
    }

    SystemdServiceDataSource::~SystemdServiceDataSource()
    {
        sd_bus_unref(sdBus);
    }

    std::string SystemdServiceDataSource::getData(const std::string &service)
    {
        const auto servicePath = GetServicePath(sdBus, service);
        const auto activeState = GetServiceActiveState(sdBus, servicePath);
        return activeState == "active" ? "ON" : "OFF";
    }
} // namespace PiScreen::datasource
