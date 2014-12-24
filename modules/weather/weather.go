package weather

import (
  "bytes"
  "fmt"
  "github.com/yanzay/autohome/modules/generic"
  "github.com/yanzay/googlespeak"
  "github.com/yanzay/yahooweather"
  "strconv"
  "text/template"
)

type WeatherModule struct {
  generic.GenericModule
}

func (m *WeatherModule) Settings() []string {
  return []string{"city_id", "temp_unit"}
}

func (w *WeatherModule) GetFunctions() []func(map[string]string) {
  return []func(map[string]string){w.Notify}
}

func (w *WeatherModule) Notify(settings map[string]string) {
  cityId, _ := strconv.ParseInt(settings["city_id"], 10, 64)
  condition, forecasts, err := yahooweather.GetWeather(cityId)
  if err != nil {
    fmt.Printf("Error getting weather: %s\n", err)
    return
  }

  var weather = struct {
    currentTemp  int
    tomorrowTemp int
  }{
    condition.Temp,
    forecasts[0].Low,
  }

  var buf bytes.Buffer
  tmpl, _ := template.New("text").Parse(settings["text"])
  tmpl.Execute(&buf, weather)
  text := buf.String()

  googlespeak.Say(text, settings["lang"])
}

func (w *WeatherModule) Send(command string) {

}
