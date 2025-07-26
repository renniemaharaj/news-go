package loggers

import "github.com/renniemaharaj/grouplogs/pkg/logger"

var GROUP_PUBLIC *logger.Group

// var GROUP_PRIVATE *logger.Group

var LOGGER_CFX *logger.Logger
var LOGGER_SOCKET *logger.Logger
var LOGGER_STORE *logger.Logger
var LOGGER_SERVER *logger.Logger
var LOGGER_BROWSER *logger.Logger
var LOGGER_COORDINATOR *logger.Logger
var LOGGER_TRANSFORMER *logger.Logger

func Initialize() {
	LOGGER_CFX = logger.New().Prefix("Config")
	LOGGER_STORE = logger.New().Prefix("Store")
	LOGGER_SOCKET = logger.New().Prefix("Socket")
	LOGGER_SERVER = logger.New().Prefix("Server")
	LOGGER_BROWSER = logger.New().Prefix("Browser")

	LOGGER_TRANSFORMER = logger.New().Prefix("Transformer")
	LOGGER_COORDINATOR = logger.New().Prefix("Coordinator")

	GROUP_PUBLIC = logger.NewGroup()
	// GROUP_PRIVATE = logger.NewGroup()

	GROUP_PUBLIC.Join(LOGGER_COORDINATOR)
	GROUP_PUBLIC.Join(LOGGER_TRANSFORMER)
	GROUP_PUBLIC.Join(LOGGER_BROWSER)
	GROUP_PUBLIC.Join(LOGGER_SERVER)
	GROUP_PUBLIC.Join(LOGGER_SOCKET)
}
