package main

import (
	"time"

	"github.com/yanzay/autohome/modules/arduino"
	"github.com/yanzay/autohome/modules/habra"
	"github.com/yanzay/autohome/modules/weather"
)

type Application struct {
	scheduler *Scheduler
	db        storage
	modules   map[string]Module
}

type Module interface {
	Initialize(map[string]string, interface{})
	Functions() []string
	Settings() []string
	Send(string)
	Handle(string, string, ...string) (string, map[string]string)
	Menus() []string
}

func (a *Application) registerModules() {
	a.modules = make(map[string]Module)
	a.registerModule("arduino", new(arduino.ArduinoModule))
	a.registerModule("weather", new(weather.WeatherModule))
	a.registerModule("habra", new(habra.HabraModule))
}

func (a *Application) registerModule(name string, module Module) {
	opts := a.db.getSettings()
	module.Initialize(opts[name], a.db)
	a.modules[name] = module
}

func main() {
	app := new(Application)
	app.db = NewSqliteStorage()
	app.registerModules()
	app.db.prepareSettings(app.modules)
	app.scheduler = NewScheduler(app)
	go startWebServer(app)
	mainLoop()
}

func mainLoop() {
	for {
		time.Sleep(5 * time.Second)
	}
}
