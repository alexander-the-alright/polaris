// ===================================================================
// Auth: alex
// File: gps_conn.ino
// Revn: 07-18-2023  1.0
// Func: connect GPS module to Golang's net.Conn via ESP8266
//
// TODO: first loop static bool to distinguish between case 0 and 1
//       store data in EEPROM, retrieve it after power failure
//       entirely ignore 0's, maybe trigger LED on 0's?
// ===================================================================
// CHANGE LOG
// -------------------------------------------------------------------
//*07-18-2023: ac copied from gps_master.ino
//                copied from https://github.com/esp8266/Arduino/blob/master/libraries/ESP8266WiFi/examples/WiFiClient/WiFiClient.ino
//                changed host to byte array
//                commented
//
// ===================================================================



#include <ESP8266WiFi.h>                  // Special Wifi Library
#include <TinyGPS++.h>                    // Tiny GPS Plus Library
// Software Serial Library so we can use other Pins for communication
// with the GPS module
#include <SoftwareSerial.h>

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
    return;
  }

/*
  // wait for data to be available
  unsigned long timeout = millis();
  while (client.available() == 0) {
    if (millis() - timeout > 5000) {
      Serial.println(">>> Client Timeout !");
      client.stop();
      delay(60000);
      return;
    }
  }
*/

  // This will send gps info
  Serial.println("sending gps data");
  if (client.connected()) {
    // send data to server
    client.print(gps.location.lat(), 5);    // 5 decimals of precision
    client.print(",");
    client.println(gps.location.lng(), 5);
    // print data to serial monitor
    Serial.print(gps.location.lat(), 5);    // 5 decimals of precision
    Serial.print(",");
    Serial.println(gps.location.lng(), 5);
  }

  // Close the connection
  Serial.println();
  Serial.println("closing connection");
  client.stop();

  // don't get data continuously, wait one second
  delay(1000);
  
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
