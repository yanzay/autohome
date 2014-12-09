package main

type settings struct {
	YahooCityId   int64
	YahooTempUnit string
}

type scheduleItem struct {
	CronString string
	FuncName   string
}

type storage interface {
	getSettings() settings
	saveSettings(settings)
	getSchedules() []scheduleItem
	saveSchedules([]scheduleItem)
}
