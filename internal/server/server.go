package server

import (
	"embed"
	"net/http"

	"github.com/ghokun/appletv3-iptv/internal/appletv"
	"github.com/ghokun/appletv3-iptv/internal/config"
	"github.com/ghokun/appletv3-iptv/internal/logging"
)

//go:embed assets/*
var assets embed.FS

func serveHTTP(mux *http.ServeMux, errs chan<- error) {
	port := ":" + config.Current.HTTPPort
	errs <- http.ListenAndServe(port, mux)
}

func serveHTTPS(mux *http.ServeMux, errs chan<- error) {
	port := ":" + config.Current.HTTPSPort
	errs <- http.ListenAndServeTLS(port, config.Current.PemPath, config.Current.KeyPath, mux)
}

func Serve() {
	// Serve both http and https
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/assets/", http.FileServer(http.FS(assets)))
	mux.Handle("/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir(".cache/logo"))))
	mux.HandleFunc("/redbulltv.cer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.Current.CerPath)
	})

	// Serve apple tv pages and functions
	mux.HandleFunc("/", appletv.MainHandler)

	// Channels
	mux.HandleFunc("/channels.xml", appletv.ChannelsHandler)
	mux.HandleFunc("/channel-options.xml", appletv.ChannelOptionsHandler)
	mux.HandleFunc("/recent.xml", appletv.RecentHandler)
	mux.HandleFunc("/favorites.xml", appletv.FavoritesHandler)
	mux.HandleFunc("/toggle-favorite.xml", appletv.ToggleFavoriteHandler)
	mux.HandleFunc("/category.xml", appletv.CategoryHandler)
	mux.HandleFunc("/player.xml", appletv.PlayerHandler)

	// Search
	mux.HandleFunc("/search.xml", appletv.SearchHandler)
	mux.HandleFunc("/search-results.xml", appletv.SearchResultsHandler)

	// Settings
	mux.HandleFunc("/settings.xml", appletv.SettingsHandler)
	mux.HandleFunc("/set-m3u.xml", appletv.SetM3UHandler)
	mux.HandleFunc("/reload-channels.xml", appletv.ReloadChannelsHandler)
	mux.HandleFunc("/clear-recent.xml", appletv.ClearRecentHandler)
	mux.HandleFunc("/clear-favorites.xml", appletv.ClearFavoritesHandler)
	mux.HandleFunc("/logs.xml", appletv.LogsHandler)

	httpErrs := make(chan error, 1)
	go serveHTTP(mux, httpErrs)
	go serveHTTPS(mux, httpErrs)
	logging.Fatal(<-httpErrs)
}
