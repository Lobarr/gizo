package helpers

import (
	"os"
	"path"

	"github.com/kpango/glg"
)

func Logger() *glg.Glg {
	var logfile *os.File
	if os.Getenv("ENV") == "dev" {
		logfile = glg.FileWriter(path.Join(os.Getenv("HOME"), ".gizo-dev", "output.log"), os.FileMode(0666))
	} else {
		logfile = glg.FileWriter(path.Join(os.Getenv("HOME"), ".gizo", "output.log"), os.FileMode(0666))
	}
	return glg.New().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.INFO, logfile).
		AddLevelWriter(glg.WARN, logfile).
		AddLevelWriter(glg.FATAL, logfile).
		AddLevelWriter(glg.ERR, logfile).
		AddLevelWriter(glg.FATAL, logfile)
}
