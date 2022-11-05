package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

// L logger
var L = logrus.New()

// https://github.com/wechaty/wechaty/issues/2167 兼容 wecahty 社区的 log 等级
var logLevels = map[string]string{
	"silent":  "panic",
	"silly":   "trace",
	"verbose": "trace",
}

func init() {
	level := os.Getenv("WECHATY_LOG")
	if v, ok := logLevels[level]; ok {
		level = v
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	L.SetLevel(logLevel)
	L.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05.000",
		DisableSorting:            false,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		CallerPrettyfier:          nil,
	})
}
