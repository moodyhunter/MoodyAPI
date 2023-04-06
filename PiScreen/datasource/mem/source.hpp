#pragma once

#include "datasource/IDataSource.hpp"

namespace PiScreen::datasource
{
    class MemoryDataSource : public IDataSource
    {
      public:
        MemoryDataSource();
        ~MemoryDataSource();

        virtual std::string getData(const std::string &extInfo) override;
    };
} // namespace PiScreen::datasource
