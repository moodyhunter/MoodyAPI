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

#include "BitBang_I2C.h"

#include <fcntl.h>
#include <linux/i2c-dev.h>
#include <stdio.h>
#include <sys/ioctl.h>
#include <unistd.h>

void I2CInit(BBI2C *pI2C)
{
    if (!pI2C)
        return;

    char filename[32];
    sprintf(filename, "/dev/i2c-%d", pI2C->iBus);
    pI2C->file_i2c = open(filename, O_RDWR);
}

uint8_t I2CTest(BBI2C *pI2C, uint8_t addr)
{
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, addr) >= 0)
        return 1;
    return 0;
}

void I2CScan(BBI2C *pI2C, uint8_t *pMap)
{
    // clear the bitmap
    for (auto i = 0; i < 16; i++)
        pMap[i] = 0;

    // try every address
    for (auto i = 1; i < 128; i++)
        if (I2CTest(pI2C, i))
            pMap[i >> 3] |= (1 << (i & 7));
}

int I2CWrite(BBI2C *pI2C, uint8_t iAddr, uint8_t *pData, int iLen)
{
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
        if (write(pI2C->file_i2c, pData, iLen) >= 0)
            return 1;
    return 0;
}

int I2CReadRegister(BBI2C *pI2C, uint8_t iAddr, uint8_t u8Register, uint8_t *pData, int iLen)
{
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
    {
        write(pI2C->file_i2c, &u8Register, 1);
        return read(pI2C->file_i2c, pData, iLen) > 0;
    }
    return false;
}

int I2CRead(BBI2C *pI2C, uint8_t iAddr, uint8_t *pData, int iLen)
{
    int i = 0;
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
        i = read(pI2C->file_i2c, pData, iLen);
    return (i > 0);
}

DEVICE_TYPE I2CDiscoverDevice(BBI2C *pI2C, uint8_t i)
{
    // Probably an OLED display
    if (i == 0x3c || i == 0x3d)
    {
        uint8_t cTemp[8]{ 0 };
        I2CReadRegister(pI2C, i, 0x00, cTemp, 1);
        cTemp[0] &= 0xbf;    // mask off power on/off bit
        if (cTemp[0] == 0x8) // SH1106
            return DEVICE_SH1106;
        else if (cTemp[0] == 3 || cTemp[0] == 6)
            return DEVICE_SSD1306;
    }
    return DEVICE_UNKNOWN;
}
