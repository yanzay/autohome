package main

import (
  "github.com/robfig/cron"
)

type Scheduler struct {
  funcs   []func()
  cr      *cron.Cron
  modules map[string]Module
  db      storage
}

func NewScheduler(app *Application) *Scheduler {
  s := new(Scheduler)
  s.modules = app.modules
  s.db = app.db
  s.Reload()
  return s
}

func (s *Scheduler) Reload() {
  if s.cr != nil {
    s.cr.Stop()
  }
  s.cr = cron.New()
  s.loadFunctions()
  s.cr.Start()
}

func (s *Scheduler) loadFunctions() {
  items := s.db.getSchedules()
  for _, item := range items {
    s.cr.AddFunc(item.CronString, s.moduleHandler(item.Module, item.FuncName))
  }
  // cr.AddFunc("0 1 14 * * *", s.funcByName("lightOn"))
  // cr.AddFunc("0 43 13 * * *", s.funcByName("lightOff"))
  // cr.AddFunc("00 45 13 * * *", s.funcByName("notify"))
}

func (s *Scheduler) moduleHandler(module string, funcname string) func() {
  return func() { s.modules[module].Send(funcname) }
}
