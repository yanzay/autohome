package main

import (
  "fmt"
  "github.com/eaigner/jet"
  _ "github.com/mattn/go-sqlite3"
  "strconv"
  "strings"
  "time"
)

type sqliteStorage struct {
  db      *jet.Db
  modules map[string]Module
}

func NewSqliteStorage() *sqliteStorage {
  fmt.Println("Creating new sqlite storage")
  var err error
  store := new(sqliteStorage)
  store.db, err = jet.Open("sqlite3", "./data.db")
  if err != nil {
    fmt.Printf("Error connecting to database %s\n", err)
  }
  fmt.Println("Creating databases")
  store.createDatabase()
  return store
}

func (s *sqliteStorage) createDatabase() {
  s.db.Query("CREATE TABLE IF NOT EXISTS schedules (cron_string TEXT, module TEXT, func_name TEXT)").Run()
  s.db.Query("CREATE TABLE IF NOT EXISTS settings (module TEXT, key TEXT, value TEXT)").Run()
  s.db.Query("CREATE TABLE IF NOT EXISTS arduino_temperatures (datetime TEXT, value REAL)").Run()
}

func (s *sqliteStorage) prepareSettings(modules map[string]Module) {
  s.modules = modules
  var rows []*struct {
    Module string
    Key    string
    Value  string
  }

  for name, module := range s.modules {
    for _, setting := range module.Settings() {
      query := fmt.Sprintf("SELECT module, key, value FROM settings WHERE module='%s' AND key='%s'", name, setting)
      s.db.Query(query).Rows(&rows)
      if len(rows) == 0 {
        insert := fmt.Sprintf("INSERT INTO settings (module, key, value) VALUES ('%s', '%s', '')", name, setting)
        s.db.Query(insert).Run()
      }
    }
  }
}

func (s *sqliteStorage) getSettings() settings {
  sets := make(settings)
  var rows []*struct {
    Module string
    Key    string
    Value  string
  }
  err := s.db.Query("SELECT module, key, value FROM settings").Rows(&rows)
  for _, item := range rows {
    if sets[item.Module] == nil {
      sets[item.Module] = make(map[string]string)
    }
    sets[item.Module][item.Key] = item.Value
  }
  if err != nil {
    fmt.Printf("Error getting settings from database %s\n", err)
  }
  return sets
}

func (s *sqliteStorage) saveSettings(sets settings) {
  s.db.Query("delete from settings").Run()
  for module, setting := range sets {
    for key, value := range setting {
      err := s.db.Query("INSERT INTO settings (module, key, value) values ($1, $2, $3)", module, key, value).Run()
      if err != nil {
        fmt.Printf("Error saving: %v", err)
      }
    }
  }
}

func (s *sqliteStorage) getSchedules() []scheduleItem {
  var items []scheduleItem
  err := s.db.Query("select cron_string, module, func_name from schedules").Rows(&items)
  if err != nil {
    fmt.Printf("Error getting schedule items from database %s\n", err)
  }
  return items
}

func (s *sqliteStorage) saveSchedules(items []scheduleItem) {
  var values []string
  err := s.db.Query("delete from schedules").Run()
  if err != nil {
    fmt.Printf("Error deleting schedules: %s", err)
  }
  query := "insert into schedules (cron_string, module, func_name) values "
  for _, item := range items {
    values = append(values, fmt.Sprintf("('%s', '%s', '%s')", item.CronString, item.Module, item.FuncName))
  }
  query += strings.Join(values, ", ")
  err = s.db.Query(query).Run()
  if err != nil {
    fmt.Printf("Error saving schedules: %v\n", err)
  }
}

func (s *sqliteStorage) SaveTemperature(temp float32) {
  if temp > 1 {
    datetime := time.Now().Format("2006-01-02T15:04:05")
    query := fmt.Sprintf("INSERT INTO arduino_temperatures (datetime, value) values ('%s', %f)\n", datetime, temp)
    err := s.db.Query(query).Run()
    if err != nil {
      fmt.Printf("Fail to save temperature: %v\n", err)
    }
  }
}

func (s *sqliteStorage) GetTemperatures(lastDate ...string) [][2]string {
  var items []tempItem
  var query string
  if len(lastDate) > 0 && lastDate != nil {
    t, _ := time.Parse("2006-01-02T15:04:05.000Z", lastDate[0])
    formatted := t.Format("2006-01-02T15:04:05")
    query = fmt.Sprintf("SELECT datetime, value FROM arduino_temperatures WHERE datetime > '%s'", formatted)
  } else {
    query = "SELECT datetime, value FROM (SELECT * FROM arduino_temperatures ORDER BY datetime DESC LIMIT 120) ORDER BY datetime"
  }
  err := s.db.Query(query).Rows(&items)
  if err != nil {
    fmt.Printf("Error getting temperatures from database %s\n", err)
  }
  var result [][2]string
  result = make([][2]string, len(items))
  for i, item := range items {
    result[i][0] = item.Datetime
    result[i][1] = strconv.FormatFloat(float64(item.Value), 'f', 4, 32)
  }
  return result
}
