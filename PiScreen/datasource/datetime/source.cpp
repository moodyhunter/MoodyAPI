#include "source.hpp"

namespace PiScreen::datasource
{
    DateTimeDateSource::DateTimeDateSource()
    {
    }

    DateTimeDateSource::~DateTimeDateSource()
    {
    }

    std::string DateTimeDateSource::getData(const std::string &extInfo)
    {

        const auto now = time(nullptr);
        const auto tstruct = *localtime(&now);
        char buf[80];

        auto format = extInfo;
        if (format.empty())
            format = "%Y-%m-%d %H:%M:%S";

        strftime(buf, sizeof(buf), format.c_str(), &tstruct);
        return buf;
    }
} // namespace PiScreen::datasource