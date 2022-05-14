#include "Dummy.hpp"
#include "SH1106.hpp"

namespace PiScreen
{
#if PISCREEN_DUMMY_DEVICE
    using PiScreenDevice = devices::DummyDevice;
#else
    using PiScreenDevice = devices::SH1106Device;
#endif
    static_assert(devices::is_screen_device_valid_v<PiScreenDevice>, "Invalid device class selected.");
    // To be defined in main.cpp
    constexpr auto SCREEN_WIDTH = PiScreenDevice::SCREEN_WIDTH;
    constexpr auto SCREEN_HEIGHT = PiScreenDevice::SCREEN_HEIGHT;
} // namespace PiScreen
