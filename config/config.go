package config

import (
	"os"
	"strconv"

	lightning "github.com/janritter/go-lightning-log"
)

// Logger contains the Lightning Logger instance configured by the Init function, it's used for logging by calling the Log function on it.
var Logger *lightning.Lightning

// Init is used to initalize Lightning Logger with the configured LogLevel.
func Init() {
	logLevel, _ := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	Logger, _ = lightning.Init(logLevel)
}
