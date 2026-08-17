package database

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type PostConfig struct {
	Title    string `toml:"title"`
	Author   string `toml:"author"`
	Overview string `toml:"overview"`
}

type Post struct {
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Overview        string    `json:"overview"`
	MarkdownContent string    `json:"markdown_content"`
	IsPending       bool      `json:"is_pending"`
	ViewCount       int       `json:"view_count"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}

// LoadPostConfigFromToml : load the post configuration from the toml file
// if fail, return nil
func LoadPostConfigFromToml(configPath string) *PostConfig {
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil
	}

	defer file.Close()

	var tomlByte []byte
	_, err = file.Read(tomlByte)
	if err != nil {
		return nil
	}

	var postConfig PostConfig

	err = toml.Unmarshal(tomlByte, &postConfig)
	if err != nil {
		return nil
	}

	return &postConfig
}
