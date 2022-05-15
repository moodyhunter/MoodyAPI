#include "source.hpp"

#include <chrono>
#include <fstream>
#include <sstream>
#include <string>
#include <sys/sysinfo.h>
#include <thread>

using namespace std::chrono_literals;

namespace PiScreen::datasource
{
    CPUDataSource::CPUDataSource() : m_ncores(get_nprocs()), worker(&CPUDataSource::p_worker, this)
    {
    }

    CPUDataSource::~CPUDataSource()
    {
        stop = true;
        worker.join();
    }

    std::pair<unsigned long long, std::chrono::steady_clock::time_point> CPUDataSource::p_getIdle() const
    {
        const auto time = std::chrono::steady_clock::now();
        std::ifstream statStream{ "/proc/stat" };
        std::string line;
        std::getline(statStream, line);
        statStream.close();

        std::string fields[10];
        std::stringstream lineStream{ line };

        for (int i = 0; i < 10; i++)
            std::getline(lineStream, fields[i], ' ');

        return std::make_pair(std::stoull(fields[5]), time);
    }

    void CPUDataSource::p_worker()
    {
        auto [lastIdle, lastTime] = p_getIdle();

        while (!stop)
        {
            std::this_thread::sleep_for(std::chrono::seconds{ 1 });
            const auto [thisIdle, thisTime] = p_getIdle();
            const auto timeDiff = std::chrono::duration_cast<std::chrono::seconds>(thisTime - lastTime).count();

            m_idle = (thisIdle - lastIdle) / (float) timeDiff;
        }
    }

    std::string CPUDataSource::getData(const std::string &)
    {
        const int result = 100.0f - m_idle / m_ncores;
        return std::to_string(result) + "%";
    }
} // namespace PiScreen::datasource