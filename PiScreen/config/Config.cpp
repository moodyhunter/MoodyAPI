#include "Config.hpp"

using namespace PiScreen::config;

namespace PiScreen::config
{
    ScreenItem MakeDataSourceText(int startX, int startY, int dataSourceId, std::string datasourceExtInfo, int fontSize, bool isEmoji, bool isBold)
    {
        ScreenItem item;
        item[ITEM_TYPE] = { ITEM_TEXT_DATASOURCE };
        item[ITEM_PROP_START_X] = { startX };
        item[ITEM_PROP_START_Y] = { startY };
        item[TEXT_PROP_DATASOURCE_ID] = { dataSourceId };
        item[TEXT_PROP_DATASOURCE_EXT_INFO] = { datasourceExtInfo };
        item[TEXT_PROP_FONT_SIZE] = { fontSize };
        item[TEXT_PROP_IS_EMOJI] = { isEmoji };
        item[TEXT_PROP_IS_BOLD] = { isBold };
        return item;
    }

    ScreenItem MakeStaticText(int startX, int startY, std::string text, int fontSize, bool isEmoji, bool isBold)
    {
        ScreenItem item;
        item[ITEM_TYPE] = { ITEM_TEXT };
        item[ITEM_PROP_START_X] = { startX };
        item[ITEM_PROP_START_Y] = { startY };
        item[TEXT_CONTENT] = { text };
        item[TEXT_PROP_FONT_SIZE] = { fontSize };
        item[TEXT_PROP_IS_EMOJI] = { isEmoji };
        item[TEXT_PROP_IS_BOLD] = { isBold };
        return item;
    }

    ScreenItem MakeLine(int startX, int startY, int endX, int endY, int width)
    {
        ScreenItem item;
        item[ITEM_TYPE] = { ITEM_LINE };
        item[ITEM_PROP_START_X] = { startX };
        item[ITEM_PROP_START_Y] = { startY };
        item[LINE_PROP_END_X] = { endX };
        item[LINE_PROP_END_Y] = { endY };
        item[LINE_PROP_WIDTH] = { width };
        return item;
    }
} // namespace PiScreen::config