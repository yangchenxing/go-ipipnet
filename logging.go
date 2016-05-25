package goipipnet

// Logger 日志回调
type Logger func(skip int, level string, format string, params ...interface{})

var (
	logger Logger
)

// SetLogger 设置日志回调
func SetLogger(l Logger) {
	logger = l
}

func logDebug(format string, params ...interface{}) {
	if logger != nil {
		logger(2, "debug", format, params...)
	}
}

func logInfo(format string, params ...interface{}) {
	if logger != nil {
		logger(2, "info", format, params...)
	}
}

func logWarn(format string, params ...interface{}) {
	if logger != nil {
		logger(2, "warn", format, params...)
	}
}

func logError(format string, params ...interface{}) {
	if logger != nil {
		logger(2, "error", format, params...)
	}
}
