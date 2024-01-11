package log

import (
	"git.garena.com/shopee/common/ulog"
	"git.garena.com/shopee/platform/splog"
	"git.garena.com/shopee/platform/splog/core"
)

// LogWrapper provides necessary methods to implement the ulog.BaseLogger interface
type LogWrapper struct {
	ulog.BaseLogger
	splogger splog.Interface
}

func (w *LogWrapper) Log(level ulog.Level, msg string, fields ...ulog.TypedField) {
	wrapper := w.With(fields...).(*LogWrapper)
	lvl, _ := toSplogLevel(level)
	wrapper.splogger.Record(lvl, msg)
}

func (w *LogWrapper) With(fields ...ulog.TypedField) ulog.BaseLogger {
	var typedFields splog.TypedFields
	splogger := w.splogger
	for i := range fields {
		typedField := toSplogTypedField(fields[i])
		// Need to handle error separately since WithError() and
		// WithTypedField() has different logic
		if typedField.Type == core.ErrorType {
			err, ok := typedField.Interface.(error)
			if ok && err != nil {
				splogger = splogger.WithError(err)
			}
		} else {
			typedFields = append(typedFields, typedField)
		}
	}
	splogger = splogger.WithTypedFields(typedFields)
	return Wrap(splogger)
}

func Wrap(logger splog.Interface) *LogWrapper {
	return &LogWrapper{
		splogger: logger,
	}
}
