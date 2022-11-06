package codec

import (
	"io/ioutil"
	"log"
	"sync"
)

const (
	OPTIONAL_QUERY_SIZE  = 19
	PANASONIC_QUERY_SIZE = 110
)

var panasonicQuery = [PANASONIC_QUERY_SIZE]byte{0x71, 0x6c, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var optionalPCBQuery = [OPTIONAL_QUERY_SIZE]byte{0xF1, 0x11, 0x01, 0x50, 0x00, 0x00, 0x40, 0xFF, 0xFF, 0xE5, 0xFF, 0xFF, 0x00, 0xFF, 0xEB, 0xFF, 0xFF, 0x00, 0x00}
var optionalPCBMutex sync.Mutex

var commandChannel chan []byte

func GetChannel() chan []byte {
	if commandChannel == nil {
		commandChannel = make(chan []byte, 20)
	}
	return commandChannel
}

func SendPanasonicQuery() {
	commandChannel <- panasonicQuery[:]
}

func SendOptionalPCBQuery() {
	optionalPCBMutex.Lock()
	toSend := make([]byte, len(optionalPCBQuery))
	copy(toSend, optionalPCBQuery[:])
	optionalPCBMutex.Unlock()

	commandChannel <- toSend
}

func Acknowledge(datagram []byte) {
	optionalPCBMutex.Lock()
	//response to heatpump should contain the data from heatpump on byte 4 and 5
	optionalPCBQuery[4] = datagram[4]
	optionalPCBQuery[5] = datagram[5]
	optionalPCBMutex.Unlock()
}

func SaveOptionalPCB(filename string) {
	optionalPCBMutex.Lock()
	err := ioutil.WriteFile(filename, optionalPCBQuery[:], 0644)
	optionalPCBMutex.Unlock()
	//TODO serialize to json instead, restore topics and []byte
	if err != nil {
		log.Print(err)
	} else {
		log.Print("Optional PCB data stored")
	}
}

func LoadOptionalPCB(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err)
	} else {
		optionalPCBMutex.Lock()
		copy(optionalPCBQuery[:], data)
		optionalPCBMutex.Unlock()

		log.Print("Optional PCB data loaded")
	}
}