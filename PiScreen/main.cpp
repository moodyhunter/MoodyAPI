#include "SSOLED.hpp"

#include <limits>
#include <string>

int main(int argc, char *argv[])
{
    if (argc != 4)
        return 1;

    unsigned char ucBackBuf[1024]{ 0 };

    // try I2C channel 0 through 2
    int iChannel = -1;
    while (iChannel < 2)
    {
        iChannel++;
        SSOLED ssoled{ iChannel, -1, true, false };
        if (ssoled.getDeviceType() == OLED_NOT_FOUND)
            continue;
        ssoled.setBackBuffer(ucBackBuf);
        ssoled.fill(0, false);
        ssoled.setTextWrap(true);

        const auto contrast = std::min(std::stoi(argv[1]), (int) std::numeric_limits<unsigned char>::max());
        ssoled.setContrast(contrast);

        const auto fontSize = (OLED_FONT_SIZE) std::stoi(argv[2]);
        ssoled.writeString(0, 0, 0, argv[3], fontSize, false, false);
        ssoled.drawBuffer();
        return 0;
    }

    printf("Unable to initialize I2C bus 0-2, please check your connections and verify the device address by typing 'i2cdetect -y <channel>\n");
    return 1;
}
