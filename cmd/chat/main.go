package main

import (
	"chatroom/internal/apps/chat"
	"net/http"
)

func main() {
	app := chat.NewApp()
	err := http.ListenAndServe(":8080", app.Router)
	if err != nil {
		panic(err)
	}
}
