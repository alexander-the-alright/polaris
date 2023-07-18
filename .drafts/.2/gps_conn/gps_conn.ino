// ===================================================================
// Auth: alex
// File: gps_conn.ino
// Revn: 07-17-2023  1.0
// Func: connect GPS module to Golang's net.Conn via ESP8266
//
// TODO: add connectivity
// ===================================================================
// CHANGE LOG
// -------------------------------------------------------------------
//*07-17-2023: ac copied from gps_master.ino
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
