package config

import (
	"log"
	"log/slog"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
)

//go:embed config.toml
var configFileContent string

const (
	configFileName = "config.toml"
)

var defaultConfigDir string = getDefaultConfigDir()

var configDirCandidates = []string{
	".",
	defaultConfigDir,
}

var GlobalWebConfig *WebConfig

type WebConfig struct {
	Address      string `toml:"serve_address"`
	DatabasePath string `toml:"database_path"`
	LogPath      string `toml:"log_path"`
}

func InitGlobalConfig() {
	var configFilePath string
	var exist bool

	for _, dir := range configDirCandidates {

		configFilePath, exist = checkIfConfigFileExist(dir)

		if exist {
			break
		}
		exist = false
	}

	if !exist {
		err := initDefaultConfig()
		if err != nil {
			log.Fatal("error when init default config: ", err)
		}
	}

	slog.Info("use config file", "file_path", configFilePath)

	tomlData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("error when load config: ", err)
	}

	var webConfig WebConfig
	err = toml.Unmarshal(tomlData, &webConfig)
	if err != nil {
		log.Fatal("error when unmarshal toml config: ", err)
	}

	GlobalWebConfig = &webConfig
}

func GetConfig() *WebConfig {
	return GlobalWebConfig
}

// initDefaultConfig : place the config.toml file to defaultConfigPath
func initDefaultConfig() (err error) {
	configDir := getDefaultConfigDir()
	err = utils.EnsureDir(configDir)
	if err != nil {
		return
	}

	defaultConfigPath := path.Join(configDir, configFileName)

	var configFile *os.File
	configFile, err = os.OpenFile(defaultConfigPath, os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		return
	}

	_, err = configFile.WriteString(configFileContent)
	if err != nil {
		return
	}

	slog.Info("use a default config file", "config_path", defaultConfigPath)

	return
}

// getDefualtConfitDir : use ~/.config/zzyz-web/ as the default config path
// if fail, use pwd
func getDefaultConfigDir() (configPath string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {

		pwd, pwdErr := os.Getwd()
		if pwdErr != nil {
			slog.Error("error when get current dir path", "error", pwdErr)
			slog.Warn("use . as the config dir path")

			return "."
		}

		configPath = pwd

		slog.Error("error when get home dir path, use current path", "error", err, "config_dir_path", configPath)

		return
	}

	configPath = path.Join(homeDir, ".config", "zzyz-web")

	return
}

func GetLogFilePath() string {
	rawLogPath := GetConfig().LogPath

	homeDir, _ := os.UserHomeDir()

	return strings.Replace(rawLogPath, "~", homeDir, 1)
}
