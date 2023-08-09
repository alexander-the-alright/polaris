# polaris
## _Path Reporting for Borrowed Goods_
_Currently in Pre-Pre-Alpha_
## Installation
### Necessary Software
- [Golang][go]
- [Arduino IDE][ard]
  - [ESP plugin][esp]
  - [TinyGPS++][tg+] library (Installing the TinyGPS++ Library section)

### Necessary Hardware
- Computer (I used a [Raspberry Pi Zero 2 W][vilros])
- [ESP8266][espbuy]
- [uBlox 6m GPS module][gps]
- External Power Supply

## Setup
### Serverside
Compile the recipient Golang program on the computer with
```sh
go build gps.go
```
### Nodeside
Open the sketch and use the Arduino IDE to upload to the ESP module
Connect the GPS module; pinout is given in Arduino sketch, must check against pinout image

## Running
### Serverside
Run the recipient with
```sh
./gps &
```
When the path is received, it will be stored in a file named with the received date. It can be viewed with
```sh
cat YYYYMMDDHHmmss.gps
```
### Nodeside
Compile and download the sketch to the ESP

## Usage
Connect the ESP and GPS to the external battery. Once the GPS warms up, it will get a location reading every 30 seconds. Power does not need to be continuously supplied, as the location is written to the EEPROM. To get the location data, reset/power cycle ESP.

## Plotting
#### _So now you have your GPS coordinates. Now what?_
For now, you have to search the GPS coordinates on Google Maps or something.
However, in the future, it would be nice to process the coordinates locally,
without having to copy/paste several sets of coords to a website. However, the
methods I have found to do this ([HERE][here] and [geoplot][geop]) have not
worked, even the examples. Eventually, I will find another solution, probably
in Python, but until then, copy/paste.

## TODO
- A whole lot
- Golang plotting
- ESP write to EEPROM
- Solve the Power Problem
- Solve the Housing Problem
- Solve the Attachment Problem


   [go]: <https://go.dev/doc/install>
   [ard]: <https://www.arduino.cc/en/software>
   [esp]: <https://randomnerdtutorials.com/how-to-install-esp8266-board-arduino-ide/>
   [tg+]: <https://randomnerdtutorials.com/guide-to-neo-6m-gps-module-with-arduino/>
   [vilros]: <https://www.amazon.com/Vilros-Raspberry-Kit-Premium-Essential-Accessories/dp/B0748M1Z1B>
   [espbuy]: <https://www.amazon.com/HiLetgo-Internet-Development-Wireless-Micropython/dp/B081CSJV2V>
   [gps]: <https://www.amazon.com/HiLetgo-GY-NEO6MV2-Controller-Ceramic-Antenna/dp/B01D1D0F5M>
   [here]: <https://www.here.com/learn/blog/reverse-geocoding-a-location-using-golang>
   [geop]: <https://github.com/morikuni/go-geoplot/tree/mai://github.com/morikuni/go-geoplot/tree/maini>
