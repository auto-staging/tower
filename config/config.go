package config

import lightning "github.com/janritter/go-lightning-log"

var Logger *lightning.Lightning

func Init() {
	Logger, _ = lightning.Init(3)
}
