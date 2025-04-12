package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// Interface -.
type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: &logger,
	}
}

// Debug -.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug().Msgf(message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info().Msgf(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn().Msgf(message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.logger.Error().Msgf(msg.Error(), args...)
	case string:
		l.logger.Error().Msgf(message.(string), args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", zerolog.ErrorLevel, message, msg), args...)
	}
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.logger.Fatal().Msgf(msg.Error(), args...)
	case string:
		l.logger.Fatal().Msgf(message.(string), args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", zerolog.FatalLevel, message, msg), args...)
	}

	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
