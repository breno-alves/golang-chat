package main

import (
	"chatroom/internal/apps/bot"
)

func main() {
	app := bot.NewApp()
	app.KeepAlive()
}
