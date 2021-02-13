package main

import (
	"log"
	"net/http"

	"github.com/ghokun/appletv3-iptv/internal/appletv"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
)

func serveHTTP(mux *http.ServeMux, errs chan<- error) {
	errs <- http.ListenAndServe(":80", mux)
}

func serveHTTPS(mux *http.ServeMux, errs chan<- error) {
	errs <- http.ListenAndServeTLS(":443", "assets/certs/redbulltv.pem", "assets/certs/redbulltv.key", mux)
}

func main() {
	playlist, _ := m3u.Parse("news.m3u")
	// if err == nil {
	// 	for _, category := range playlist.Categories {
	// 		for _, channel := range category.Channels {
	// 			fmt.Println("ID: ", channel.ID)
	// 			fmt.Println("Title: ", channel.Title)
	// 			fmt.Println("Media URL: ", channel.MediaURL)
	// 			fmt.Println("Logo: ", channel.Logo)
	// 			fmt.Println("Description: ", channel.Description)
	// 		}
	// 	}
	// }

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/js"))))
	mux.HandleFunc("/redbulltv.cer", appletv.CertificateHandler)

	mux.HandleFunc("/category.xml", func(w http.ResponseWriter, r *http.Request) {
		appletv.CategoryHandler(w, r, playlist)
	})
	mux.HandleFunc("/player.xml", func(w http.ResponseWriter, r *http.Request) {
		appletv.PlayerHandler(w, r, playlist)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		appletv.MainHandler(w, r, playlist)
	})

	errs := make(chan error, 1) // a channel for errors
	go serveHTTP(mux, errs)     // start the http server in a thread
	go serveHTTPS(mux, errs)    // start the https server in a thread
	log.Fatal(<-errs)           // block until one of the servers writes an error
}
