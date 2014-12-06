package main

import (
	"github.com/yanzay/gorduino"
	"time"
)

func main() {
	arduino := connectArduino()
	go startWebServer(arduino)
	mainLoop()
}

func connectArduino() *gorduino.Gorduino {
	return gorduino.NewGorduino("/dev/tty.usbmodem1411", 13)
}

func mainLoop() {
	for {
		time.Sleep(5 * time.Second)
	}
}
