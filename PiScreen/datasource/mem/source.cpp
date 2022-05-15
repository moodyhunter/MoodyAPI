#include "source.hpp"

#include <fstream>
#include <string>

namespace PiScreen::datasource
{
    MemoryDataSource::MemoryDataSource()
    {
    }

    MemoryDataSource::~MemoryDataSource()
    {
    }

    std::string MemoryDataSource::getData(const std::string &extInfo)
    {
        // get linux memory usage percentage
        std::ifstream meminfo("/proc/meminfo");
        std::string memTotalS;
        std::string memAvailableS;

        std::string line;
        while (std::getline(meminfo, line))
        {
            if (line.find("MemTotal") != std::string::npos)
            {
                memTotalS = line.substr(line.find_last_of(':') + 2);
            }
            else if (line.find("MemAvailable") != std::string::npos)
            {
                memAvailableS = line.substr(line.find_last_of(':') + 2);
            }
        }

        meminfo.close();

        // calculate memory usage percentage
        const auto memTotal = std::stoi(memTotalS);
        const auto memUsage = memTotal - std::stoi(memAvailableS);
        const auto memUsagePercent = memUsage * 100 / memTotal;

        // return memory usage percentage
        return std::to_string(memUsagePercent) + "%";
    }
} // namespace PiScreen::datasource
