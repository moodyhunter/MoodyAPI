#include "config/Config.hpp"
#include "datasource/datasource.hpp"
#include "datasource/mem/source.hpp"
#include "device/device.hpp"
#include "renderer/Renderer.hpp"

#ifndef PISCREEN_OUTPUT_LIMIT
#warning "PISCREEN_OUTPUT_LIMIT not defined, defaulting to 0"
#define PISCREEN_OUTPUT_LIMIT 0
#endif

using namespace std::chrono_literals;

int main(int argc, char *argv[])
{
    PiScreen::PiScreenDevice device;
    device.SetContrast(std::byte{ 0xaa });

    PiScreen::config::ScreenContent config;
    {
        using namespace PiScreen::config;

#if 1
        config.push_back(MakeLine(0, 0, PiScreen::SCREEN_WIDTH, 0, 1));
        config.push_back(MakeLine(PiScreen::SCREEN_WIDTH, 0, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, 0, 0, 1));
#endif

        {
            config.push_back(MakeDataSourceText(2, 12, PiScreen::datasource::DateTimeDateSource_ID, "%H:%M:%S", 11, false, true));
            config.push_back(MakeLine(50, 0, 50, 15, 1));
            config.push_back(MakeDataSourceText(58, 12, PiScreen::datasource::CPUDataSource_ID, "", 10, false, true));
            config.push_back(MakeLine(93, 0, 86, 15, 1));
            config.push_back(MakeDataSourceText(97, 12, PiScreen::datasource::MemoryDataSource_ID, "", 10, false, true));
        }

        config.push_back(MakeLine(0, 15, PiScreen::SCREEN_WIDTH, 14, 1));

        {
            config.push_back(MakeStaticText(2, 27, "IP", 13, false, true));
            config.push_back(MakeLine(18, 14, 18, 30, 1));
            config.push_back(MakeDataSourceText(22, 26, PiScreen::datasource::IPAddressDataSource_ID, "", 11, false, true));

            config.push_back(MakeStaticText(93, 47, "👌", 30, true, false));
        }

        config.push_back(MakeLine(0, 31, PiScreen::SCREEN_WIDTH - 27, 30, 1));

        {
            config.push_back(MakeStaticText(2, 41, "Camera", 11, false, true));
            config.push_back(MakeDataSourceText(10, 59, PiScreen::datasource::SystemdServiceDataSource_ID, "motion.service", 18, false, true));

            config.push_back(MakeLine(48, 31, 47, 64, 1));

            config.push_back(MakeStaticText(50, 41, "Notifier", 11, false, true));

            config.push_back(MakeLine(92, 31, 92, 64, 1));
            config.push_back(MakeStaticText(59, 59, "???", 18, false, true));
        }
    }

    PiScreen::renderer::ScreenRenderer renderer;
    renderer.InitDevice(&device);
    renderer.SetConfiguration(config);

#if PISCREEN_OUTPUT_LIMIT > 0
    std::size_t frameCount = 0;
    while (frameCount < PISCREEN_OUTPUT_LIMIT)
    {
        renderer.Render();
        std::this_thread::sleep_for(1s);
        ++frameCount;
    }
#else
    while (true)
    {
        renderer.Render();
        std::this_thread::sleep_for(1s);
    }
#endif

    return 0;
}
