package main

import (
  "fmt"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "html/template"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
)

type Stats struct {
  Message string
}

type StatsForm struct {
  YahooCityId   string `form:"YahooCityId" binding:"required"`
  YahooTempUnit string `form:"YahooTempUnit" binding:"required"`
}

type SchedulerForm struct {
  CronStrings []string `form:"cronStrings[]" binding:"required"`
  FuncNames   []string `form:"funcNames[]" binding:"required"`
}

func startWebServer(app *Application) {
  m := martini.Classic()
  m.Map(app.db)
  m.Map(app.modules["arduino"])
  m.Map(app)
  options := render.Options{
    Layout: "layout",
    Funcs: []template.FuncMap{
      template.FuncMap{
        "menus":     func() interface{} { return moduleMenus(app) },
        "functions": func() interface{} { return moduleFunctions(app) },
      },
    },
  }
  m.Use(render.Renderer(options))
  m.Get("/", settingsHandler)
  m.Get("/settings", settingsHandler)
  m.Post("/settings", saveSettingsHandler)
  m.Get("/scheduler", schedulerHandler)
  m.Post("/scheduler", binding.Bind(SchedulerForm{}), saveSchedulerHandler)
  m.Get("/:module/:action", getModuleHandler)
  m.Post("/:module/:action", postModuleHandler)
  m.Run()
}

type MenuItem struct {
  Title string
  Link  string
}

type ModuleFunction struct {
  Module   string
  FuncName string
}

func moduleMenus(app *Application) []MenuItem {
  var menus []MenuItem
  for name, module := range app.modules {
    for _, menu := range module.Menus() {
      menus = append(menus, MenuItem{Title: strings.Title(menu), Link: name + "/" + menu})
    }
  }
  return menus
}

func moduleFunctions(app *Application) []ModuleFunction {
  var funcs []ModuleFunction
  for name, module := range app.modules {
    for _, funcname := range module.Functions() {
      funcs = append(funcs, ModuleFunction{Module: name, FuncName: funcname})
    }
  }
  return funcs
}

func getModuleHandler(req *http.Request, params martini.Params, r render.Render, app *Application) {
  module := app.modules[params["module"]]
  if module != nil {
    template, data := module.Handle("GET", params["action"], req.URL.Query()["lastDate"]...)
    if template != "" {
      r.HTML(200, template, data)
    } else {
      r.JSON(200, data)
    }
  }
}

func postModuleHandler(w http.ResponseWriter, r *http.Request, app *Application, params martini.Params) {
  module := app.modules[params["module"]]
  module.Handle("POST", params["action"])
  route := fmt.Sprintf("/%s/%s", params["module"], params["action"])
  http.Redirect(w, r, route, 302)
}

func settingsHandler(r render.Render, db storage, app *Application) {
  sets := make(map[string][]string)
  for name, module := range app.modules {
    sets[name] = module.Settings()
  }
  r.HTML(200, "settings", db.getSettings())
}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request, db storage, app *Application) {
  body, _ := ioutil.ReadAll(r.Body)
  v, _ := url.ParseQuery(string(body))
  sets := make(settings)
  for name, module := range app.modules {
    sets[name] = make(map[string]string)
    for _, setting := range module.Settings() {
      form_key := fmt.Sprintf("%s[%s]", name, setting)
      sets[name][setting] = v.Get(form_key)
      fmt.Printf("sets[%s][%s] = %s\n", name, setting, sets[name][setting])
    }
    module.Initialize(sets[name], db)
  }
  db.saveSettings(sets)
  http.Redirect(w, r, "/settings", 302)
}

func schedulerHandler(r render.Render, db storage) {
  r.HTML(200, "scheduler", db.getSchedules())
}

func saveSchedulerHandler(w http.ResponseWriter, r *http.Request, f SchedulerForm, db storage, app *Application) {
  var items []scheduleItem
  for i := range f.CronStrings {
    if f.CronStrings[i] != "" && f.FuncNames[i] != "" {
      arr := strings.Split(f.FuncNames[i], ":")
      moduleName, funcName := arr[0], arr[1]
      items = append(items, scheduleItem{CronString: f.CronStrings[i], Module: moduleName, FuncName: funcName})
    }
  }
  db.saveSchedules(items)
  app.scheduler.Reload()
  http.Redirect(w, r, "/scheduler", 302)
}
