#pragma once

#include <string>

namespace PiScreen::config
{
    class IDataSource
    {
      public:
        virtual ~IDataSource() = default;
        virtual std::string getData(const std::string &sourceId) = 0;
    };
} // namespace PiScreen::config
