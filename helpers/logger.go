package helpers

import (
	"os"
	"path"

	"github.com/kpango/glg"
)

// ILogger interface for logger
type ILogger interface {
	Info(...interface{}) error
	Infof(string, ...interface{}) error
	Log(...interface{}) error
	Logf(string, ...interface{}) error
	Warn(...interface{}) error
	Warnf(string, ...interface{}) error
	Debug(...interface{}) error
	Debugf(string, ...interface{}) error
	Error(...interface{}) error
	Errorf(string, ...interface{}) error
	Success(...interface{}) error
	Successf(string, ...interface{}) error
	Fail(...interface{}) error
	Failf(string, ...interface{}) error
	Print(...interface{}) error
	Println(...interface{}) error
	Printf(string, ...interface{}) error
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

//Logger initializes logger
func Logger() *glg.Glg {
	var logfile *os.File
	var logger *glg.Glg

	if os.Getenv("ENV") == "dev" {
		logfile = glg.FileWriter("./tmp/gizo.log", os.FileMode(0666))
		logger = glg.Get().
			SetMode(glg.STD).
			AddLevelWriter(glg.INFO, logfile).
			AddLevelWriter(glg.WARN, logfile).
			AddLevelWriter(glg.FATAL, logfile).
			AddLevelWriter(glg.ERR, logfile).
			AddLevelWriter(glg.DEBG, logfile)
	} else {
		logfile = glg.FileWriter(path.Join(os.Getenv("HOME"), ".gizo", "gizo.log"), os.FileMode(0666))
		logger = glg.Get().
			SetMode(glg.BOTH).
			AddLevelWriter(glg.INFO, logfile).
			AddLevelWriter(glg.WARN, logfile).
			AddLevelWriter(glg.DEBG, logfile).
			AddLevelWriter(glg.ERR, logfile).
			AddLevelWriter(glg.FATAL, logfile)
	}

	return logger
}
