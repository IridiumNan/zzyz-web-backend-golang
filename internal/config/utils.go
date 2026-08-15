package config

import (
	"os"
	"path"
)

// checkIfConfigFileExist : return true if exist config file
func checkIfConfigFileExist(dirPath string) (string, bool) {
	configFilePath := path.Join(dirPath, configFileName)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return configFilePath, false
	}

	return configFilePath, true
}
