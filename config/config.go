package config

import (
	"os"
	"strconv"

	lightning "github.com/janritter/go-lightning-log"
)

var Logger *lightning.Lightning

func Init() {
	logLevel, _ := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	Logger, _ = lightning.Init(logLevel)
}
