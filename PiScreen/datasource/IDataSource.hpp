#pragma once

#include <map>
#include <string>

namespace PiScreen::datasource
{
    class IDataSource
    {
      public:
        virtual ~IDataSource() = default;
        virtual std::string getData(const std::string &extInfo) = 0;
    };

    inline std::map<int, IDataSource *(*) (void)> registrations;
} // namespace PiScreen::datasource
