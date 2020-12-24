package server

import "os"

// GetBindAddress gets the value for the PORT environment variable, if it exists - if not, defaults to port 8080
func GetBindAddress() string {
	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8080"
	}
	return ":" + port
}
