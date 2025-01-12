package main

import (
	"chatroom/internal/apps/chat"
	"net/http"
)

func main() {
	//time.Sleep(5000 * time.Second)
	app := chat.NewApp()
	err := http.ListenAndServe(":8080", app.Router)

	if err != nil {
		panic(err)
	}
}
