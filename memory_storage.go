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
  s.Settings = sets
}

func (s *memoryStorage) getSchedules() []scheduleItem {
  return s.ScheduleItems
}

func (s *memoryStorage) saveSchedules(items []scheduleItem) {
  s.ScheduleItems = items
}
