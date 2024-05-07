package logger

import (
	"os"
	"time"

	"twitter-clone/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *config.Config) (*logrus.Logger, error) {
	level, _ := logrus.ParseLevel(cfg.LogLevel)
	return &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		},
		Hooks:    make(logrus.LevelHooks),
		Level:    level,
		ExitFunc: os.Exit,
	}, nil
}
