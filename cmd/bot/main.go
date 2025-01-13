package main

import (
	"chatroom/internal/apps/bot"
)

func main() {
	_ = bot.NewApp()

	//err := http.ListenAndServe(":8080", app.Router)
	//if err != nil {
	//	panic(err)
	//}
}
