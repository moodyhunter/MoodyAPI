#pragma once

#include "common/common.hpp"

#include <algorithm>
#include <bits/iterator_concepts.h>
#include <exception>
#include <list>
#include <map>
#include <stdexcept>
#include <string>
#include <variant>

namespace PiScreen::config
{
    enum ScreenItemType
    {
        ITEM_LINE,
        ITEM_TEXT,
        ITEM_TEXT_DATASOURCE,
    };

    enum ContentItemProperties
    {
        ITEM_TYPE, // Content type

        ITEM_PROP_START_X, // Start X
        ITEM_PROP_START_Y, // Start Y

        LINE_PROP_END_X, // End X
        LINE_PROP_END_Y, // End Y
        LINE_PROP_WIDTH, // Line Width

        TEXT_CONTENT,            // Text with static content
        TEXT_PROP_DATASOURCE_ID, // Text with content from a data source

        TEXT_PROP_FONT_SIZE, // Font size
        TEXT_PROP_IS_EMOJI,  // Is the text an emoji?
        TEXT_PROP_IS_BOLD,   // Is the text bold?
    };

    using _complex_content_entry_t = std::variant<int, bool, std::string, ScreenItemType>;

    struct ContentEntry : public _complex_content_entry_t
    {
        template<typename T>
        T getProperty(ContentItemProperties prop) const
        {
            if constexpr (std::is_same_v<T, int>)
            {
                switch (prop)
                {
                    case ITEM_PROP_START_X:
                    case ITEM_PROP_START_Y:
                    case LINE_PROP_END_X:
                    case LINE_PROP_END_Y:
                    case LINE_PROP_WIDTH:
                    case TEXT_PROP_FONT_SIZE:
                    {
                        return std::get<int>(*this);
                    }
                    default: throw std::runtime_error("Invalid entry type for integer type.");
                }
                common::Unreachable();
            }
            else if constexpr (std::is_same_v<T, std::string>)
            {
                switch (prop)
                {
                    case TEXT_CONTENT:
                    case TEXT_PROP_DATASOURCE_ID:
                    {
                        return std::get<std::string>(*this);
                    }
                    default: throw std::runtime_error("Invalid entry type for string type.");
                }
                common::Unreachable();
            }
            else if constexpr (std::is_same_v<T, bool>)
            {
                switch (prop)
                {
                    case TEXT_PROP_IS_EMOJI:
                    case TEXT_PROP_IS_BOLD:
                    {
                        return std::get<bool>(*this);
                    }
                    default: throw std::runtime_error("Invalid entry type for boolean type.");
                }
                common::Unreachable();
            }
            else if constexpr (std::is_same_v<T, ScreenItemType>)
            {
                switch (prop)
                {
                    case ITEM_TYPE:
                    {
                        return std::get<ScreenItemType>(*this);
                    }
                    default: throw std::runtime_error("Invalid entry type for ScreenContentType type.");
                }
                common::Unreachable();
            }
            else
            {
                throw std::runtime_error("Invalid entry type.");
            }
        }
    };

    using ScreenItem = std::map<ContentItemProperties, ContentEntry>;
    using ScreenContent = std::list<ScreenItem>;

    template<typename T>
    T GetScreenItemProp(const ScreenItem &item, ContentItemProperties prop)
    {
        return item.at(prop).getProperty<T>(prop);
    }

    inline ScreenItem MakeLine(int startX, int startY, int endX, int endY, int width)
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

    inline ScreenItem MakeText(int startX, int startY, std::string text, int fontSize, bool isEmoji, bool isBold)
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

    inline ScreenItem MakeDataSourceText(int startX, int startY, std::string dataSourceId, int fontSize, bool isEmoji, bool isBold)
    {
        ScreenItem item;
        item[ITEM_TYPE] = { ITEM_TEXT_DATASOURCE };
        item[ITEM_PROP_START_X] = { startX };
        item[ITEM_PROP_START_Y] = { startY };
        item[TEXT_PROP_DATASOURCE_ID] = { dataSourceId };
        item[TEXT_PROP_FONT_SIZE] = { fontSize };
        item[TEXT_PROP_IS_EMOJI] = { isEmoji };
        item[TEXT_PROP_IS_BOLD] = { isBold };
        return item;
    }

} // namespace PiScreen::config
