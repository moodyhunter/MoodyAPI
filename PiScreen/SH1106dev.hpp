#pragma once

#include "I2CDevice.hpp"

#include <cstddef>
#include <cstdint>

class SH1106Device
{
  public:
    SH1106Device(int busId, int iAddr = 0);
    ~SH1106Device();

    bool initDevice(bool bFlip = false, bool bInvert = false);

    void setContrast(std::byte ucContrast);
    void setPower(bool bOn);
    void DrawBuffer(const uint8_t *const pBuffer);

  private:
    void p_WriteDataBlock(const unsigned char *ucBuf, int iLen, bool bRender);
    void p_SetPosition(int x, int y, bool bRender);
    void p_I2CWrite(const unsigned char *pData, int iLen);

    void p_WriteCommand(unsigned char c, unsigned char d);
    void p_WriteCommand(unsigned char c);

  private:
    I2CDevice *m_I2CDevice = nullptr;
    uint8_t m_Addr = 0;
    bool m_flip = false;
};
