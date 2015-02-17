package generic

import "net/url"

type GenericModule struct {
}

func (m *GenericModule) Initialize(options map[string]string, db interface{}) {
}

func (m *GenericModule) Functions() []string {
	return []string{}
}

func (m *GenericModule) Settings() []string {
	return []string{}
}

func (m *GenericModule) Send(command string) {
}

func (m *GenericModule) Handle(method, action string, params url.Values) (string, map[string]string) {
	return "", map[string]string{}
}

func (m *GenericModule) Menus() []string {
	return []string{}
}
