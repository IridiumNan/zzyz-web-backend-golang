package config

type WebConfig struct {
	// NOTE: MainAddress is the address which is used to reverse proxy
	// You should make it expose to the Intenet
	MainAddress string `toml:"main_address"`

	// WARN: You should use this address just on lo net or tailnet
	// Don't expose to Intenet
	InternalAddress string `toml:"internal_address"`

	// database file path, load from config.toml
	DatabasePath string `toml:"database_path"`

	// log file path, load from cofnig.toml
	LogPath string `toml:"log_path"`
}
