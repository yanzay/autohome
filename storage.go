package main

// type settings struct {
// 	YahooCityId   int64
// 	YahooTempUnit string
// }

type settings map[string]map[string]string

type scheduleItem struct {
  CronString string
  Module     string
  FuncName   string
}

type tempItem struct {
  Datetime string
  Value    float32
}

type storage interface {
  getSettings() settings
  saveSettings(settings)
  getSchedules() []scheduleItem
  saveSchedules([]scheduleItem)
  prepareSettings(map[string]Module)
}
