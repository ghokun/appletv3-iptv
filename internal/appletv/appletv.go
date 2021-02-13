package appletv

import (
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/ghokun/appletv3-iptv/internal/m3u"
)

func CertificateHandler(w http.ResponseWriter, r *http.Request) {

	f, err := os.Open("assets/certs/redbulltv.cer")
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
	f.Close()
}

func MainHandler(w http.ResponseWriter, r *http.Request, playlist m3u.Playlist) {
	t, err := template.ParseFiles("templates/main.xml")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	err = t.Execute(w, playlist)
	if err != nil {
		panic(err)
	}
}

func CategoryHandler(w http.ResponseWriter, r *http.Request, playlist m3u.Playlist) {
	t, err := template.ParseFiles("templates/category.xml")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	category := r.URL.Query().Get("category")
	err = t.Execute(w, playlist.Categories[category])
	if err != nil {
		panic(err)
	}
}

func PlayerHandler(w http.ResponseWriter, r *http.Request, playlist m3u.Playlist) {
	t, err := template.ParseFiles("templates/player.xml")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	category := r.URL.Query().Get("category")
	channel := r.URL.Query().Get("channel")
	err = t.Execute(w, playlist.Categories[category].Channels[channel])
	if err != nil {
		panic(err)
	}
}
