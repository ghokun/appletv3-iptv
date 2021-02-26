package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ghokun/appletv3-iptv/internal/config"
	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
	"github.com/ghokun/appletv3-iptv/internal/server"
)

func main() {

	configFilePtr := flag.String("config", "config.yaml", "Config file path")
	versionPtr := flag.Bool("v", false, "prints current application version")
	flag.Parse()

	if *versionPtr {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	err := config.LoadConfig(*configFilePtr)
	if err != nil {
		// Fail early if config file is not found
		log.Fatal(err)
	}

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
