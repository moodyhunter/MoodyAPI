#include "SH1106dev.hpp"

#include <bitset>
#include <cairo-ft.h>
#include <cairo.h>
#include <cmath>
#include <iostream>

#ifndef OLED_DEBUG
#define OLED_DEBUG 0
#endif

void oled_print_buffer(const unsigned char *const &buf)
{
    for (int i = 0; i < 128 * 64 / 8; i++)
    {
        const auto c = buf[i];

        const auto str = std::bitset<8>(c).to_string();
        const auto revstr = std::string(str.rbegin(), str.rend());

        // Replace '0' with ' ' and '1' with space and '*'
        std::string str2;
        for (auto c : revstr)
            str2 += (c == '0' ? ' ' : '*');

        std::cout << str2;

        if (((i + 1) * 8) % 128 == 0)
            std::cout << std::endl;
    }
}

int main(int argc, char *argv[])
{
    int iChannel = -1;
    while (iChannel < 2)
    {
        iChannel++;

        SH1106Device ssoled{ iChannel };
        if (!ssoled.initDevice())
            continue;

        ssoled.setContrast(std::byte{ 0xaa });

        const auto surface = cairo_image_surface_create(CAIRO_FORMAT_A1, 128, 64);
        const auto cr = cairo_create(surface);

        cairo_set_line_width(cr, 1.0);
        cairo_move_to(cr, 0, 5);
        cairo_line_to(cr, 0, 0);
        cairo_line_to(cr, 128, 0);
        cairo_line_to(cr, 128, 5);
        cairo_line_to(cr, 60, 32);

        FcPattern *pattern;
        {
            FcResult result;
            FcInit();
            pattern = FcNameParse((const FcChar8 *) "JoyPixels");
            FcDefaultSubstitute(pattern);
            FcConfigSubstitute(FcConfigGetCurrent(), pattern, FcMatchPattern);
            pattern = FcFontMatch(FcConfigGetCurrent(), pattern, &result);
        }

        const auto face = cairo_ft_font_face_create_for_pattern(pattern);
        cairo_set_font_face(cr, face);
        cairo_set_font_size(cr, 20.0);
        cairo_move_to(cr, 0, 50.0);
        cairo_show_text(cr, "🍆🍑");
        cairo_move_to(cr, 5, 20.0);
        cairo_set_font_size(cr, 20.0);
        cairo_show_text(cr, "🍆🍑");
        {
            double xc = 80.0;
            double yc = 30.0;
            double radius = 20.0;
            double angle1 = 45.0 * (M_PI / 180.0);  /* angles are specified */
            double angle2 = 180.0 * (M_PI / 180.0); /* in radians           */

            cairo_set_line_width(cr, 3.0);
            cairo_arc(cr, xc, yc, radius, angle1, angle2);
            cairo_stroke(cr);

            /* draw helping lines */
            cairo_set_source_rgba(cr, 1, 0.2, 0.2, 0.6);
            cairo_set_line_width(cr, 6.0);

            cairo_arc(cr, xc, yc, 10.0, 0, 2 * M_PI);
            cairo_fill(cr);

            cairo_arc(cr, xc, yc, radius, angle1, angle1);
            cairo_line_to(cr, xc, yc);
            cairo_arc(cr, xc, yc, radius, angle2, angle2);
            cairo_line_to(cr, xc, yc);
            cairo_stroke(cr);
        }
        cairo_stroke(cr);

        unsigned char *const buf = cairo_image_surface_get_data(surface);

#if OLED_DEBUG
        oled_print_buffer(buf);
#endif

        ssoled.DrawBuffer(buf);
        return 0;
    }

    std::cout << "Unable to initialize I2C bus." << std::endl;
    std::cout << "Please check your connections and verify the device address by typing 'i2cdetect -y <channel>" << std::endl;
    return 1;
}
