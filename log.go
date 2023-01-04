package log

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Log *zap.SugaredLogger
}

// Setting up variables for logging
//var L *zap.SugaredLogger

var cfg zap.Config
var atom zap.AtomicLevel

func NewLogger(level string) Logger {
	// Set the initial logging level (Starting up with Info)
	atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	err := StartLevel(level)
	if err != nil {
		fmt.Errorf("Error setting logging level, going with default Info level")
	}

	// Set the Config to use the Development config
	cfg = zap.NewDevelopmentConfig()

	cfg.OutputPaths = append(cfg.OutputPaths)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = atom

	log, err := cfg.Build()
	if err != nil {
		fmt.Errorf("Error Building zap log configuration: %v", err.Error())
		panic(err)
	}

	//Logger.Log = log.Sugar()
	//Logger.Log.Info("Logging Service has Started....")

	return Logger{
		Log: log.Sugar(),
	}

}

// StartLevel allows for setting the log level on the fly
func StartLevel(level string) error {
	switch level {
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
	default:
		return errors.New("Invalid Level")
	}

	return nil

}

// SetLevel allows for setting the log level on the fly
func (L *Logger) SetLevel(level string) error {
	switch level {
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
	default:
		return errors.New("Invalid Level")
	}

	return nil

}

// GetLevel allows for getting the current logging level
func (L *Logger) GetLevel() (string, error) {

	getlvl := atom.Level().String()
	// use strconv to see if string is empty on return
	if len(getlvl) == 0 {
		return "", fmt.Errorf("Error Getting Logging Level, returned empty")
	}

	return getlvl, nil
}

// Handlers for settings and getting the logs
func (L *Logger) SetLogLevel(g *gin.Context) {
	level := g.Param("level")

	fmt.Printf("Setting log level to %s\n", level)
	err := L.SetLevel(level)
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Log level set to %s\n", level)
	g.JSON(http.StatusOK, "ok")
}

func (L *Logger) GetLogLevel(g *gin.Context) {
	fmt.Printf("Getting log level.....\n")
	level, err := L.GetLevel()
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Log level is %s\n", level)

	m := make(map[string]string)
	m["log-level"] = level

	g.JSON(http.StatusOK, m)
}
