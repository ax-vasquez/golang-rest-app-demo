package main

import (
	_ "gorm.io/driver/sqlite"

	"codingtest/server"
)

func main() {
	r := server.SetupRouter()

	err := r.Run(server.GetBindAddress())
	if err != nil {
		panic(err)
	}
}
