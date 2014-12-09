package main

import (
	"fmt"
	"github.com/yanzay/yahooweather"
	"time"
)

func welcomeText() string {
	return "Привет!"
}

func conditionText(cond yahooweather.ConditionItem, today yahooweather.ForecastItem) string {
	var text string
	currentTime := time.Now().Format("00:00")

	text = fmt.Sprintf("Киевское время %s. ", currentTime)
	text += fmt.Sprintf("Сегодня %s, %d-е %s, в Киеве ожидается %s. ", dayOfWeek(), time.Now().Day(), month(), condition(today.Code))
	text += fmt.Sprintf("Температура воздуха от %d до %d градусa цельсия.", today.Low, today.High)
	return text
}

func dayOfWeek() string {
	switch time.Now().Weekday() {
	case time.Monday:
		return "понедельник"
	case time.Tuesday:
		return "вторник"
	case time.Wednesday:
		return "среда"
	case time.Thursday:
		return "четверг"
	case time.Friday:
		return "пятница"
	case time.Saturday:
		return "суббота"
	case time.Sunday:
		return "воскресенье"
	}
	return ""
}

func month() string {
	switch time.Now().Month() {
	case time.January:
		return "января"
	case time.February:
		return "февраля"
	case time.March:
		return "марта"
	case time.April:
		return "апреля"
	case time.May:
		return "мая"
	case time.June:
		return "июня"
	case time.July:
		return "июля"
	case time.August:
		return "августа"
	case time.September:
		return "сентября"
	case time.October:
		return "октября"
	case time.November:
		return "ноября"
	case time.December:
		return "декабря"
	}
	return ""
}

func condition(code int) string {
	fmt.Println("Code: %d", code)
	conditions := map[int]string{
		0:  "торнадо",
		1:  "тропический шторм",
		2:  "ураган",
		5:  "снег с дождём",
		6:  "дождь с мокрым снегом",
		7:  "мокрый снег",
		10: "дождь с градом",
		11: "дождь",
		12: "дождь",
		16: "снег",
		26: "облачная погода",
		28: "облачная погода с прояснениями",
		29: "облачная погода с прояснениями",
		30: "облачная погода с прояснениями",
	}
	return conditions[code]
}
