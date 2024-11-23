# Custom firmware for the Panasonic CZ-TAW1

This is an alternative firmware for RaspberryPi Zero W, an IoT adapter for the H-series heat pumps. It consists of:

* a gateway written in Go that translates serial comms with the heat pump on the CN-CNT link to MQTT

## About

### The gateway

The gateway (called GoHeishaMon or heishamon) is responsible for parsing the data received from the Heat Pump and posting it to MQTT topics. It is a reimplementation of the <https://github.com/Egyras/HeishaMon> project in Go.

#### Gateway features

* posting Heat Pump data to MQTT
* changing settings on the Heat Pump
* supports Home Assistant's MQTT discovery
* emulation of the Optional PCB

GoHeishaMon will be used without the CZ-TAW1 module on a RaspberryPI. It requires a serial port connection to the Heat Pump. The new version is running as a daemon. As a consequence, the logs are no longer written to stdout, they end up in system log [/var/log/messages] (and MQTT topic).

#### Note

* the binary is /usr/bin/heishamon by default
* the configuration is stored in /etc/heishamon/ and is preserved on upgrades
* the service name is **heishamon**

## Installation

If you use a arm64 OS, please download the right go version.

### Prerequisites

* Serial port TTL Adapter connected to the Pi for communication with the pump
* Go Version 1.20.4
* Installed RaspberryOS
* I had to re-arange the PINs of the TTL Adapter with a small breakout board to match the pinout of the pump cable

Overview of the process:

* Install Go 1.20.4 on the PI with the Installation-Script. It will download the right version and install it
* Build GoHeishaMon with Go
```
  cd ./package/heishamon/src  
  go build -v  
```
* Copy all needed files to the system
```  
  cp ./package/heishamon/files/config.example /etc/heishamon/config.yaml  
  cp ./package/heishamon/files/heishamon.init /etc/init.d/heishamon  
  cp ./package/heishamon/files/topics.yaml /etc/heishamon/  
  cp ./package/heishamon/files/topicsOptionalPCB.yaml /etc/heishamon  
  cp ./package/heishamon/src/GoHeishaMon /usr/bin/heishamon  
```
* Make the filex excutable
```  
  chmod +x /etc/init.d/heishamon  
  chmod +x /usr/bin/heishamon  
```
* Enable Service
```  
  update-rc.d heishamon defaults  
```

In order to configure GoHeishaMon:
```
* systemctl stop heishamon
```
* Edit the config file (/etc/heishamon/config.yaml)
```
* systemctl start heishamon  
* systemctl status heishamon  
```
Logs are stored in /var/log/messages  
