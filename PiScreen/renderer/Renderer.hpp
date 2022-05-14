#pragma once

#include "config/Config.hpp"
#include "device/IPiScreenDevice.hpp"

#include <cairomm/context.h>
#include <string>

namespace PiScreen::renderer
{
    class ScreenRenderer
    {
      public:
        explicit ScreenRenderer();
        ~ScreenRenderer();

        bool InitDevice(devices::IPiScreenDevice *pScreenDevice);
        void SetConfiguration(const config::ScreenContent &config);
        void RenderOne();

      private:
        std::string GetDataSourceValue(const std::string &dataSourceId);

      private:
        Cairo::RefPtr<Cairo::ImageSurface> m_CairoSurface;
        Cairo::RefPtr<Cairo::Context> m_CairoContext;
        config::ScreenContent m_Config;

        PiScreen::devices::IPiScreenDevice *m_pScreenDevice = nullptr;
    };
} // namespace PiScreen::renderer
