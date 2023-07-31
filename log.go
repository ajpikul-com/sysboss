package main

import (
	"github.com/ajpikul-com/ilog"
	"github.com/ajpikul-com/wsssh/wsconn"
)

var defaultLogger ilog.LoggerInterface

func initLogger() {
	defaultLogger = &ilog.SimpleLogger{}
	defaultLogger.Level(ilog.INFO)
	err := defaultLogger.Init()
	if err != nil {
		panic(err)
	}
	wsconn.SetDefaultLogger(defaultLogger)
}
