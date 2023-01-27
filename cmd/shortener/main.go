package main

import (
	"github.com/mazaevdmitry/go-yandex-shortener/internal/app"
	"log"
)

func main() {
	server := app.Server()
	log.Println("Starting server...")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
