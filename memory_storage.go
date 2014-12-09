package main

import (
	"fmt"
)

type memoryStorage struct {
	Settings      settings
	ScheduleItems []scheduleItem
}

func (s *memoryStorage) getSettings() settings {
	fmt.Printf("Getting settings: %v\n", s.Settings)
	return s.Settings
}

func (s *memoryStorage) saveSettings(sets settings) {
	fmt.Printf("Saving settings: %d, %s\n", sets.YahooCityId, sets.YahooTempUnit)
	s.Settings.YahooCityId = sets.YahooCityId
	s.Settings.YahooTempUnit = sets.YahooTempUnit
	fmt.Printf("Settings saved. %v\n", s.Settings)
}

func (s *memoryStorage) getSchedules() []scheduleItem {
	return s.ScheduleItems
}

func (s *memoryStorage) saveSchedules(items []scheduleItem) {
	s.ScheduleItems = items
}
