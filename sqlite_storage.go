package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type sqliteStorage struct {
	db *sql.DB
	// Settings      settings
	// ScheduleItems []scheduleItem
}

func NewSqliteStorage() *sqliteStorage {
	fmt.Println("Creating new sqlite storage")
	var err error
	store := new(sqliteStorage)
	store.db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Printf("Error connecting to database %s\n", err)
	}
	fmt.Println("Creating databases")
	store.createDatabase()
	return store
}

func (s *sqliteStorage) createDatabase() {
	query := `
  create table if not exists settings (yahoo_city_id integer, yahoo_temp_unit char(1));
  create table if not exists schedules (cron_string text, func_name text);
  `
	s.db.Exec(query)
}

func (s *sqliteStorage) getSettings() settings {
	var sets settings
	rows, err := s.db.Query("select yahoo_city_id, yahoo_temp_unit from settings")
	if err != nil {
		fmt.Printf("Error getting settings from database %s\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cityId int64
		var tempUnit string
		rows.Scan(&cityId, &tempUnit)
		fmt.Println(cityId, tempUnit)
		sets = settings{YahooCityId: cityId, YahooTempUnit: tempUnit}
	}
	return sets
}

func (s *sqliteStorage) saveSettings(sets settings) {
	s.db.Exec("delete from settings")
	fmt.Printf("Saving settings: %d, %s\n", sets.YahooCityId, sets.YahooTempUnit)
	query := fmt.Sprintf("insert into settings (yahoo_city_id, yahoo_temp_unit) values (%d, '%s')", sets.YahooCityId, sets.YahooTempUnit)
	fmt.Printf("query %s\n", query)
	s.db.Exec(query)
}

func (s *sqliteStorage) getSchedules() []scheduleItem {
	fmt.Println("getting schedules")
	var items []scheduleItem
	rows, err := s.db.Query("select cron_string, func_name from schedules")
	if err != nil {
		fmt.Printf("Error getting schedule items from database %s\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println("just another row...")
		var cronString string
		var funcName string
		rows.Scan(&cronString, &funcName)
		fmt.Println(cronString, funcName)
		items = append(items, scheduleItem{CronString: cronString, FuncName: funcName})
	}
	return items
}

func (s *sqliteStorage) saveSchedules(items []scheduleItem) {
	var values []string
	_, err := s.db.Exec("delete from schedules")
	if err != nil {
		fmt.Printf("Error deleting schedules: %s", err)
	}
	fmt.Println("Saving schedules")
	query := "insert into schedules (cron_string, func_name) values "
	for _, item := range items {
		values = append(values, fmt.Sprintf("('%s', '%s')", item.CronString, item.FuncName))
	}
	query += strings.Join(values, ", ")
	fmt.Printf("Final query: %s", query)
	_, err = s.db.Exec(query)
	if err != nil {
		fmt.Printf("Error saving schedules: %v\n", err)
	}
	fmt.Println("schedules saved")
}
