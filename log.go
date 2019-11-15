package log

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Default is the package default logger.
	// It can be used right out of the box.
	// It can be replaced by a custom configured one
	// using package Set(*Logger) function
	// or using *Logger.Set() method.
	Default *Logger
)

func init() {
	Default = NewLogger(Debug, "logger")
}

// CtxLogger returns a logger stored in the context provided by the argument.
// We want to avoid runtime errors, then, to get the logger from the context
// the package uses an unexported key to store and retrive it.
func CtxLogger(ctx context.Context) (logger *Logger, ok bool) {
	l, ok := ctx.Value(loggerCtxKey).(*Logger)
	return l, ok
}

// NewLogger logger.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewLogger(level int, name string, stfields ...interface{}) *Logger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	errl := zerolog.New(os.Stderr).With().Timestamp().Logger()

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	l := &Logger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
	}

	if len(stfields) > 1 && !cfg.configured {
		setup(name, stfields)
		Default = l
	}

	return l
}

// NewDevLogger logger.
// Pretty logging for development mode.
// Not recommended for production use.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewDevLogger(level int, name string, stfields ...interface{}) *Logger {
	if level < Disabled || level > Error {
		level = Info
	}

	stdl := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	errl := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	l := &Logger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
	}

	if len(stfields) > 1 && !cfg.configured {
		setup(name, stfields)
		Default = l
	}

	return l
}

// Set the base package logger.
func Set(l *Logger) {
	Default = l
}

// Set default package logger.
// Can be used chained with NewLogger to create a new one,
// set it up as package default logger and get it for use in one step.
// i.e:
// logger := log.NewLogger(log.Debug, "name", "version", "revision").Set()
func (l *Logger) Set() *Logger {
	Default = l
	return Default
}

// InCtx returns a copy of context that also includes a configured logger.
func InCtx(ctx context.Context, fields ...string) context.Context {
	l, _ := FromCtx(ctx)
	if len(fields) > 0 {
		l.SetDyna(fields)
	}
	return context.WithValue(ctx, loggerCtxKey, l)
}

// FromCtx returns current logger in context.
// If there is not logger in context it returns
// a new one with current config values.
func FromCtx(ctx context.Context) (log *Logger, fresh bool) {
	l, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return NewLogger(cfg.level, cfg.name), true
	}
	return &l, false
}

// Debug logs debug messages.
func (l Logger) Debug(meta ...interface{}) {
	if len(meta) > 0 {
		l.debugf(stringify(meta[0]), meta[1:len(meta)])
	}
}

// Info logs info messages.
func (l Logger) Info(meta ...interface{}) {
	if len(meta) > 0 {
		l.infof(stringify(meta[0]), meta[1:len(meta)])
	}
}

// Warn logs warning messages.
func (l Logger) Warn(meta ...interface{}) {
	if len(meta) > 0 {
		l.warnf(stringify(meta[0]), meta[1:len(meta)])
	}
}

// Error logs error messages.
func (l Logger) Error(err error, meta ...interface{}) {
	if len(meta) > 0 {
		l.errorf(err, stringify(meta[0]), meta[1:len(meta)])
		return
	}
	l.errorf(err, "", nil)
}

func (l Logger) debugf(message string, fields []interface{}) {
	if l.Level > Debug {
		return
	}
	le := l.StdLog.Info()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) infof(message string, fields []interface{}) {
	if l.Level > Info {
		return
	}
	le := l.StdLog.Info()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) warnf(message string, fields []interface{}) {
	if l.Level > Warn {
		return
	}
	le := l.StdLog.Info()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) errorf(err error, message string, fields []interface{}) {
	le := l.ErrLog.Error()
	appendKeyValues(le, l.dynafields, fields)
	le.Err(err)
	le.Msg(message)
}

// TODO: Optimize.
// Static key-value calculation shoud be cached.
// Dynamic key-value calculation shoud be cached if didn't changed.
func appendKeyValues(le *zerolog.Event, dynafields []interface{}, fields []interface{}) {
	if cfg.name != "" {
		le.Str("name", cfg.name)
	}

	fs := make(map[string]interface{})

	if len(fields) > 1 {
		for i := 0; i < len(fields)-1; i++ {
			if fields[i] != nil && fields[i+1] != nil {
				k := stringify(fields[i])
				fs[k] = fields[i+1]
				// fmt.Printf("field - (%s, %v)\n", k, fs[k])
				i++
			}
		}

		if len(dynafields) > 1 {
			// fs := make(map[string]interface{})
			for i := 0; i < len(dynafields)-1; i++ {
				if dynafields[i] != nil && dynafields[i+1] != nil {
					k := stringify(dynafields[i])
					fs[k] = dynafields[i+1]
					// fmt.Printf("dyna - (%s, %v)\n", k, fs[k])
					i++
				}
			}
		}

		if len(cfg.stfields) > 1 {
			for i := 0; i < len(cfg.stfields)-1; i++ {
				if cfg.stfields[i] != nil && cfg.stfields[i+1] != nil {
					k := stringify(cfg.stfields[i])
					fs[k] = cfg.stfields[i+1]
					// fmt.Printf("static - (%s, %v)\n", k, fs[k])
					i++
				}
			}
		}

	}
	le.Fields(fs)
}

func stringify(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%+v", v)
	}
}

// UpdateLogLevel updates log level.
func (l Logger) UpdateLogLevel(level int) {
	// Allow info level to log the update
	// But don't downgrade to it if Error is set.
	current := Error
	l.Info("Log level updated", "", "log level", level)
	l.Level = current
	if level < Disabled || level > Error {
		l.Level = level
		setLogLevel(&l.StdLog, level)
		setLogLevel(&l.ErrLog, level)
	}
}
func setLogLevel(l *zerolog.Logger, level int) {
	switch level {
	case -1:
		l.Level(zerolog.Disabled)
	case 0:
		l.Level(zerolog.DebugLevel)
	case 1:
		l.Level(zerolog.InfoLevel)
	case 2:
		l.Level(zerolog.WarnLevel)
	case 3:
		l.Level(zerolog.ErrorLevel)
	default:
		l.Level(zerolog.DebugLevel)
	}
}
