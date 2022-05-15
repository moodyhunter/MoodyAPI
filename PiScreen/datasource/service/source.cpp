#include "source.hpp"

#include <cstdlib>
#include <iostream>
#include <sstream>

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
        sd_bus_error_free(&sdBusError);
    }

    std::string SystemdServiceDataSource::getData(const std::string &service)
    {
        // Get service path from systemd
        const char *ServicePath = nullptr;

        {
            const auto errorCodeCall = sd_bus_call_method( //
                sdBus,                                     // The bus connection.
                "org.freedesktop.systemd1",                // The service to call.
                "/org/freedesktop/systemd1",               // The object path.
                "org.freedesktop.systemd1.Manager",        // The interface.
                "LoadUnit",                                // The method to call.
                &sdBusError,                               // The error message.
                &message,                                  // The message.
                "s",                                       // The signature of the parameters.
                service.c_str()                            // The parameters.
            );

            if (errorCodeCall < 0)
            {
                std::cerr << "Failed to call method, error: " << strerror(-errorCodeCall);
                std::exit(1);
            }

            const auto errorCodeReadMsg = sd_bus_message_read(message, "o", &ServicePath);
            if (errorCodeReadMsg < 0)
            {
                std::cerr << "Failed to read message, error: " << strerror(-errorCodeReadMsg);
                std::exit(1);
            }
            sd_bus_message_unref(message);
            message = nullptr;
        }

        char *ServiceStatus = nullptr;

        {
            const auto errorCodeGetProp = sd_bus_get_property_string( //
                sdBus,                                                // The bus connection.
                "org.freedesktop.systemd1",                           // The service to call.
                ServicePath,                                          // The object path.
                "org.freedesktop.systemd1.Unit",                      // The interface.
                "ActiveState",                                        // The property.
                &sdBusError,                                          // The error message.
                &ServiceStatus                                        // The property value.
            );

            if (errorCodeGetProp < 0)
            {
                std::cerr << "Failed to get property, error: " << strerror(-errorCodeGetProp);
                std::exit(1);
            }
        }

        sd_bus_message_unref(message);
        std::string ret{ ServiceStatus };
        std::free(ServiceStatus);

        return ret == "active" ? "ON" : "OFF";
    }
} // namespace PiScreen::datasource