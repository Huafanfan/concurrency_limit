package log

import (
	"log"

	"git.garena.com/shopee/platform/splog"
	"git.garena.com/shopee/platform/splog/handlers/multilevel"
)

func New(name string, cfgs ...Config) splog.Interface {
	c := defaultConfig()
	for _, cfg := range cfgs {
		cfg(c)
	}

	handler, err := multilevel.NewLevelsHandler(c)
	if err != nil {
		log.Fatalf("Fail to create logger handler. Error: %v\n", err)
	}

	levelStr := "info"
	level, err := splog.ParseLevel(levelStr)
	if err != nil {
		log.Fatalf("Fail to parse logger level. LevelStr: %s, Error: %v\n", levelStr, err)
	}

	loggerInternal := splog.New([]splog.Handler{handler}, level)

	logger := loggerInternal.WithMetaField(splog.MetaLoggerName, name)
	// need to set this, otherwise the caller location will always point to log/wrapper.go
	source := 1<<16 + 5
	logger = logger.WithSource(source)

	return logger
}

func CopyWithNewName(logger splog.Interface, name string) splog.Interface {
	return logger.WithMetaField(splog.MetaLoggerName, name)
}

func defaultConfig() *multilevel.FileHandlerConfig {
	cfg := &multilevel.FileHandlerConfig{
		Type:   "FileHandler",
		Levels: []string{"debug", "trace", "info", "warn", "error", "fatal", "assert", "data", "access"},
		Sync:   multilevel.LogSyncConfig{SyncWrite: true},
		File:   "/dev/stderr",
		Message: multilevel.LogMessageConfig{
			Format:       "long",
			FieldsFormat: "text",
			MaxBytes:     10000,
			MetaOption:   "source",
		},
	}
	return cfg
}

type Config func(cfg *multilevel.FileHandlerConfig)

func WithFile(filename string) Config {
	return func(cfg *multilevel.FileHandlerConfig) {
		cfg.File = filename
	}
}
