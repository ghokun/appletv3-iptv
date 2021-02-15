package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ghokun/appletv3-iptv/internal/appletv"
	"github.com/ghokun/appletv3-iptv/pkg/m3u"
)

func serveHTTP(mux *http.ServeMux, port string, errs chan<- error) {
	errs <- http.ListenAndServe(":"+port, mux)
}

func serveHTTPS(mux *http.ServeMux, port string, certificate string, key string, errs chan<- error) {
	errs <- http.ListenAndServeTLS(":"+port, certificate, key, mux)
}

func main() {

	m3uPathPtr := flag.String("m3u", "https://iptv-org.github.io/iptv/countries/uk.m3u", "URL that starts with http(s) or a local file path")
	httpPortPtr := flag.String("http", "80", "Port for http requets.")
	httpsPortPtr := flag.String("https", "443", "Port for http requets.")
	certificatePtr := flag.String("crt", "assets/certs/redbulltv.pem", "Certificate path.")
	keyPtr := flag.String("key", "assets/certs/redbulltv.key", "Key path.")

	flag.Parse()

	// Generate playlist singleton
	err := m3u.GeneratePlaylist(*m3uPathPtr)
	if err != nil {
		log.Fatal(err)
	}

	// Serve both http and https
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir(".cache/logo"))))

	// Serve apple tv pages
	mux.HandleFunc("/", appletv.MainHandler)
	mux.HandleFunc("/recent.xml", appletv.RecentHandler)
	mux.HandleFunc("/category.xml", appletv.CategoryHandler)
	mux.HandleFunc("/search.xml", appletv.SearchHandler)
	mux.HandleFunc("/search-results.xml", appletv.SearchResultsHandler)
	mux.HandleFunc("/player.xml", appletv.PlayerHandler)
	mux.HandleFunc("/settings.xml", appletv.SettingsHandler)

	errs := make(chan error, 1)
	go serveHTTP(mux, *httpPortPtr, errs)
	go serveHTTPS(mux, *httpsPortPtr, *certificatePtr, *keyPtr, errs)
	log.Fatal(<-errs)
}
