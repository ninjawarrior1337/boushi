#include <Arduino.h>

#define FASTLED_ESP8266_D1_PIN_ORDER
#include <FastLED.h>
#include <ESPAsyncWebServer.h>
#include <ESP8266WiFi.h>
#include <LittleFS.h>
#include <ArduinoJson.h>
#include <AsyncJson.h>
#include <DNSServer.h>
#include <ESP8266mDNS.h>

#include "./art.h"
#include "./creds.h"

#define DATA_PIN D2

// Define the array of leds
CRGB leds[STRIP_LEN];

//Define DNSServer
DNSServer dnsServer;

//Define webserver
AsyncWebServer server(80);

//Define the lock
bool lock = false;
String password = LOCK_PW;

//Declare mode flag
String mode = "cycle";

//Declare cycle manual override
bool interractedWith = false;

bool isLocked(AsyncWebServerRequest *req) {
  if(lock) {
    req->send(403, "text/plain", "LOCKED");
    return lock;
  } else {
    return false;
  }
}

void cancelCycle() {
  mode = "static";
  interractedWith = true;
}

void setup() {
  Serial.begin(9600);
  //Setup wifi
  #ifdef DEV
  WiFi.begin(WIFI_SSID, WIFI_PW);
  while(WiFi.status() != WL_CONNECTED) {
    delay(300);
    Serial.print(".");
  }
  Serial.println("Connected");
  #endif

  #ifdef PROD
  WiFi.mode(WIFI_AP);
  WiFi.softAPConfig(IPAddress(21, 21, 21, 21), IPAddress(21, 21, 21, 21), IPAddress(255, 255, 255, 0));
  WiFi.softAP("Boushi");

  dnsServer.setErrorReplyCode(DNSReplyCode::ServerFailure);
  Serial.print(dnsServer.start(53, "*", WiFi.softAPIP()));
  #endif

  if (!MDNS.begin("hat")) {
    Serial.println("Error setting up MDNS responder!");
  }

  //LittleFS
  LittleFS.begin();

  //Setup Server
  DefaultHeaders::Instance().addHeader("Access-Control-Allow-Origin", "*");
  server.serveStatic("/", LittleFS, "/").setDefaultFile("index.html");

  //Handle locking and unlocking
  server.on("/lock", HTTP_POST, [](AsyncWebServerRequest *req) {
    if(req->params() > 0) {
      if(req->getParam("pass", true)->value() == password) {
        lock = !lock;
        if(lock) {
          req->send(200, "text/plain", "LOCKED");
        } else {
          req->send(200, "text/plain", "UNLOCKED");
        }
      } else {
        req->send(403, "text/plain", "INCORRECT");
      }
    }
  });
  
  //Read only get all available art on device
  server.on("/api/art", HTTP_GET, [](AsyncWebServerRequest *req) {
    DynamicJsonDocument doc(1024);
    for(uint i = 0; i < ART_COUNT; i++) {
      doc["data"][i] = ART_LIST[i];
    }
    req->send(200, "application/json", doc.as<String>());
  });
  server.on("/api/gifs", HTTP_GET, [](AsyncWebServerRequest *req) {
    DynamicJsonDocument doc(1024);
    for(uint i = 0; i < GIF_COUNT; i++) {
      doc["data"][i] = GIF_LIST[i];
    }
    req->send(200, "application/json", doc.as<String>());
  });

  //Write all pixels to color specified by req
  server.on("/api/fill", [](AsyncWebServerRequest *req) {
    cancelCycle();
    if(isLocked(req)) {
      return;
    }
    if(req->params() > 0) {
      for(auto i = 0; i < STRIP_LEN; i++)
        leds[i] = CRGB(gamma32(atoi(req->getParam(0)->value().c_str())));
    }
    FastLED.show();
    req->send(200, "text/plain", "FILLED");
  });

  //Write pixels to pattern specified by req
  server.on("/api/draw", HTTP_POST, [](AsyncWebServerRequest *req) {
    cancelCycle();
    if(isLocked(req)) {
      return;
    }
    if(req->hasParam("pixels", true, false)) {
      DynamicJsonDocument doc(4096);
      String s = req->getParam("pixels", true)->value();
      deserializeJson(doc, s);
      for(auto i = 0; i < STRIP_LEN; i++) {
        uint32_t color = doc["data"][i];
        leds[i] = CRGB(gamma32(color));
      }
      FastLED.show();
      req->send(200, "text/plain", "DRAWN");
    }
  });

  //Write pixels to pattern stored on device
  server.on("/api", [](AsyncWebServerRequest *req) {
    cancelCycle();
    if(isLocked(req)) {
      return;
    }
    if(req->hasParam("art")) {
      renderGif = false;
      switchGrid(leds, req->getParam("art")->value());
    }
    if(req->hasParam("gif")) {
      gifName = req->getParam("gif")->value();
      renderGif = true;
    }
    req->send(200, "text/plain", "OK");
  });

  server.begin();

  //Setup LEDs
  FastLED.addLeds<NEOPIXEL, DATA_PIN>(leds, STRIP_LEN);  // GRB ordering is assumed
  FastLED.setBrightness(0x21);
  FastLED.setCorrection(LEDColorCorrection::TypicalLEDStrip);
  // FastLED.setTemperature(ColorTemperature::ClearBlueSky);
}

unsigned long lastMillis;
void loop() { 
  #ifdef PROD
  dnsServer.processNextRequest();
  #endif
  MDNS.update();

  if (millis() - lastMillis >= 10*1000UL && mode == "cycle") {
    lastMillis = millis();  //get ready for the next iteration
    const uint idx = random(0, ART_COUNT);
    switchGrid(leds, ART_LIST[idx]);
  }
  if(renderGif) {
    playGif(leds, gifName);
  }
}