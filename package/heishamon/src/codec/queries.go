package codec

import (
	"sync"

	"github.com/rondoval/GoHeishaMon/topics"
)

const (
	optionalQuerySize  = 19
	panasonicQuerySize = 110
)

var panasonicQuery = [panasonicQuerySize]byte{0x71, 0x6c, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var optionalPCBQuery = [optionalQuerySize]byte{0xF1, 0x11, 0x01, 0x50, 0x00, 0x00, 0x40, 0xFF, 0xFF, 0xE5, 0xFF, 0xFF, 0x00, 0xFF, 0xEB, 0xFF, 0xFF, 0x00, 0x00}
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

func RestoreOptionalPCB(optinalTopics []*topics.TopicEntry) {
	optionalPCBMutex.Lock()
	for _, sensor := range optinalTopics {
		if sensor.EncodeFunction != "" && sensor.CurrentValue() != "" {
			encode(sensor, optionalPCBQuery[:])
		}
	}
	optionalPCBMutex.Unlock()
}
