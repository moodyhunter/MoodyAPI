#pragma once

#include "IPiScreenDevice.hpp"

namespace PiScreen::devices
{
    // A dummy device for debugging purposes.
    class DummyDevice final : public IPiScreenDevice
    {
      public:
        static constexpr auto SCREEN_WIDTH = 128;
        static constexpr auto SCREEN_HEIGHT = 64;

        DummyDevice();
        ~DummyDevice();

        void SetPower(bool bOn);
        void SetContrast(std::byte bContrast);
        void DrawBuffer(const uint8_t *buf);
    };

    static_assert(is_screen_device_valid_v<DummyDevice>, "Must be a complete class.");
} // namespace PiScreen::devices
