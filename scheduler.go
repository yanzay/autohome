package main

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/yanzay/googlespeak"
	"github.com/yanzay/gorduino"
	"github.com/yanzay/yahooweather"
)

type Scheduler struct {
	funcs   []func()
	cr      *cron.Cron
	arduino *gorduino.Gorduino
	db      storage
}

func NewScheduler(arduino *gorduino.Gorduino, db storage) *Scheduler {
	s := new(Scheduler)
	s.arduino = arduino
	s.db = db
	s.Reload()
	return s
}

func (s *Scheduler) Reload() {
	if s.cr != nil {
		s.cr.Stop()
	}
	s.cr = cron.New()
	s.loadFunctions(s.cr)
	s.cr.Start()
}

func (s *Scheduler) loadFunctions(cr *cron.Cron) {
	items := s.db.getSchedules()
	for _, item := range items {
		cr.AddFunc(item.CronString, s.funcByName(item.FuncName))
	}
	// cr.AddFunc("0 1 14 * * *", s.funcByName("lightOn"))
	// cr.AddFunc("0 43 13 * * *", s.funcByName("lightOff"))
	// cr.AddFunc("00 45 13 * * *", s.funcByName("notify"))
}

func (s *Scheduler) lightOff() {
	fmt.Println("Light off")
	s.arduino.Off(13)
}

func (s *Scheduler) lightOn() {
	fmt.Println("Light on")
	s.arduino.On(13)
}

func (s *Scheduler) notify() {
	fmt.Println("Notify!")
	sets := s.db.getSettings()
	condition, forecasts, err := yahooweather.GetWeather(sets.YahooCityId)
	if err != nil {
		fmt.Printf("Error getting weather: %s\n", err)
		return
	}

	text := welcomeText()
	text += conditionText(condition, forecasts[0])

	googlespeak.Say(text, "ru")
}

func (s *Scheduler) funcByName(name string) func() {
	switch name {
	case "lightOn":
		return s.lightOn
	case "lightOff":
		return s.lightOff
	case "notify":
		return s.notify
	}
	return func() {}
}
