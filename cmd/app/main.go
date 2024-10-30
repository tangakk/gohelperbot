package main

import (
	"gohelperbot/internal/tgbot"
	"gohelperbot/pkg/questions"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/ru"
)

func main() {
	l, _ := golem.New(ru.New())
	q := questions.New(l)
	q.ReadFromJSON("q.json")
	config := tgbot.Config{}
	config.ReadFromJSON("config.json")
	bot := tgbot.New(config)
	bot.Q = *q
	bot.Run()
}
