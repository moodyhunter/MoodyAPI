#include "Renderer.hpp"

#include "device/device.hpp"

#include <string>

using namespace PiScreen::renderer;

ScreenRenderer::ScreenRenderer()
{
    m_CairoSurface = Cairo::ImageSurface::create(Cairo::Format::FORMAT_A1, SCREEN_WIDTH, SCREEN_HEIGHT);
    m_CairoContext = Cairo::Context::create(m_CairoSurface);
}

ScreenRenderer::~ScreenRenderer()
{
    m_CairoContext.clear();
    m_CairoSurface.clear();

    for (auto &ds : m_DataSources)
    {
        delete ds.second;
    }
}

bool ScreenRenderer::InitDevice(devices::IPiScreenDevice *pScreenDevice)
{
    m_pScreenDevice = pScreenDevice;
    return true;
}

void ScreenRenderer::SetConfiguration(const config::ScreenContent &config)
{
    m_Config = config;
}

std::string ScreenRenderer::GetDataSourceValue(int dataSourceId, const std::string &extInfo)
{
    if (m_DataSources.find(dataSourceId) == m_DataSources.end())
        m_DataSources[dataSourceId] = datasource::registrations[dataSourceId]();
    const auto &dataSource = m_DataSources[dataSourceId];
    return dataSource->getData(extInfo);
}

void ScreenRenderer::Render()
{
    // clear screen
    m_CairoContext->set_operator(Cairo::OPERATOR_CLEAR);
    m_CairoContext->rectangle(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT);
    m_CairoContext->paint_with_alpha(1.0);
    m_CairoContext->set_operator(Cairo::OPERATOR_SOURCE);

    for (auto &item : m_Config)
        RenderOne(item);
    m_pScreenDevice->DrawBuffer(m_CairoSurface->get_data());
}

void ScreenRenderer::RenderOne(const config::ScreenItem &item)
{
    using namespace PiScreen::config;
    const auto startX = GetScreenItemProp<int>(item, ITEM_PROP_START_X);
    const auto startY = GetScreenItemProp<int>(item, ITEM_PROP_START_Y);

    switch (GetScreenItemProp<ScreenItemType>(item, ITEM_TYPE))
    {
        case ITEM_TEXT_DATASOURCE:
        {
            const auto dataSourceId = GetScreenItemProp<int>(item, TEXT_PROP_DATASOURCE_ID);
            const auto dataSourceExtInfo = GetScreenItemProp<std::string>(item, config::TEXT_PROP_DATASOURCE_EXT_INFO);
            const auto value = GetDataSourceValue(dataSourceId, dataSourceExtInfo);
            auto newItem = item;
            newItem[config::TEXT_CONTENT] = { value };
            newItem[config::ITEM_TYPE] = { ITEM_TEXT };
            RenderOne(newItem);
            break;
        }
        case ITEM_TEXT:
        {
            const auto fontSize = GetScreenItemProp<int>(item, TEXT_PROP_FONT_SIZE);
            const auto text = GetScreenItemProp<std::string>(item, TEXT_CONTENT);
            const auto isEmoji = GetScreenItemProp<bool>(item, TEXT_PROP_IS_EMOJI);
            const auto isBold = GetScreenItemProp<bool>(item, TEXT_PROP_IS_BOLD);

            if (isEmoji)
            {
                m_CairoContext->select_font_face("emoji", Cairo::FONT_SLANT_NORMAL, Cairo::FONT_WEIGHT_NORMAL);
            }
            else
            {
                m_CairoContext->select_font_face("sans-serif", Cairo::FONT_SLANT_NORMAL, isBold ? Cairo::FONT_WEIGHT_BOLD : Cairo::FONT_WEIGHT_NORMAL);
            }

            m_CairoContext->set_font_size(fontSize);
            m_CairoContext->move_to(startX, startY);
            m_CairoContext->show_text(text);
            m_CairoContext->stroke();
            break;
        }
        case ITEM_LINE:
        {
            const auto endX = GetScreenItemProp<int>(item, LINE_PROP_END_X);
            const auto endY = GetScreenItemProp<int>(item, LINE_PROP_END_Y);
            const auto lineWidth = GetScreenItemProp<int>(item, LINE_PROP_WIDTH);

            m_CairoContext->set_line_width(lineWidth);
            m_CairoContext->move_to(startX, startY);
            m_CairoContext->line_to(endX, endY);
            m_CairoContext->stroke();
            break;
        }

        default: break;
    }
}
