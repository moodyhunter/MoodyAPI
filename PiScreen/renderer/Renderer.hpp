#pragma once

#include "config/Config.hpp"
#include "datasource/IDataSource.hpp"
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

        void Render();
        void RenderOne(const config::ScreenItem &item);

      private:
        std::string GetDataSourceValue(int dataSourceId, const std::string &extInfo);
        void RenderText(int startX, int startY, const config::ScreenItem &item);

      private:
        Cairo::RefPtr<Cairo::ImageSurface> m_CairoSurface;
        Cairo::RefPtr<Cairo::Context> m_CairoContext;
        std::map<int, datasource::IDataSource *> m_DataSources;

        config::ScreenContent m_Config;
        PiScreen::devices::IPiScreenDevice *m_pScreenDevice = nullptr;
    };
} // namespace PiScreen::renderer
