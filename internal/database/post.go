package database

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
)

// LoadPostConfigFromToml : load the post configuration from the toml file
// if fail, return nil
func LoadPostConfigFromToml(configPath string) (*models.PostConfig, error) {
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error when open file, file path: %s, err: %w", configPath, err)
	}

	defer file.Close()

	tomlByte, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error when read file, file path: %s, err: %w", configPath, err)
	}

	var postConfig models.PostConfig

	err = toml.Unmarshal(tomlByte, &postConfig)
	if err != nil {
		return nil, fmt.Errorf("error when Unmarshal the config file content, raw content: %s, err: %w", string(tomlByte), err)
	}

	return &postConfig, nil
}
