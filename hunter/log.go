package hunter

import (
	"log"
	"os"
	"strings"

	"github.com/uber-go/zap"
)

var Logger zap.Logger

func InitLogger(lp string, lv string, isDebug bool) {
	var level zap.Level

	switch strings.ToLower(lv) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}

	if isDebug {
		Logger = zap.New(
			zap.NewJSONEncoder(
				zap.RFC3339Formatter("@timestamp"), // human-readable timestamps
				zap.MessageKey("@message"),         // customize the message key
				zap.LevelString("@level"),          // stringify the log level
			),
			zap.AddCaller(),
			level,
		)
	} else {
		f, err := os.OpenFile(lp, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Panic(err)
		}

		Logger = zap.New(
			zap.NewJSONEncoder(
				zap.RFC3339NanoFormatter("@timestamp"), // human-readable timestamps
				zap.MessageKey("@message"),             // customize the message key
				zap.LevelString("@level"),              // stringify the log level
			),
			zap.Output(f),
			zap.AddCaller(),
			level,
		)
	}
}
