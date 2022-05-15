#include "source.hpp"

#include <arpa/inet.h>
#include <ifaddrs.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>

namespace PiScreen::datasource
{
    IPAddressDataSource::IPAddressDataSource()
    {
    }

    IPAddressDataSource::~IPAddressDataSource()
    {
    }

    std::string IPAddressDataSource::getData(const std::string &extInfo)
    {
        struct ifaddrs *addrs = NULL;

        getifaddrs(&addrs);

        for (auto iface = addrs; iface != NULL; iface = iface->ifa_next)
        {
            if (iface->ifa_addr && iface->ifa_addr->sa_family == AF_INET)
            {
                const auto tmpAddrPtr = &((struct sockaddr_in *) iface->ifa_addr)->sin_addr;
                char addressBuffer[INET_ADDRSTRLEN];
                inet_ntop(AF_INET, tmpAddrPtr, addressBuffer, INET_ADDRSTRLEN);

                std::string ipAddress = addressBuffer;
                std::string ifName = iface->ifa_name;

                if (!extInfo.empty())
                {
                    if (ifName == extInfo)
                    {
                        freeifaddrs(addrs);
                        return ipAddress;
                    }
                }
                else if (ifName.starts_with("enp") || ifName.starts_with("wlan") || ifName.starts_with("eth") || ifName.starts_with("wlp"))
                {
                    freeifaddrs(addrs);
                    return ipAddress;
                }
            }
        }

        if (addrs)
            freeifaddrs(addrs);

        return "N/A";
    }

} // namespace PiScreen::datasource