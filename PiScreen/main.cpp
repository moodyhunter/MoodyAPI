#include "device/Dummy.hpp"
#include "device/SH1106.hpp"

#include <bitset>
#include <cairomm/cairomm.h>
#include <cmath>
#include <iostream>

#if PISCREEN_DUMMY_DEVICE
using PiScreenDevice = DummyDevice;
#else
using PiScreenDevice = SH1106Device;
#endif

int main(int argc, char *argv[])
{
    int iChannel = -1;
    while (iChannel < 2)
    {
        iChannel++;

        PiScreenDevice device{ iChannel };
        if (!device.InitDevice())
            continue;

        device.SetContrast(std::byte{ 0xaa });

        auto surface = Cairo::ImageSurface::create(Cairo::Format::FORMAT_A1, 128, 64);
        const auto cr = Cairo::Context::create(surface);

        cr->set_line_width(1);

        cr->move_to(0, 0);

        cr->line_to(0, 64);
        cr->line_to(128, 64);
        cr->line_to(128, 0);
        cr->line_to(0, 0);
        cr->stroke();

        FcPattern *pattern;
        {
            FcInit();
            FcConfig *config = FcInitLoadConfigAndFonts();
            FcPattern *pat = FcNameParse((const FcChar8 *) "WenQuanYi Micro Hei");

            FcConfigSubstitute(config, pat, FcMatchFont);
            FcDefaultSubstitute(pat);

            FcResult result;
            pattern = FcFontMatch(config, pat, &result);
        }

        const auto face = Cairo::FtFontFace::create(pattern);
        cr->set_font_face(face);
        cr->set_font_size(20.0);

        cr->move_to(0, 20.0);
        cr->show_text("哈哈AaBb");

        cr->move_to(0, 50.0);
        cr->show_text("🍆🍑？");

        cr->move_to(0, 80.0);
        {
            double xc = 80.0;
            double yc = 30.0;
            double radius = 20.0;
            double angle1 = 45.0 * (M_PI / 180.0);  /* angles are specified */
            double angle2 = 180.0 * (M_PI / 180.0); /* in radians           */

            cr->set_line_width(3.0);
            cr->arc(xc, yc, radius, angle1, angle2);
            cr->stroke();

            /* draw helping lines */
            cr->set_source_rgba(1, 0.2, 0.2, 0.6);
            cr->set_line_width(6.0);

            cr->arc(xc, yc, 10.0, 0, 2 * M_PI);
            cr->fill();

            cr->arc(xc, yc, radius, angle1, angle1);
            cr->line_to(xc, yc);
            cr->arc(xc, yc, radius, angle2, angle2);
            cr->line_to(xc, yc);
            cr->stroke();
        }
        cr->stroke();

        unsigned char *const buf = surface->get_data();

        device.DrawBuffer(buf);
        return 0;
    }

    std::cout << "Unable to initialize I2C bus." << std::endl;
    std::cout << "Please check your connections and verify the device address by typing 'i2cdetect -y <channel>" << std::endl;
    return 1;
}
