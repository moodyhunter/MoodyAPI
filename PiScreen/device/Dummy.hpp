#pragma once

#include "common/IPiScreenDevice.hpp"

// A dummy device for debugging purposes.
class DummyDevice final : public IPiScreenDevice
{
  public:
    DummyDevice(int);

    bool InitDevice(bool bFlip = false, bool bInvert = false);
    void SetPower(bool bOn);
    void SetContrast(std::byte bContrast);
    void DrawBuffer(const uint8_t *buf);
};

static_assert(is_screen_device_constructable_v<DummyDevice>, "Must be a complete class.");