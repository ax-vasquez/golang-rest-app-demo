package server

import "os"

func GetBindAddress() string {
	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8080"
	}
	return ":" + port
}
