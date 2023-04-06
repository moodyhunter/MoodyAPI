#include "datasource/IDataSource.hpp"

namespace PiScreen::datasource
{
    class DateTimeDateSource : public IDataSource
    {
      public:
        DateTimeDateSource();
        ~DateTimeDateSource();

        std::string getData(const std::string &extInfo) override;
    };

} // namespace PiScreen::datasource
