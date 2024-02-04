package worker

import (
	"fmt"

	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/rs/zerolog"
)

/*
zerolog adapter
*/
type loggerAdapter struct {
}

func NewLoggerAdapter() *loggerAdapter {
	return &loggerAdapter{}
}

/*
use to call zerolog.log.Withlevel
*/
func Print(level zerolog.Level, args ...interface{}) {
	logger.Logger.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l loggerAdapter) Debug(args ...interface{}) {
	logger.Logger.WithLevel(zerolog.DebugLevel).Msg(fmt.Sprint(args...))
}

// Info logs a message at Info level.
func (l loggerAdapter) Info(args ...interface{}) {
	logger.Logger.WithLevel(zerolog.InfoLevel).Msg(fmt.Sprint(args...))
}

// Warn logs a message at Warning level.
func (l loggerAdapter) Warn(args ...interface{}) {
	logger.Logger.WithLevel(zerolog.WarnLevel).Msg(fmt.Sprint(args...))
}

// Error logs a message at Error level.
func (l loggerAdapter) Error(args ...interface{}) {
	logger.Logger.WithLevel(zerolog.ErrorLevel).Msg(fmt.Sprint(args...))
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l loggerAdapter) Fatal(args ...interface{}) {
	logger.Logger.WithLevel(zerolog.FatalLevel).Msg(fmt.Sprint(args...))
}
