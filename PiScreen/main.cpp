#include "config/Config.hpp"
#include "device/device.hpp"
#include "renderer/Renderer.hpp"

int main(int argc, char *argv[])
{
    const auto device = new PiScreen::PiScreenDevice;
    device->SetContrast(std::byte{ 0xaa });

    PiScreen::config::ScreenContent config;
    {
        using namespace PiScreen::config;

#if 1
        config.push_back(MakeLine(0, 0, PiScreen::SCREEN_WIDTH, 0, 1));
        config.push_back(MakeLine(PiScreen::SCREEN_WIDTH, 0, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, 0, 0, 1));
#endif

        config.push_back(MakeStaticText(2, 12, "IP", 13, false, true));
        config.push_back(MakeStaticText(30, 12, "192.168.200.102", 11, false, true));
        config.push_back(MakeLine(0, 15, PiScreen::SCREEN_WIDTH, 14, 1));

        config.push_back(MakeStaticText(2, 27, "CPU", 13, false, true));
        config.push_back(MakeStaticText(33, 27, "50%", 10, false, false));

        config.push_back(MakeLine(60, 14, 60, 30, 1));

        config.push_back(MakeStaticText(64, 27, "MEM", 13, false, true));
        config.push_back(MakeStaticText(98, 27, "50%", 10, false, false));

        config.push_back(MakeLine(0, 31, PiScreen::SCREEN_WIDTH, 30, 1));

        config.push_back(MakeStaticText(2, 41, "Camera", 11, false, true));
        config.push_back(MakeStaticText(11, 59, "OFF", 18, false, true));

        config.push_back(MakeStaticText(98, 55, "👌🏻", 25, true, false));
    }

    PiScreen::renderer::ScreenRenderer renderer;
    renderer.InitDevice(device);
    renderer.SetConfiguration(config);
    renderer.RenderOne();

    return 0;
}
