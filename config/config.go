package config

import (
	"log"
	"os"
	"strconv"

	"github.com/janritter/go-lightning-log"
)

// Logger contains the Lightning Logger instance configured by the Init function, it's used for logging by calling the Log function on it.
var Logger *lightning.Lightning

// Init is used to initialize Lightning Logger with the configured LogLevel.
func Init() {
	logLevel, err := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	if err != nil {
		log.Println("Init() - Convert configLevel")
		log.Println(err)
	}
	Logger, err = lightning.Init(logLevel)
	if err != nil {
		log.Println("Init() - Init Logger")
		log.Println(err)
	}
}
