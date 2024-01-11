package limit

import (
	"sync/atomic"

	"git.garena.com/shopee/common/ulog"
	"git.garena.com/shopee/platform/splog"
	_ "git.garena.com/shopee/platform/splog/handlers/multilevel"
	"git.garena.com/yifan.zhangyf/concurrency_limit/log"
)

var logger atomic.Value
var splogger atomic.Value

const loggerName = "concurrency_limit"

func init() {
	l := log.New(loggerName)
	splogger.Store(l)
	wrapper := log.Wrap(l)
	unifiedLogger := ulog.NewLogger(wrapper)
	logger.Store(unifiedLogger)
}

// Deprecated: Logger returns the splog logger of sps to support users that haven't
// migrate to unified log, we don't actually use this to log in golang_splib
func GetLogger() splog.Interface {
	return splogger.Load().(splog.Interface)
}

// GetUnifiedLogger returns the global package-level unified logger.
func GetUnifiedLogger() *ulog.Logger {
	return logger.Load().(*ulog.Logger)
}

// SetLogger replace the global package-level logger.
// Note that we turn splog logger into ulog logger for logging
// This is not thread-safe and should be called during the initializing of the
// program.
func SetLogger(l splog.Interface) {
	l = log.CopyWithNewName(l, loggerName)
	splogger.Store(l)
	wrapper := log.Wrap(l)
	unifiedLogger := ulog.NewLogger(wrapper)
	logger.Store(unifiedLogger)
}

// SetUnifiedLogger replace the global package-level unified logger.
//
// This is not thread-safe and should be called during the initializing of the
// program.
func SetUnifiedLogger(l *ulog.Logger) {
	logger.Store(l)
}
