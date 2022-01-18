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
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <unistd.h>

void I2CInit(BBI2C *pI2C)
{
    if (pI2C == NULL)
        return;

    char filename[32];
    sprintf(filename, "/dev/i2c-%d", pI2C->iBus);
    if ((pI2C->file_i2c = open(filename, O_RDWR)) < 0)
        return;

    return;
}

//
// Test a specific I2C address to see if a device responds
// returns 0 for no response, 1 for a response
//
uint8_t I2CTest(BBI2C *pI2C, uint8_t addr)
{
    uint8_t response = 0;

    if (ioctl(pI2C->file_i2c, I2C_SLAVE, addr) >= 0)
        response = 1;
    return response;
}

//
// Scans for I2C devices on the bus
// returns a bitmap of devices which are present (128 bits = 16 bytes, LSB first)
// A set bit indicates that a device responded at that address
//
void I2CScan(BBI2C *pI2C, uint8_t *pMap)
{
    int i;
    for (i = 0; i < 16; i++) // clear the bitmap
        pMap[i] = 0;
    for (i = 1; i < 128; i++) // try every address
    {
        if (I2CTest(pI2C, i))
        {
            pMap[i >> 3] |= (1 << (i & 7));
        }
    }
}

//
// Write I2C data
// quits if a NACK is received and returns 0
// otherwise returns the number of bytes written
//
int I2CWrite(BBI2C *pI2C, uint8_t iAddr, uint8_t *pData, int iLen)
{
    int rc = 0;

    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
    {
        if (write(pI2C->file_i2c, pData, iLen) >= 0)
            rc = 1;
    }
    return rc;
}

//
// Read N bytes starting at a specific I2C internal register
//
int I2CReadRegister(BBI2C *pI2C, uint8_t iAddr, uint8_t u8Register, uint8_t *pData, int iLen)
{
    int rc;

    int i = 0;
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
    {
        write(pI2C->file_i2c, &u8Register, 1);
        i = read(pI2C->file_i2c, pData, iLen);
    }
    return (i > 0);
}

//
// Read N bytes
//
int I2CRead(BBI2C *pI2C, uint8_t iAddr, uint8_t *pData, int iLen)
{
    int rc;

    int i = 0;
    if (ioctl(pI2C->file_i2c, I2C_SLAVE, iAddr) >= 0)
    {
        i = read(pI2C->file_i2c, pData, iLen);
    }
    return (i > 0);
}

//
// Figure out what device is at that address
// returns the enumerated value
//
int I2CDiscoverDevice(BBI2C *pI2C, uint8_t i)
{
    uint8_t j, cTemp[8];
    int iDevice = DEVICE_UNKNOWN;

    if (i == 0x3c || i == 0x3d) // Probably an OLED display
    {
        I2CReadRegister(pI2C, i, 0x00, cTemp, 1);
        cTemp[0] &= 0xbf;    // mask off power on/off bit
        if (cTemp[0] == 0x8) // SH1106
            iDevice = DEVICE_SH1106;
        else if (cTemp[0] == 3 || cTemp[0] == 6)
            iDevice = DEVICE_SSD1306;
        return iDevice;
    }

    return iDevice;
}