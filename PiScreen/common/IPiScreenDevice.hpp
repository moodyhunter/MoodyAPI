#pragma once

#include <cstddef>
#include <cstdint>
#include <type_traits>

class IPiScreenDevice
{
  public:
    virtual bool InitDevice(bool bFlip = false, bool bInvert = false) = 0;
    virtual void SetPower(bool bOn) = 0;
    virtual void SetContrast(std::byte bContrast) = 0;
    virtual void DrawBuffer(const uint8_t *buf) = 0;
};

template<typename T>
constexpr auto is_screen_device_constructable_v = std::is_constructible_v<T, int>;
