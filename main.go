package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	os.Exit(Run(os.Args[1:]))
}
