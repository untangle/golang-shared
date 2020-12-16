package main

import (
	"../services/logger"
	"../services/overseer"
	"../services/settings"
)

func main() {

	logger.Startup()
	overseer.Startup()
	settings.Startup()

	logger.Info("Testing logger\n")
	overseer.AddCounter("test", 1)
	logger.Info("Testing overseer: %d\n", overseer.GetCounter("test"))
	settingOut, err := settings.GetSettings([]string{"system"})
	if err != nil {
		logger.Err("Failed to get settings: %s\n", err)
	}
	logger.Info("Testing settings: %s\n", settingOut)

	logger.Shutdown()
	overseer.Shutdown()
	settings.Shutdown()

}
