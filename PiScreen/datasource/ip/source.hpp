#pragma once

#include "datasource/IDataSource.hpp"

namespace PiScreen::datasource
{
    class IPAddressDataSource : public IDataSource
    {
      public:
        IPAddressDataSource();
        ~IPAddressDataSource();

        virtual std::string getData(const std::string &extInfo) override;
    };
} // namespace PiScreen::datasource