package log

const (
	// Disabled level
	Disabled = -1
	// Debug level
	Debug = iota
	// Info level
	Info
	// Warn level
	Warn
	// Error level
	Error
)

const (
	loggerCtxKey contextKey = "logger"
)
