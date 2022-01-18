//
// Bit Bang I2C library
// Copyright (c) 2018-2019 BitBank Software, Inc.
// Written by Larry Bank (bitbank@pobox.com)
// Project started 10/12/2018
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

#include "I2CDevice.hpp"

#include <fcntl.h>
#include <linux/i2c-dev.h>
#include <string>
#include <sys/ioctl.h>
#include <unistd.h>

I2CDevice::I2CDevice(int busId)
{
    iBus = busId;
    constexpr std::string_view v = "/dev/i2c-";
    const auto dname = std::string{ v } + std::to_string(busId);
    fd = open(dname.c_str(), O_RDWR);
}

I2CDevice::~I2CDevice()
{
    close(fd);
}

uint8_t I2CDevice::Test(uint8_t addr)
{
    if (ioctl(fd, I2C_SLAVE, addr) >= 0)
        return 1;
    return 0;
}

void I2CDevice::Scan(uint8_t *pMap)
{
    // clear the bitmap
    for (auto i = 0; i < 16; i++)
        pMap[i] = 0;

    // try every address
    for (auto i = 1; i < 128; i++)
        if (Test(i))
            pMap[i >> 3] |= (1 << (i & 7));
}

int I2CDevice::Write(uint8_t iAddr, uint8_t *pData, int iLen)
{
    if (ioctl(fd, I2C_SLAVE, iAddr) >= 0)
        if (write(fd, pData, iLen) >= 0)
            return 1;
    return 0;
}

int I2CDevice::ReadRegister(uint8_t iAddr, uint8_t u8Register, uint8_t *pData, int iLen)
{
    if (ioctl(fd, I2C_SLAVE, iAddr) >= 0)
    {
        write(fd, &u8Register, 1);
        return read(fd, pData, iLen) > 0;
    }
    return false;
}

int I2CDevice::Read(uint8_t iAddr, uint8_t *pData, int iLen)
{
    int i = 0;
    if (ioctl(fd, I2C_SLAVE, iAddr) >= 0)
        i = read(fd, pData, iLen);
    return (i > 0);
}

DEVICE_TYPE I2CDevice::DiscoverDevice(uint8_t i)
{
    // Probably an OLED display
    if (i == 0x3c || i == 0x3d)
    {
        uint8_t cTemp[8]{ 0 };
        ReadRegister(i, 0x00, cTemp, 1);
        cTemp[0] &= 0xbf;    // mask off power on/off bit
        if (cTemp[0] == 0x8) // SH1106
            return DEVICE_SH1106;
        else if (cTemp[0] == 3 || cTemp[0] == 6)
            return DEVICE_SSD1306;
    }
    return DEVICE_UNKNOWN;
}
