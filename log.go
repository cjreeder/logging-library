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

	err := SetLevel(level)
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

	Logger.log = log.Sugar()
	Logger.log.Info("Logging Service has Started....")

	return Logger{
		log: log.Sugar(),
	}

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
// Do we want to catch an error here and return it???????
// If not, let's not add an error return
func (L *Logger) GetLevel() (string, error) {

	getlvl, err := atom.Level().String()
	if err != nil {
		fmt.Errorf("Error Getting Logging Level: %v", err.Error())
	}

	return getlvl, nil
}

// Handlers for settings and getting the logs
func (L *Logger) SetLogLevel(g *gin.Context) {
	level := g.Param("level")

	L.log.Infof("Setting log level to %s", level)
	err := SetLevel(level)
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}

	L.log.Infof("Log level set to %s", level)
	g.JSON(http.StatusOK, "ok")
}

func (L *Logger) GetLogLevel(g *gin.Context) {
	L.Infof("Getting log level.....")
	level, err := GetLevel()
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}

	L.Log.Infof("Log level is %s", level)

	m := make(map[string]string)
	m["log-level"] = level

	g.JSON(http.StatusOK, m)
}
