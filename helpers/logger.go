package helpers

import (
	"os"

	"github.com/kpango/glg"
)

func Logger() *glg.Glg {
	logfile := glg.FileWriter("./workspace/output.log", os.FileMode(0666))
	return glg.New().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.INFO, logfile).
		AddLevelWriter(glg.WARN, logfile).
		AddLevelWriter(glg.FATAL, logfile).
		AddLevelWriter(glg.ERR, logfile).
		AddLevelWriter(glg.FATAL, logfile)
}
