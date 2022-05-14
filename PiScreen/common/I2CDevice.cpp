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

using namespace PiScreen::common;

constexpr std::string_view dev_i2c = "/dev/i2c-";

I2CDevice::I2CDevice(int busId) : m_busId(busId)
{
    const auto deviceName = std::string(dev_i2c) + std::to_string(busId);
    m_fd = open(deviceName.c_str(), O_RDWR);
}

I2CDevice::~I2CDevice()
{
    close(m_fd);
}

bool I2CDevice::TestDevice(uint8_t iAddr)
{
    return ioctl(m_fd, I2C_SLAVE, iAddr) >= 0;
}

void I2CDevice::Scan(uint8_t *pMap)
{
    // clear the bitmap
    for (auto i = 0; i < 16; i++)
        pMap[i] = 0;

    // try every address
    for (auto i = 1; i < 128; i++)
        if (TestDevice(i))
            pMap[i >> 3] |= (1 << (i & 7));
}

bool I2CDevice::Write(uint8_t iAddr, const uint8_t *pData, int iLen)
{
    if (!TestDevice(iAddr))
        return false;

    return write(m_fd, pData, iLen) >= 0;
}

bool I2CDevice::ReadRegister(uint8_t iAddr, uint8_t u8Register, uint8_t *pData, int iLen)
{
    if (!TestDevice(iAddr))
        return false;

    write(m_fd, &u8Register, 1);
    return read(m_fd, pData, iLen) > 0;
}

bool I2CDevice::Read(uint8_t iAddr, uint8_t *pData, int iLen)
{
    if (!TestDevice(iAddr))
        return false;

    return read(m_fd, pData, iLen) > 0;
}
