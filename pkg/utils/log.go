package utils

import (
	log "github.com/sirupsen/logrus"
)

const (
	logLevelEnvVariable = "LOG_LEVEL"
)

func InitLogging() {
	logLevel := GetEnvWithDefaults(logLevelEnvVariable, log.DebugLevel.String())
	level, levelerr := log.ParseLevel(logLevel)
	if levelerr != nil {
		log.Infof("No valid log level from ENV, defaulting to debug level was: %s", level)
		level = log.DebugLevel
	}
	log.Infof("2 Using log level  %s / %d", logLevel, level)
	log.SetLevel(level)
}
