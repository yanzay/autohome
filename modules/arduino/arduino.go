package arduino

import (
  "encoding/binary"
  "fmt"
  "github.com/yanzay/autohome/modules/generic"
  "github.com/yanzay/automata"
  "math"
  "regexp"
  "strconv"
  "strings"
)

type ArduinoModule struct {
  generic.GenericModule
  arduino *automata.Arduino
  db      TemperatureGetSaver
  options map[string]string
}

type TemperatureGetSaver interface {
  SaveTemperature(float32)
  GetTemperatures(lastDate ...string) [][2]string
}

type tempItem struct {
  Datetime string
  Value    float32
}

func (a *ArduinoModule) Initialize(options map[string]string, db interface{}) {
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
  a.db = db.(TemperatureGetSaver)
  a.options = options
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
  // pin, _ := regexp.Compile(".*_([0-9]+)")
  var p int

  matches := com.FindStringSubmatch(command)
  if len(matches) > 2 {
    command = matches[1]
    fmt.Printf(command)
    p, _ = strconv.Atoi(matches[2])
    fmt.Printf(" with pin %v\n", byte(p))
  }
  // if withPin != "" {
  //   p, _ = strconv.Atoi(pin.FindString(command))
  //   command = withPin
  // }
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
      a.db.SaveTemperature(value)
    }
  }()
}

func (a *ArduinoModule) Handle(method, action string, params ...string) (string, map[string]string) {
  if method == "GET" {
    switch action {
    case "control":
      return "arduino/control", map[string]string{}
    case "stats":
      return "arduino/stats", a.getStats()
    case "last_stats":
      return "", a.getLastStats(params)
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
  stats := make(map[string]string)
  for _, item := range a.db.GetTemperatures() {
    stats[item[0]] = item[1]
  }
  return stats
}

func (a *ArduinoModule) getLastStats(params []string) map[string]string {
  stats := make(map[string]string)
  for _, item := range a.db.GetTemperatures(params[0]) {
    stats[item[0]] = item[1]
  }
  return stats
}
