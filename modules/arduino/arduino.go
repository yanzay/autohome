package arduino

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/autohome/modules/generic"
	"github.com/yanzay/automata"
	"github.com/yanzay/reactiverecord"
)

type ArduinoModule struct {
	generic.GenericModule
	arduino *automata.Arduino
	db      *reactiverecord.ReactiveRecord
	options map[string]string
}

type Creater interface {
	Create(interface{}) error
}

type tempItem struct {
	Datetime string
	Value    float32
}

func (a *ArduinoModule) Initialize(options map[string]string, db interface{}) {
	var err error
	a.db, err = reactiverecord.Connect("sqlite", "data.db")
	a.db.CreateTable(ArduinoTemperature{})
	// reactiverecord.CreateTable(ArduinoTemperature{})
	a.options = options
	// ar, err := automata.NewSerial("/dev/tty.usbmodem1411")
	ar, err := automata.New(automata.EthernetArduino, "192.168.0.13:13666")
	if err != nil {
		fmt.Printf("Can't connect to arduino, continue without it. Error: %v\n", err)
		return
	}
	pins := strings.Split(options["digital_out_pins"], ",")
	for _, pin := range pins {
		ipin, _ := strconv.Atoi(pin)
		ar.SetDigitalOutput(byte(ipin))
	}
	a.arduino = ar
}

func (a *ArduinoModule) Functions() []string {
	var funcs []string
	pins := strings.Split(a.options["digital_out_pins"], ",")
	for _, pin := range pins {
		funcs = append(funcs, "lightOn_"+pin)
		funcs = append(funcs, "lightOff_"+pin)
		funcs = append(funcs, "lightToggle_"+pin)
	}
	funcs = append(funcs, "getTemp")
	return funcs
}

func (a *ArduinoModule) Settings() []string {
	return []string{"port", "digital_out_pins"}
}

func (a *ArduinoModule) Send(command string) {
	if a.arduino == nil {
		fmt.Printf("Can't operate on Arduino, it's missed. Command: %s\n", command)
		return
	}
	com, _ := regexp.Compile("(.*)_([0-9]+)")
	var p int

	matches := com.FindStringSubmatch(command)
	if len(matches) > 2 {
		command = matches[1]
		p, _ = strconv.Atoi(matches[2])
	}

	go func() {
		switch command {
		case "lightOn":
			a.arduino.On(byte(p))
		case "lightOff":
			a.arduino.Off(byte(p))
		case "lightToggle":
			a.arduino.Toggle(byte(p))
		case "getTemp":
			b := a.arduino.Temp()
			b = append(b)
			bits := binary.LittleEndian.Uint32(b)
			value := math.Float32frombits(bits)
			a.saveTemperature(value)
			// a.db.SaveTemperature(value)
		}
	}()
}

func (a *ArduinoModule) Handle(method, action string, params url.Values) (string, map[string]string) {
	if method == "GET" {
		switch action {
		case "control":
			return "arduino/control", map[string]string{}
		case "stats":
			period := params.Get("period")
			return "arduino/stats", a.getLastStatsForPeriod(period, "")
		case "last_stats":
			return "", a.getLastStatsForPeriod(params.Get("period"), params.Get("lastDate"))
		}
	} else {
		switch action {
		case "control":
			a.Send("lightToggle")
		}
	}
	return "", map[string]string{}
}

func (a *ArduinoModule) Menus() []string {
	return []string{"control", "stats"}
}

func (a *ArduinoModule) getStats() map[string]string {
	var stats []ArduinoTemperature
	result := make(map[string]string)
	a.db.All(ArduinoTemperature{}, &stats)
	for _, stat := range stats {
		result[stat.DateTime] = fmt.Sprintf("%v", stat.Value)
	}
	return result
}

func (a *ArduinoModule) getLastStatsForPeriod(period string, lastDate string) map[string]string {
	var stats []ArduinoTemperature
	result := make(map[string]string)
	var duration time.Duration
	groupFormat := "%Y-%m-%dT%H:%M:%S"
	switch period {
	case "5m":
		duration = 5 * time.Minute
	case "1h":
		duration = 1 * time.Hour
		groupFormat = "%Y-%m-%dT%H:%M:00"
	case "1d":
		duration = 24 * time.Hour
		groupFormat = "%Y-%m-%dT%H:00:00"
	case "30d":
		duration = 30 * 24 * time.Hour
		groupFormat = "%Y-%m-%dT00:00:00"
	case "1Y":
		duration = 365 * 30 * 24 * time.Hour
		groupFormat = "%Y-%m-01T00:00:00"
	default:
		duration = 5 * time.Minute
	}
	var start time.Time
	var err error
	if lastDate != "" {
		start, err = time.Parse("2006-01-02T15:04:05.000Z", lastDate)
		if err != nil {
			log.Println(err)
			start = time.Now().Add(-duration)
		}
	} else {
		start = time.Now().Add(-duration)
	}
	query := fmt.Sprintf("SELECT STRFTIME('%s', DateTime) DateTime, AVG(Value) Value FROM arduino_temperature WHERE STRFTIME('%s', DateTime) > STRFTIME('%s', '%s') GROUP BY STRFTIME('%s', DateTime)", groupFormat, groupFormat, groupFormat, start.Format("2006-01-02T15:04:05"), groupFormat)
	a.db.Query(query).Rows(&stats)
	for _, stat := range stats {
		result[stat.DateTime] = fmt.Sprintf("%v", stat.Value)
	}
	return result
}

func (a *ArduinoModule) getLastStats(lastDate, period string) map[string]string {
	var stats []ArduinoTemperature
	result := make(map[string]string)
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", lastDate)
	formatted := t.Format("2006-01-02T15:04:05")
	a.db.Where(ArduinoTemperature{}, fmt.Sprintf("DateTime > '%s'", formatted)).Run(&stats)
	for _, stat := range stats {
		result[stat.DateTime] = fmt.Sprintf("%v", stat.Value)
	}
	return result
}

func (a *ArduinoModule) saveTemperature(value float32) {
	at := new(ArduinoTemperature)
	at.DateTime = time.Now().Format("2006-01-02T15:04:05")
	at.Value = value
	a.db.Create(*at)
}

func hourAgo() time.Time {
	return time.Now().Add(-1 * time.Hour)
}
