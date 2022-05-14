#pragma once

#include "common/common.hpp"

#include <cstddef>
#include <cstdint>
#include <type_traits>

namespace PiScreen::devices
{
    class IPiScreenDevice
    {
      public:
        virtual void SetPower(bool bOn) = 0;
        virtual void SetContrast(std::byte bContrast) = 0;
        virtual void DrawBuffer(const uint8_t *buf) = 0;
    };

    template<typename T>
    constexpr auto is_screen_device_valid_v = std::is_constructible_v<T>;
} // namespace PiScreen::devices
