package core

import (
	"os"

	"github.com/gizo-network/gizo/helpers"
)

//InitializeDataPath creates .gizo folder and block subfolder
func InitializeDataPath() {
	if os.Getenv("ENV") == "dev" {
		os.MkdirAll(BlockPathDev, os.FileMode(0777))
	} else {
		os.MkdirAll(BlockPathProd, os.FileMode(0777))
	}
}

//RemoveDataPath delete's .gizo / .gizo-dev folder
func RemoveDataPath() {
	logger := helpers.Logger()
	if os.Getenv("ENV") == "dev" {
		err := os.RemoveAll(IndexPathDev)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		err := os.RemoveAll(IndexPathProd)
		if err != nil {
			logger.Fatal(err)
		}
	}
}
