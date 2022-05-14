#include "Dummy.hpp"

DummyDevice::DummyDevice(int)
{
}
bool DummyDevice::InitDevice(bool bFlip, bool bInvert)
{
    return true;
}
void DummyDevice::SetPower(bool bOn)
{
}
void DummyDevice::SetContrast(std::byte bContrast)
{
}
void DummyDevice::DrawBuffer(const uint8_t *buf)
{
}
