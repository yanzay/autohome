package main

type settings struct {
	YahooCityId   int64
	YahooTempUnit string
}

type storage interface {
	getSettings() settings
	saveSettings(settings)
}
