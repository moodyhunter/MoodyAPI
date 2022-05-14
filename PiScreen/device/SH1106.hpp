#pragma once

#include "IPiScreenDevice.hpp"
#include "common/I2CDevice.hpp"

#include <cstdlib>

namespace PiScreen::devices
{
    class SH1106DeviceCore
    {
      public:
        static constexpr auto SCREEN_WIDTH = 128;
        static constexpr auto SCREEN_HEIGHT = 64;

        SH1106DeviceCore(int busId, int iAddr = 0);
        ~SH1106DeviceCore();

        bool InitDevice(bool bFlip = false, bool bInvert = false);
        void SetPower(bool bOn);
        void SetContrast(std::byte ucContrast);
        void DrawBuffer(const uint8_t *pBuffer);

      private:
        void p_WriteDataBlock(const unsigned char *ucBuf, int iLen, bool bRender);
        void p_SetPosition(int x, int y, bool bRender);
        void p_I2CWrite(const unsigned char *pData, int iLen);

        void p_WriteCommand(unsigned char c, unsigned char d);
        void p_WriteCommand(unsigned char c);

      private:
        common::I2CDevice *m_I2CDevice = nullptr;
        uint8_t m_Addr = 0;
        bool m_flip = false;
    };

    class SH1106Device final : public IPiScreenDevice
    {
      public:
        static constexpr auto SCREEN_WIDTH = SH1106DeviceCore::SCREEN_WIDTH;
        static constexpr auto SCREEN_HEIGHT = SH1106DeviceCore::SCREEN_HEIGHT;

        SH1106Device()
        {
            int iChannel = -1;
            while (iChannel < 2)
            {
                iChannel++;
                m_pDevice = new SH1106DeviceCore{ iChannel };
                if (!m_pDevice->InitDevice())
                {
                    delete m_pDevice;
                    continue;
                }
                return;
            }

            std::cerr << "Unable to initialize I2C bus." << std::endl;
            std::cerr << "Please check your connections and verify the device address by typing 'i2cdetect -y <channel>" << std::endl;
            std::exit(EXIT_FAILURE);
        }

        // clang-format off
        ~SH1106Device() { delete m_pDevice; }

        void SetPower(bool bOn) override { m_pDevice->SetPower(bOn); }
        void SetContrast(std::byte ucContrast) override { m_pDevice->SetContrast(ucContrast); }
        void DrawBuffer(const uint8_t *pBuffer) override { m_pDevice->DrawBuffer(pBuffer); }
        // clang-format on

      private:
        SH1106DeviceCore *m_pDevice = nullptr;
    };

    static_assert(is_screen_device_valid_v<SH1106Device>, "Must be a complete class.");
} // namespace PiScreen::devices
