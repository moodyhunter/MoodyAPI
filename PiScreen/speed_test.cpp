//
// oled test program
// Written by Larry Bank
#include "ss_oled.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

int millis()
{
    int iTime;
    struct timespec res;

    clock_gettime(CLOCK_MONOTONIC, &res);
    iTime = 1000 * res.tv_sec + res.tv_nsec / 1000000;

    return iTime;
} /* millis() */

#define OLED_WIDTH 128
#define OLED_HEIGHT 64

void SpeedTest(SSOLED *ssoled)
{
    int i, x, y;
    char szTemp[32];
    unsigned long ms;

    ssoled->fill(0x0, 1);
    ssoled->writeString(0, 16, 0, (char *) "ss_oled Demo", FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "Written by Larry Bank", FONT_SMALL, 1, 1);
    ssoled->writeString(0, 0, 3, (char *) "**Demo**", FONT_LARGE, 0, 1);
    usleep(2000000);

    // Pixel and line functions won't work without a back buffer
    ssoled->fill(0, 1);
    ssoled->writeString(0, 0, 0, (char *) "Backbuffer Test", FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "3000 Random dots", FONT_NORMAL, 0, 1);
    usleep(2000000);
    ssoled->fill(0, 1);
    ms = millis();
    for (i = 0; i < 3000; i++)
    {
        x = random() & (OLED_WIDTH - 1);
        y = random() & (OLED_HEIGHT - 1);
        ssoled->setPixel(x, y, 1, 1);
    }
    ms = millis() - ms;
    sprintf(szTemp, "%dms", (int) ms);
    ssoled->writeString(0, 0, 0, szTemp, FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "Without backbuffer", FONT_SMALL, 0, 1);
    usleep(2000000);
    ssoled->fill(0, 1);
    ms = millis();
    for (i = 0; i < 3000; i++)
    {
        x = random() & (OLED_WIDTH - 1);
        y = random() & (OLED_HEIGHT - 1);
        ssoled->setPixel(x, y, 1, 0);
    }
    ssoled->dumpBuffer(NULL);
    ms = millis() - ms;
    sprintf(szTemp, "%dms", (int) ms);
    ssoled->writeString(0, 0, 0, szTemp, FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "With backbuffer", FONT_SMALL, 0, 1);
    usleep(2000000);
    ssoled->fill(0, 1);
    ssoled->writeString(0, 0, 0, (char *) "Backbuffer Test", FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "96 lines", FONT_NORMAL, 0, 1);
    usleep(2000000);
    ms = millis();
    for (x = 0; x < OLED_WIDTH - 1; x += 2)
    {
        ssoled->drawLine(x, 0, OLED_WIDTH - x, OLED_HEIGHT - 1, 1);
    }
    for (y = 0; y < OLED_HEIGHT - 1; y += 2)
    {
        ssoled->drawLine(OLED_WIDTH - 1, y, 0, OLED_HEIGHT - 1 - y, 1);
    }
    ms = millis() - ms;
    sprintf(szTemp, "%dms", (int) ms);
    ssoled->writeString(0, 0, 0, szTemp, FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "Without backbuffer", FONT_SMALL, 0, 1);
    usleep(2000000);
    ssoled->fill(0, 1);
    ms = millis();
    for (x = 0; x < OLED_WIDTH - 1; x += 2)
    {
        ssoled->drawLine(x, 0, OLED_WIDTH - 1 - x, OLED_HEIGHT - 1, 0);
    }
    for (y = 0; y < OLED_HEIGHT - 1; y += 2)
    {
        ssoled->drawLine(OLED_WIDTH - 1, y, 0, OLED_HEIGHT - 1 - y, 0);
    }
    ssoled->dumpBuffer(NULL);
    ms = millis() - ms;
    sprintf(szTemp, "%dms", (int) ms);
    ssoled->writeString(0, 0, 0, szTemp, FONT_NORMAL, 0, 1);
    ssoled->writeString(0, 0, 1, (char *) "With backbuffer", FONT_SMALL, 0, 1);
    usleep(2000000);
    ssoled->fill(0, 1);
    ssoled->writeString(0, 0, 0, (char *) "Fill Test", FONT_NORMAL, 0, 1);
    ms = millis();
    for (x = 0; x < 50; x++)
    {
        ssoled->fill(0, 1);
        ssoled->fill(0xff, 1);
    }
    ms = millis() - ms;
    ssoled->writeString(0, 0, 0, (char *) "Fill rate", FONT_NORMAL, 0, 1);
    sprintf(szTemp, "%d FPS", (int) (100000 / ms));
    ssoled->writeString(0, 0, 1, szTemp, FONT_NORMAL, 0, 1);
}

int main(int argc, char *argv[])
{
    unsigned char ucBackBuf[1024]{ 0 };

    // try I2C channel 0 through 2
    int iChannel = -1;
    while (iChannel < 2)
    {
        iChannel++;
        const auto ssoled = new SSOLED(-1, true, false);
        if (ssoled->DeviceType() == OLED_NOT_FOUND)
            continue;
        ssoled->setBackBuffer(ucBackBuf);
        SpeedTest(ssoled);
        printf("Press ENTER to quit\n");
        getchar();
        ssoled->setPower(false);
        return 0;
    }

    printf("Unable to initialize I2C bus 0-2, please check your connections and verify the device address by typing 'i2cdetect -y <channel>\n");
    return 1;
} /* main() */
