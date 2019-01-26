package config

import (
	"log"
	"os"
	"strconv"

	"github.com/janritter/go-lightning-log"
)

// Logger contains the Lightning Logger instance configured by the Init function, it's used for logging by calling the Log function on it.
var Logger *lightning.Lightning

var version string
var commitHash string
var branch string
var buildTime string

// Init is used to initialize Lightning Logger with the configured LogLevel.
func Init() {
	log.Printf("version - %s | branch - %s | commit hash - %s | build time - %s \n", version, branch, commitHash, buildTime)

	logLevel, err := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	if err != nil {
		log.Println("ERROR - Init() - Convert configLevel")
		log.Println(err)
	}
	Logger, err = lightning.Init(logLevel)
	if err != nil {
		log.Println("ERROR - Init() - Init Logger")
		log.Println(err)
	}
}
