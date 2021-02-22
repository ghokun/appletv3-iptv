package config

import (
	"os"

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
	Favorites   []string `yaml:"favorites"`
}

// Version - Application version.
const Version = "0.1.0"

// Current - Global configuration variable.
var Current *Config
var currentConfigFile *string

// LoadConfig - Loads configuration file.
func LoadConfig(configFile string) (err error) {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Current)
	if err != nil {
		return err
	}
	currentConfigFile = &configFile
	return nil
}

func saveConfig() (err error) {
	file, err := os.OpenFile(*currentConfigFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&Current)
	if err != nil {
		return err
	}
	return nil
}

// SaveM3UPath - Edits M3U path and saves to configuration file.
func SaveM3UPath(newM3UPath string) (err error) {
	Current.M3UPath = newM3UPath
	return saveConfig()
}

// SaveFavorites - Save favorite channels to file, in order to preserve between restarts.
func SaveFavorites(newFavorites []string) (err error) {
	Current.Favorites = newFavorites
	return saveConfig()
}

//
func ClearFavorites() (err error) {
	Current.Favorites = nil
	return saveConfig()
}
