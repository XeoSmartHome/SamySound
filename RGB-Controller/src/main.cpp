#include <Arduino.h>
#include <FastLED.h>

#define NUM_LEDS 60 

#define DATA_PIN D4

CRGB leds[NUM_LEDS];
CRGB now[NUM_LEDS];
CRGB history[NUM_LEDS];

void setup() { 
	Serial.begin(115200);
	Serial.println("resetting");
	FastLED.addLeds<WS2812,DATA_PIN,GRB>(leds,NUM_LEDS);
	FastLED.setBrightness(200);

	for (int i = 0; i < NUM_LEDS; i++) {
		now[i] = CRGB::Black;
		history[i] = CRGB::Black;
	}
}


void loop() { 
  	if (Serial.available()) {
		const int volume = Serial.parseInt();

		for (int i = 0; i < volume; i++) {
			now[i] = CRGB(min(i * 255 / 60 + 60, 255), max((NUM_LEDS - i) * 255 / 60 - 60, 0), 0);
		}

		for (int i = volume; i < NUM_LEDS; i++) {
			now[i] = CRGB::Black;
		}

		for (int i = 0; i < NUM_LEDS; i++) {
			leds[i] = CRGB((now[i].r + history[i].r * 3) / 4, (now[i].g + history[i].g * 3) / 4, 0);
		}

		for (int i = 0; i < NUM_LEDS; i++) {
			history[i] = leds[i];
		}

		FastLED.show();
	}

}