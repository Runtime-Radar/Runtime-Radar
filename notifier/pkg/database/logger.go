package database

import (
	"strings"

	"github.com/rs/zerolog"
)

type GORMLogger struct {
	*zerolog.Logger
}

// Printf wraps Printf calls, parses incoming data and changes it for better logs readability. Unfortunately there is no easy way of customizing
// GORM logger because another option of re-implementing `gorm.io/gorm/logger.Interface` is not very convenient due to its awkward architecture.
func (l *GORMLogger) Printf(format string, toPrint ...interface{}) {
	// First item is always "*.go" source file with corresponding format prefix
	toPrint = toPrint[1:]

	// Prefixes are taken from `gorm.io/gorm/logger`
	repl := strings.NewReplacer(
		"%s\n[info] ", "<gorm info> %s",
		"%s\n[warn] ", "<gorm warn> %s",
		"%s\n[error] ", "<gorm error> %s",
		"%s\n[%.3fms] [rows:%v] %s", "SQL [%.3fms] [rows:%v] %s",
		"%s %s\n[%.3fms] [rows:%v] %s", "<%s> SQL [%.3fms] [rows:%v] %s",
	)
	format = repl.Replace(format)

	l.Logger.Printf(format, toPrint...)
}
