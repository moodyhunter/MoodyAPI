#pragma once

#include "common/I2CDevice.hpp"
#include "common/IPiScreenDevice.hpp"

#include <type_traits>

class SH1106Device final : public IPiScreenDevice
{
  public:
    SH1106Device(int busId, int iAddr = 0);
    ~SH1106Device();

    bool InitDevice(bool bFlip = false, bool bInvert = false) override;
    void SetPower(bool bOn) override;
    void SetContrast(std::byte ucContrast) override;
    void DrawBuffer(const uint8_t *pBuffer) override;

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

static_assert(is_screen_device_constructable_v<SH1106Device>, "Must be a complete class.");