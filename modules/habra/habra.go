package habra

import (
  "github.com/yanzay/autohome/modules/generic"
)

type HabraModule struct {
  generic.GenericModule
}

func (h *HabraModule) GetFunctions() []func(map[string]string) {
  return []func(map[string]string){h.Notify}
}

func (h *HabraModule) Notify(settings map[string]string) {
}

func (h *HabraModule) Send(command string) {

}
