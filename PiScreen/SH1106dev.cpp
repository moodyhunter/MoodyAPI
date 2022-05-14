#include "SH1106dev.hpp"

#include <bitset>
#include <cstring>

#ifndef OLED_DEBUG
#define OLED_DEBUG 0
#endif

#if OLED_DEBUG
#include <iostream>
#endif

constexpr unsigned char oled_initbuf[]{
    0x00, /*0xae,*/ 0xa8, 0x3f, 0xd3, 0x00, 0x40, 0xa1, 0xc8, 0xda, 0x12, 0x81, 0xff, 0xa4, 0xa6, 0xd5, 0x80, 0x8d, 0x14, 0xaf, 0x20, 0x02,
};
constexpr auto DOTS_PER_BLOCK_H = 8;
constexpr auto DOTS_PER_BLOCK_V = 8;

SH1106Device::SH1106Device(int busId, int iAddr) : m_busId(busId), m_Addr(iAddr)
{
    m_I2CDevice = new I2CDevice(m_busId); // on Linux, SDA = bus number, SCL = device address
    // find the device address if requested
    if (!this->m_Addr) // find it
    {
        m_I2CDevice->TestDevice(0x3c);
        if (m_I2CDevice->TestDevice(0x3c))
            m_Addr = 0x3c;
        else if (m_I2CDevice->TestDevice(0x3d))
            m_Addr = 0x3d;
        else
            return;
    }
    else
    {
        m_I2CDevice->TestDevice(m_Addr);
        if (!m_I2CDevice->TestDevice(m_Addr))
            return;
    }
}

bool SH1106Device::initDevice(bool bFlip, bool bInvert)
{
    if (m_Addr == 0)
        return false;

    unsigned char uc[4];

    // 132x64, 128x64 and 64x32
    p_I2CWrite((unsigned char *) oled_initbuf, sizeof(oled_initbuf));

    if (bInvert)
    {
        uc[0] = 0;
        uc[1] = 0xa7;
        p_I2CWrite(uc, 2);
    }

    if (bFlip) // rotate display 180
    {
        uc[0] = 0;
        uc[1] = 0xa0;
        p_I2CWrite(uc, 2);
        uc[1] = 0xc0;
        p_I2CWrite(uc, 2);
    }

    return true;
}

SH1106Device::~SH1106Device()
{
    delete m_I2CDevice;
}

void SH1106Device::p_I2CWrite(unsigned char *pData, int iLen)
{
    m_I2CDevice->Write(m_Addr, pData, iLen);
}

void SH1106Device::setPower(bool bOn)
{
    p_WriteCommand(bOn ? 0xaf : 0xae); // turn on OLED
}

void SH1106Device::p_WriteCommand(unsigned char c)
{
    unsigned char buf[2];

    buf[0] = 0x00; // command introducer
    buf[1] = c;
    p_I2CWrite(buf, 2);
}

void SH1106Device::p_WriteCommand(unsigned char c, unsigned char d)
{
    unsigned char buf[3];

    buf[0] = 0x00;
    buf[1] = c;
    buf[2] = d;
    p_I2CWrite(buf, 3);
}

void SH1106Device::setContrast(std::byte ucContrast)
{
    p_WriteCommand(0x81, (unsigned char) ucContrast);
}

void SH1106Device::p_SetPosition(int x, int y, bool bRender)
{
    unsigned char buf[4];

    if (!bRender)
        return; // don't send the commands to the OLED if we're not rendering the graphics now

    // SH1106 has 128 pixels centered in 132
    x += 2;

    buf[0] = 0x00;            // command introducer
    buf[1] = 0xb0 | y;        // set page to Y
    buf[2] = x & 0xf;         // lower column address
    buf[3] = 0x10 | (x >> 4); // upper column addr
    p_I2CWrite(buf, 4);
}

void SH1106Device::p_WriteDataBlock(const unsigned char *ucBuf, int iLen, bool bRender)
{
    unsigned char ucTemp[129];

    ucTemp[0] = 0x40; // data command
                      // Copying the data has the benefit in SPI mode of not letting
                      // the original data get overwritten by the SPI.transfer() function
    if (bRender)
    {
        std::memcpy(&ucTemp[1], ucBuf, iLen);
        p_I2CWrite(ucTemp, iLen + 1);
    }

#if OLED_DEBUG
    std::cout << "BEGIN DATA" << std::endl;

    int x = 0;
    bool buf[8][8]{ { 0 } };
    for (auto i = 0; i < iLen; i++)
    {
        for (auto b = 0; b < 8; b++)
            buf[b][x++ / 8] = (ucBuf[i] >> b) & 1;

        if (x % 64 == 0)
        {
            for (const auto aa : buf)
            {
                for (auto b = 0; b < 8; b++)
                    std::cout << (aa[b] ? '*' : ' ');
                std::cout << std::endl;
            }
            std::memset(buf, 0, 64);
        }
    }
#endif
}

void SH1106Device::setCursorPos(int x, int y)
{
    m_CursorX = x;
    m_CursorY = y;
}

void SH1106Device::DrawBuffer(const uint8_t *const buf)
{
    const auto HorizontalBlocksN = 128 / DOTS_PER_BLOCK_H;
    const auto VerticalBlocksN = 64 / DOTS_PER_BLOCK_V;

    // Y first because we want to fill the rows first
    for (int x = 0; x < HorizontalBlocksN; x++)
    {
        for (int y = 0; y < VerticalBlocksN; y++)
        {
            unsigned char bytes_in[8]{ 0 };
            for (int byteN = 0; byteN < 8; byteN++)
            {
                const auto currentRow = buf[(y * DOTS_PER_BLOCK_H + byteN) * (HorizontalBlocksN) + x];

                for (int bitN = 0; bitN < 8; bitN++)
                {
                    const auto bit = (currentRow >> (bitN)) & 1;
                    bytes_in[bitN] |= (bit << byteN);
                }
            }
            p_SetPosition(x * VerticalBlocksN, y, true);
            p_WriteDataBlock(bytes_in, 8, true);
        }
    }
}
