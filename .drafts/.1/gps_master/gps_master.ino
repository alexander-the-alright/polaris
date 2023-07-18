// ===================================================================
// Auth: alex
// File: gps_master.ino
// Revn: 07-17-2023  1.0
// Func: show how to use Ublox GPS module
//
// TODO: draft 2
// ===================================================================
// CHANGE LOG
// -------------------------------------------------------------------
//*07-17-2023: ac adapted from https://github.com/mkconer/ESP8266_GPS/blob/master/ESP8266_GPS_OLED_Youtube.ino
//
// ===================================================================


#include <TinyGPS++.h>                    // Tiny GPS Plus Library
// Software Serial Library so we can use other Pins for communication
// with the GPS module
#include <SoftwareSerial.h>

// set ESP Rx to GPIO 12 (connect to Ublox Tx),
// Tx to GPIO 13 (connect to Ublox Rx)
// GPIO pinout image at URL: https://m.media-amazon.com/images/I/812vnqUs5tL.jpg
// Should also have an image in git repository
// Ublox 6m GPS module to pins 12 and 13
static const int RXPin = 12, TXPin = 13;
// Ublox GPS default Baud Rate is 9600
static const uint32_t GPSBaud = 9600;

// Create an Instance of the TinyGPS++ object called gps
TinyGPSPlus gps;
// The serial connection to the GPS device
SoftwareSerial ss(RXPin, TXPin);

void setup() {
  // Set ESP and Software Serial comms speed
  Serial.begin(GPSBaud);
  ss.begin(GPSBaud);
}

void loop() {
  // print location information
  Serial.print("Latitude  : ");
  Serial.println(gps.location.lat(), 5);    // 5 decimals of precision
  Serial.print("Longitude : ");
  Serial.println(gps.location.lng(), 5);
  
  delay(2000);              // wait 2 seconds between getting location
  
  smartDelay(500);          // Run Procedure smartDelay

  // if GPS is not connected correctly
  if (millis() > 5000 && gps.charsProcessed() < 10)
    Serial.println(F("No GPS data received: check wiring"));
}

// Custom version of delay() ensures gps object is being "fed"
static void smartDelay(unsigned long ms) {              
  unsigned long start = millis();
  do {    // wait until GPS has data
    while (ss.available())
      gps.encode(ss.read());
  } while (millis() - start < ms);
}
