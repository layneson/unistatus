#include <stdlib.h>
#include <stdint.h>
#include "ws2811.h"

ws2811_t ledstring =
{
    .freq = WS2811_TARGET_FREQ,
    .dmanum = 5,
    .channel = 
    {
        [0] = 
        {
            .gpionum = 18,
            .count = 64,
            .invert = 0,
            .brightness = 100,
        }
    }
};

void setPixel(int x, int y, uint32_t rgb) {
    int pixel = getPixelPosition(x, y);
    ledstring.channel[0].leds[pixel] = rgb;
}

void setBrightness(int b) {
    ledstring.channel[0].brightness = b;
}

void init() {
    ws2811_init(&ledstring);
}

void show() {
    ws2811_render(&ledstring);
}

void stop() {
    ws2811_fini(&ledstring);
}

int getPixelPosition(int x, int y) {

    int map[8][8] = {
        {7 ,6 ,5 ,4 ,3 ,2 ,1 ,0 },
        {8 ,9 ,10,11,12,13,14,15},
        {23,22,21,20,19,18,17,16},
        {24,25,26,27,28,29,30,31},
        {39,38,37,36,35,34,33,32},
        {40,41,42,43,44,45,46,47},
        {55,54,53,52,51,50,49,48},
        {56,57,58,59,60,61,62,63}
    };

    return map[x][y];
}
