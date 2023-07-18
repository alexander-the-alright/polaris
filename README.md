# polaris
## _Path Reporting for Borrowed Goods_

## Installation
### Necessary Software
- [Golang][go]
- [Arduino IDE][ard] and [ESP plugin][esp]

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
Open the sketch and user the Arduino IDE to upload to the ESP module
Connect the GPS module; pinout is given in Arduino sketch, must check against pinout image

## Running
### Serverside
Run the recipient with
```sh
./gps &
```
When the path is received, it will be stored in a file named with the received date. It can be viewed with
```sh
cat YYYYMMDD.path
```
### Nodeside
Compile and download the sketch to the ESP

## Usage
Connect the ESP and GPS to the external battery. Once the GPS warms up, it will get a location reading every 2 seconds. Power does not need to be continuously supplied, as the location is written to the EEPROM. To get the location data, reset/power cycle ESP.

## TODO
- A whole lot
- Golang write to file
- Golang plotting
- ESP connectivity
- ESP write to EEPROM


   [go]: <https://go.dev/doc/install>
   [ard]: <https://www.arduino.cc/en/software>
   [esp]: <https://randomnerdtutorials.com/how-to-install-esp8266-board-arduino-ide/>
   [vilros]: <https://www.amazon.com/Vilros-Raspberry-Kit-Premium-Essential-Accessories/dp/B0748M1Z1B>
   [espbuy]: <https://www.amazon.com/HiLetgo-Internet-Development-Wireless-Micropython/dp/B081CSJV2V>
   [gps]: <https://www.amazon.com/HiLetgo-GY-NEO6MV2-Controller-Ceramic-Antenna/dp/B01D1D0F5M>