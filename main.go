package main

import (
	"flag"
	"net/http"

	"github.com/ghokun/appletv3-iptv/internal/appletv"
	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
)

func serveHTTP(mux *http.ServeMux, port string, errs chan<- error) {
	errs <- http.ListenAndServe(":"+port, mux)
}

func serveHTTPS(mux *http.ServeMux, port string, certificate string, key string, errs chan<- error) {
	errs <- http.ListenAndServeTLS(":"+port, certificate, key, mux)
}

func main() {

	m3uPathPtr := flag.String("m3u", "", "URL that starts with http(s) or a local file path")
	httpPortPtr := flag.String("http", "80", "Port for http requets.")
	httpsPortPtr := flag.String("https", "443", "Port for https requets.")
	certificatePtr := flag.String("cert", "assets/certs/redbulltv.pem", "Certificate path.")
	keyPtr := flag.String("cert-key", "assets/certs/redbulltv.key", "Key path.")
	logToFilePtr := flag.Bool("log-to-file", true, "Enable/Disable logging to file system.")
	loggingPathPtr := flag.String("log-path", "log", "Logging path")

	flag.Parse()

	if *logToFilePtr {
		logging.EnableLoggingToFile(loggingPathPtr)
	}

	logging.Info("Starting appletv3-iptv")

	if *m3uPathPtr != "" {
		// Generate playlist early if supplied via argument
		err := m3u.GeneratePlaylist(*m3uPathPtr)
		if err != nil {
			logging.Warn(err)
		}
	}

	// Serve both http and https
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir(".cache/logo"))))

	// Serve apple tv pages and functions
	mux.HandleFunc("/", appletv.MainHandler)
	mux.HandleFunc("/channels.xml", appletv.ChannelsHandler)
	mux.HandleFunc("/channel-options.xml", appletv.ChannelOptionsHandler)
	mux.HandleFunc("/recent.xml", appletv.RecentHandler)
	mux.HandleFunc("/favorites.xml", appletv.FavoritesHandler)
	mux.HandleFunc("/toggle-favorite", appletv.ToggleFavoriteHandler)
	mux.HandleFunc("/category.xml", appletv.CategoryHandler)
	mux.HandleFunc("/player.xml", appletv.PlayerHandler)
	mux.HandleFunc("/search.xml", appletv.SearchHandler)
	mux.HandleFunc("/search-results.xml", appletv.SearchResultsHandler)
	mux.HandleFunc("/settings.xml", func(rw http.ResponseWriter, r *http.Request) {
		appletv.SettingsHandler(rw, r, *logToFilePtr)
	})
	mux.HandleFunc("/set-m3u", appletv.SetM3UHandler)
	mux.HandleFunc("/reload-channels", appletv.ReloadChannelsHandler)
	mux.HandleFunc("/clear-recent", appletv.ClearRecentHandler)
	mux.HandleFunc("/clear-favorites", appletv.ClearFavoritesHandler)

	mux.HandleFunc("/logs.xml", func(rw http.ResponseWriter, r *http.Request) {
		appletv.LogsHandler(rw, r, *loggingPathPtr)
	})

	httpErrs := make(chan error, 1)
	go serveHTTP(mux, *httpPortPtr, httpErrs)
	go serveHTTPS(mux, *httpsPortPtr, *certificatePtr, *keyPtr, httpErrs)
	logging.Fatal(<-httpErrs)
}
