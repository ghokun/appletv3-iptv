package main

import (
	"flag"

	"github.com/ghokun/appletv3-iptv/internal/config"
	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
	"github.com/ghokun/appletv3-iptv/internal/server"
)

func main() {

	configFilePtr := flag.String("config", "config.yaml", "Config file path")
	config.LoadConfig(*configFilePtr)
	flag.Parse()

	if config.Current.LogToFile {
		logging.EnableLoggingToFile()
	}

	logging.Info("Starting appletv3-iptv")

	if config.Current.M3UPath != "" {
		err := m3u.GeneratePlaylist()
		if err != nil {
			logging.Warn(err)
		}
	}

	server.Serve()
}
