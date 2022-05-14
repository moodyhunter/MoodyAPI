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
        config.push_back(MakeLine(0, 0, PiScreen::SCREEN_WIDTH, 0, 1));
        config.push_back(MakeLine(PiScreen::SCREEN_WIDTH, 0, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, PiScreen::SCREEN_WIDTH, PiScreen::SCREEN_HEIGHT, 1));
        config.push_back(MakeLine(0, PiScreen::SCREEN_HEIGHT, 0, 0, 1));

        config.push_back(MakeText(20, 20, "Hello World!", 10, false, false));
        config.push_back(MakeText(80, 50, "👌🏻", 40, true, false));
        config.push_back(MakeLine(0, 20, 40, 40, 1));

        config.push_back(MakeDataSourceText(10, 50, "data:id1", 15, false, true));
    }

    PiScreen::renderer::ScreenRenderer renderer;
    renderer.InitDevice(device);
    renderer.SetConfiguration(config);
    renderer.RenderOne();

    return 0;
}
