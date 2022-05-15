#include "datasource/IDataSource.hpp"

#include <systemd/sd-bus.h>

namespace PiScreen::datasource
{
    class SystemdServiceDataSource : public IDataSource
    {
      public:
        SystemdServiceDataSource();
        ~SystemdServiceDataSource();

        std::string getData(const std::string &) override;

      private:
        sd_bus_error sdBusError = SD_BUS_ERROR_NULL;
        sd_bus_message *message = nullptr;
        sd_bus *sdBus = nullptr;
    };

} // namespace PiScreen::datasource