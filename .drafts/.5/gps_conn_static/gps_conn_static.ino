// ===================================================================
// Auth: alex
// File: gps_conn_static.ino
// Revn: 08-01-2023  2.0
// Func: connect GPS module to Golang's net.Conn via ESP8266
//
// TODO: first loop static bool to distinguish between case 0 and 1
//       store data in EEPROM, retrieve it after power failure
// ===================================================================
// CHANGE LOG
// -------------------------------------------------------------------
//*07-20-2023: ac copied from gps_conn.ino
//                added switch/case 
// 07-25-2023: ac changed start/fin to ot/za
//                gutted client timeout block, already commented
//                added EEPROM.h
//                removed code block that closes client connection
// 07-27-2023: ac added spin loop to spin until location is updated
//                  (very broken)
//*08-01-2023: ac removed spin loop, added isUpdated() condition to
//                  client sending block
//
// ===================================================================


#include <ESP8266WiFi.h>                  // Special Wifi Library
#include <TinyGPS++.h>                    // Tiny GPS Plus Library
// Software Serial Library so we can use other Pins for communication
// with the GPS module
#include <SoftwareSerial.h>
#include "EEPROM.h"


// wifi credentials
#ifndef STASSID
#define STASSID "kaer morhen"
#define STAPSK "project:sentinel"
#endif


// set ESP Rx to GPIO 12 (connect to Ublox Tx),
// Tx to GPIO 13 (connect to Ublox Rx)
// GPIO pinout image at URL: https://m.media-amazon.com/images/I/812vnqUs5tL.jpg
// Should also have an image in git repository
// Ublox 6m GPS module to pins 12 and 13
static const int RXPin = 12, TXPin = 13;
// Ublox GPS default Baud Rate is 9600
static const uint32_t GPSBaud = 9600;

// wifi credentials
const char* ssid = STASSID;
const char* password = STAPSK;

const unsigned char host[4] = { 192, 168, 1, 120 };
//const char* host = "xyz.xyz.xyz.xyz";
const uint16_t port = 1202;

// literally keep track of how many messages have been sent
int count = 0;

// Create an Instance of the TinyGPS++ object called gps
TinyGPSPlus gps;
// The serial connection to the GPS device
SoftwareSerial ss(RXPin, TXPin);

void setup() {
  // Set ESP and Software Serial comms speed
  Serial.begin(115200);
  ss.begin(GPSBaud);

  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  /* Explicitly set the ESP8266 to be a WiFi-client, otherwise, it by default,
     would try to act as both a client and an access-point and could cause
     network-issues with your other WiFi-devices on your WiFi-network. */
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  // wait for connection to be established
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  // print connection details
  Serial.println("");
  Serial.println("WiFi connected");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}


void loop() {

  // print connecty things every loop
  // it disconnects every loop. Don't love this tbh
  Serial.print("connecting to ");
  //Serial.print(host);
  Serial.print(':');
  Serial.println(port);

  // Use WiFiClient class to create TCP connections
  WiFiClient client;
  if (!client.connect(host, port)) {
    Serial.println("connection failed");
    delay(5000);
    //return;
  }

  Serial.print( "isUpdated: " );
  Serial.println( gps.location.isUpdated() );
  Serial.print( "isValid: " );
  Serial.println( gps.location.isValid() );
  delay( 1000 * 30 );
  //     1000 milliseconds in one second
  //            30 seconds in one-half minute
  // }

  // This will send gps info
  Serial.println("sending gps data");
  // only send message if client is connected and location is valid
  if (client.connected() && gps.location.isValid()) {
    switch( count ) {
      case 0:     // on the first message
        // send start message
        client.println("ot");
        Serial.println("start/ot");
        count++;
        break;
      case 15:    // on the last message
        // send finale message
        client.println("za");
        Serial.println("fin/za");
        count = 0;    // restart count
        break;
      default:    // data messages
        // send data to server
        client.print(gps.location.lat(), 6);    // 6 decimals of precision
        client.print(",");
        client.println(gps.location.lng(), 6);
        // print data to serial monitor
        Serial.print(gps.location.lat(), 6);    // 6 decimals of precision
        Serial.print(",");
        Serial.println(gps.location.lng(), 6);
        count++;
        break;
    }
  }
  
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
