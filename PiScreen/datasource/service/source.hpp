#include "datasource/IDataSource.hpp"

typedef struct sd_bus sd_bus;

namespace PiScreen::datasource
{

    class SystemdServiceDataSource : public IDataSource
    {
      public:
        SystemdServiceDataSource();
        ~SystemdServiceDataSource();

        std::string getData(const std::string &) override;

      private:
        sd_bus *sdBus = nullptr;
    };

} // namespace PiScreen::datasource

const char *GetServicePath(sd_bus *sdBus, const char *serviceName);
const char *GetServiceActiveState(sd_bus *sdBus, const char *servicePath);
