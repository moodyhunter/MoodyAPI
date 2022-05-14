#include "SH1106.hpp"

#include <bitset>
#include <cstring>

#ifndef OLED_DEBUG
#define OLED_DEBUG 0
#endif

using namespace PiScreen::devices;

constexpr unsigned char InitSequence[]{
    0x00,
    // 0xae, // Display off
    0xa8, // Set multiplex ratio
    0x3f, // 64
    0xd3, // Set display offset
    0x00, // No offset
    0x40, // Set display start line
    0xa1, // Set segment re-map
    0xc8, // Set COM output scan direction
    0xda, // Set COM pins hardware configuration
    0x12, // Alternative COM pin configuration
    0x81, // Set contrast control register
    0xff, // Contrast value
    0xa4, // Disable entire display on
    0xa6, // Set normal display
    0xd5, // Set osc division
    0x80, 0x8d, 0x14,
    0xaf, // Display on
    0x20, // Set memory addressing mode
    0x02, // Sequential/Alternative mode
};

constexpr auto HDotsPerBlock = 8;
constexpr auto VDotsPerBlock = 8;

constexpr auto HBlocks = SH1106DeviceCore::SCREEN_WIDTH / HDotsPerBlock;
constexpr auto VBlocks = SH1106DeviceCore::SCREEN_HEIGHT / VDotsPerBlock;

SH1106DeviceCore::SH1106DeviceCore(int busId, int iAddr) : m_Addr(iAddr)
{
    // on Linux, SDA = bus number, SCL = device address
    m_I2CDevice = new common::I2CDevice(busId);

    // find the device address if requested
    if (!this->m_Addr)
    {
        if (m_I2CDevice->TestDevice(0x3c))
            m_Addr = 0x3c;
        else if (m_I2CDevice->TestDevice(0x3d))
            m_Addr = 0x3d;
    }
}

SH1106DeviceCore::~SH1106DeviceCore()
{
    delete m_I2CDevice;
}

bool SH1106DeviceCore::InitDevice(bool bFlip, bool bInvert)
{
    if (m_Addr == 0)
        return false;

    if (!m_I2CDevice->TestDevice(m_Addr))
        return false;

    p_I2CWrite(InitSequence, sizeof(InitSequence));

    if (bInvert)
    {
        unsigned char uc[4];
        uc[0] = 0;
        uc[1] = 0xa7;
        p_I2CWrite(uc, 2);
    }

    if (bFlip) // rotate display 180
    {
        unsigned char uc[4];
        uc[0] = 0;
        uc[1] = 0xa0;
        p_I2CWrite(uc, 2);
        uc[1] = 0xc0;
        p_I2CWrite(uc, 2);
    }

    return true;
}

void SH1106DeviceCore::p_I2CWrite(const unsigned char *pData, int iLen)
{
    m_I2CDevice->Write(m_Addr, pData, iLen);
}

void SH1106DeviceCore::SetPower(bool bOn)
{
    p_WriteCommand(bOn ? 0xaf : 0xae);
}

void SH1106DeviceCore::p_WriteCommand(unsigned char c)
{
    unsigned char buf[2];

    buf[0] = 0x00; // command introducer
    buf[1] = c;
    p_I2CWrite(buf, 2);
}

void SH1106DeviceCore::p_WriteCommand(unsigned char c, unsigned char d)
{
    unsigned char buf[3];

    buf[0] = 0x00;
    buf[1] = c;
    buf[2] = d;
    p_I2CWrite(buf, 3);
}

void SH1106DeviceCore::SetContrast(std::byte ucContrast)
{
    p_WriteCommand(0x81, (unsigned char) ucContrast);
}

void SH1106DeviceCore::p_SetPosition(int x, int y, bool bRender)
{

    if (!bRender)
        return; // don't send the commands to the OLED if we're not rendering the graphics now

    unsigned char buf[4];

    x += 2; // SH1106 has 128 pixels centered in 132

    buf[0] = 0x00;            // command introducer
    buf[1] = 0xb0 | y;        // page address (y)
    buf[2] = x & 0xf;         // lower column address (x)
    buf[3] = 0x10 | (x >> 4); // upper column address (x)
    p_I2CWrite(buf, 4);
}

void SH1106DeviceCore::p_WriteDataBlock(const unsigned char *ucBuf, int iLen, bool bRender)
{
    if (!bRender)
        return;

    unsigned char ucTemp[129];
    ucTemp[0] = 0x40; // data introducer

    std::memcpy(&ucTemp[1], ucBuf, iLen);
    p_I2CWrite(ucTemp, iLen + 1);

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

void SH1106DeviceCore::DrawBuffer(const uint8_t *const buf)
{
    for (int x = 0; x < HBlocks; x++)
    {
        for (int y = 0; y < VBlocks; y++)
        {
            unsigned char bytes_in[8]{ 0 };
            for (int byteN = 0; byteN < 8; byteN++)
            {
                const auto currentRow = buf[(y * HDotsPerBlock + byteN) * (HBlocks) + x];

                for (int bitN = 0; bitN < 8; bitN++)
                {
                    const auto bit = (currentRow >> (bitN)) & 1;
                    bytes_in[bitN] |= (bit << byteN);
                }
            }
            p_SetPosition(x * VBlocks, y, true);
            p_WriteDataBlock(bytes_in, 8, true);
        }
    }
}
