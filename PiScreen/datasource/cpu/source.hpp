#pragma once

#include "datasource/IDataSource.hpp"

#include <chrono>
#include <thread>

namespace PiScreen::datasource
{
    class CPUDataSource : public IDataSource
    {
      public:
        CPUDataSource();
        ~CPUDataSource();

        virtual std::string getData(const std::string &extInfo) override;

      private:
        std::pair<unsigned long long, std::chrono::steady_clock::time_point> p_getIdle() const;
        void p_worker();

      private:
        bool stop = false;
        float m_idle;
        const int m_ncores;
        std::thread worker;
    };
} // namespace PiScreen::datasource
