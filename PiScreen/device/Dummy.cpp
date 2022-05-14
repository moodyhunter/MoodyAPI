#include "Dummy.hpp"

#include <bitset>

using namespace PiScreen::devices;

DummyDevice::DummyDevice()
{
    std::cerr << "Using Dummy device" << std::endl;
}

DummyDevice::~DummyDevice()
{
}

void DummyDevice::SetPower(bool bOn)
{
}

void DummyDevice::SetContrast(std::byte bContrast)
{
}

void DummyDevice::DrawBuffer(const uint8_t *buf)
{
    system("clear");

    for (int i = 0; i < 128 * 64 / 8; i++)
    {
        const auto c = buf[i];

        const auto str = std::bitset<8>(c).to_string();
        const auto revstr = std::string(str.rbegin(), str.rend());

        // Replace '0' with ' ' and '1' with space and '*'
        std::string str2;
        for (auto c : revstr)
            str2 += (c == '0' ? "　" : "＃");

        std::cout << str2;

        if (((i + 1) * 8) % 128 == 0)
            std::cout << std::endl;
    }
}
