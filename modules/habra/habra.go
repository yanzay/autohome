package habra

import (
	"fmt"
	"log"
	"strings"

	"github.com/yanzay/autohome/modules/generic"
	"github.com/yanzay/googlespeak"
	"github.com/yanzay/habraparser"
	"github.com/yanzay/reactiverecord"
)

type HabraTopic struct {
	Hab   string
	Title string
	Link  string
}

type HabraModule struct {
	generic.GenericModule
	options map[string]string
	db      *reactiverecord.ReactiveRecord
}

// type ReactiveStorage interface {
//   CreateTable(obj interface{}) error
//   Create(obj interface{}) error
//   All(obj interface{}) jet.Runnable
// }

func (h *HabraModule) Initialize(options map[string]string, db interface{}) {
	var err error
	h.options = options
	h.db, err = reactiverecord.Connect("sqlite", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	h.db.CreateTable(HabraTopic{})
}

func (h *HabraModule) Settings() []string {
	return []string{"habs"}
}

func (h *HabraModule) Functions() []string {
	return []string{"notify"}
}

func (h *HabraModule) Notify() {
	topics := h.getTopics()
	for _, topic := range topics {
		h.db.Where("Link = ?", topic.Link)
		text := fmt.Sprintf("%s. %s", topic.Hab, topic.Title)
		googlespeak.Say(text, "ru")
	}
}

func (h *HabraModule) Send(command string) {
	if command == "notify" {
		h.Notify()
	}
}

func (h *HabraModule) getTopics() []HabraTopic {
	habraTopics := make([]HabraTopic, 0)
	habs := strings.Split(h.options["habs"], ",")
	for _, hab := range habs {
		items, _ := habraparser.Read(hab, 3)
		for _, item := range items {
			topic := HabraTopic{Hab: hab, Title: item.Title, Link: item.Link}
			habraTopics = append(habraTopics, topic)
		}
	}
	return habraTopics
}
