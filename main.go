package main

import (
	"fmt"
	"github.com/yanzay/gorduino"
	"time"
)

type Application struct {
	scheduler *Scheduler
	db        storage
}

func main() {
	app := new(Application)
	fmt.Println("connecting")
	arduino := connectArduino()
	fmt.Println("connected")
	app.db = NewSqliteStorage()
	app.scheduler = NewScheduler(arduino, app.db)
	go startWebServer(arduino, app.db, app)
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
