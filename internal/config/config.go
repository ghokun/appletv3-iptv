package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config is the struct for configuration.
type Config struct {
	M3UPath     string   `yaml:"m3uPath"`
	HTTPPort    string   `yaml:"httpPort"`
	HTTPSPort   string   `yaml:"httpsPort"`
	CerPath     string   `yaml:"cerPath"`
	PemPath     string   `yaml:"pemPath"`
	KeyPath     string   `yaml:"keyPath"`
	LogToFile   bool     `yaml:"logToFile"`
	LoggingPath string   `yaml:"loggingPath"`
	Recents     []string `yaml:"recents,flow"`
	Favorites   []string `yaml:"favorites,flow"`
}

var (
	// Current - Global configuration variable.
	Current           *Config
	currentConfigFile *string
	// Version - Set by ldflags
	Version string
)

// LoadConfig - Loads configuration file.
func LoadConfig(configFile string) (err error) {
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(contents, &Current)
	if err != nil {
		return err
	}
	currentConfigFile = &configFile
	return nil
}

func saveConfig(config *Config) (err error) {
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(*currentConfigFile, contents, 0644)
	if err != nil {
		return err
	}
	return nil
}

// SaveM3UPath - Edits M3U path and saves to configuration file.
func (config *Config) SaveM3UPath(newM3UPath string) (err error) {
	config.M3UPath = newM3UPath
	return saveConfig(config)
}

// SaveRecents - Save recent channels to file, in order to preserve between restarts.
func (config *Config) SaveRecents(newRecents []string) (err error) {
	config.Recents = newRecents
	return saveConfig(config)
}

// ClearRecents -
func (config *Config) ClearRecents() (err error) {
	config.Recents = make([]string, 0)
	return saveConfig(config)
}

// SaveFavorites - Save favorite channels to file, in order to preserve between restarts.
func (config *Config) SaveFavorites(newFavorites []string) (err error) {
	config.Favorites = newFavorites
	return saveConfig(config)
}

// ClearFavorites -
func (config *Config) ClearFavorites() (err error) {
	config.Favorites = make([]string, 0)
	return saveConfig(config)
}
