package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/yanzay/gorduino"
	"net/http"
)

type Application struct {
	Storage storage
}

type Stats struct {
	Message string
}

type StatsForm struct {
	YahooCityId   int64  `form:"YahooCityId" binding:"required"`
	YahooTempUnit string `form:"YahooTempUnit" binding:"required"`
}

func startWebServer(arduino *gorduino.Gorduino) {
	db := new(memoryStorage)
	m := martini.Classic()
	m.Map(db)
	m.Map(arduino)
	m.Use(render.Renderer(render.Options{Layout: "layout"}))
	m.Get("/", statsHandler)
	m.Get("/stats", statsHandler)
	m.Get("/control", controlHandler)
	m.Post("/control", postControlHandler)
	m.Get("/settings", settingsHandler)
	m.Post("/settings", binding.Bind(StatsForm{}), saveSettingsHandler)
	m.Run()
}

func statsHandler(r render.Render) {
	r.HTML(200, "stats", Stats{Message: "hey ho"})
}

func controlHandler(r render.Render, db storage) {
	r.HTML(200, "control", nil)
}

func postControlHandler(w http.ResponseWriter, r *http.Request, arduino *gorduino.Gorduino) {
	arduino.Toggle(13)
	http.Redirect(w, r, "/control", 302)
}

func settingsHandler(r render.Render, db storage) {
	r.HTML(200, "settings", db.getSettings())
}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request, f StatsForm, db storage) {
	fmt.Printf("form: %v", f)
	sets := settings{YahooCityId: f.YahooCityId, YahooTempUnit: f.YahooTempUnit}
	fmt.Printf("request: %v", sets)
	db.saveSettings(sets)
	http.Redirect(w, r, "/settings", 302)
}
