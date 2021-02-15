package appletv

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/ghokun/appletv3-iptv/internal/m3u"
	"golang.org/x/text/language"
)

type TemplateData struct {
	BasePath     string
	BodyID       string
	Data         interface{}
	Translations map[string]string
}

var matcher = language.NewMatcher([]language.Tag{
	language.AmericanEnglish,
	language.Turkish,
})

func generateXML(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	template, err := template.ParseFiles("templates/base.xml", templateName)
	if err != nil {
		panic(err)
	}
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, accept)
	file, err := os.Open("templates/locales/" + tag.String() + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var translations map[string]string
	if err := json.NewDecoder(file).Decode(&translations); err != nil {
		log.Fatal(err)
	}

	templateData := TemplateData{
		BasePath:     "https://appletv.redbull.tv",
		BodyID:       templateName,
		Data:         data,
		Translations: translations,
	}
	w.Header().Set("Content-Type", "application/xml")
	err = template.Execute(w, templateData)
	if err != nil {
		panic(err)
	}
}

// MainHandler handles requests to the main navigation page.
// https://appletv.redbull.tv
func MainHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/main.xml", m3u.GetPlaylist())
}

// CategoryHandler handles requests to the navigated tab on main page.
// https://appletv.redbull.tv/category.xml?category=
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	generateXML(w, r, "templates/category.xml", m3u.GetPlaylist().Categories[category])
}

// RecentHandler handles requests to the navigated tab on main page.
// https://appletv.redbull.tv/recent.xml
func RecentHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/recent.xml", m3u.GetPlaylist())
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	// for categoryKey, categoryValue := range m3u.GetPlaylist().Categories {
	// 	summary := summary + categoryKey + ": " + len(categoryValue.Channels) + "\n"
	// }
	generateXML(w, r, "templates/settings.xml", "")
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/search.xml", m3u.GetPlaylist())
}

func SearchResultsHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	searchResults := m3u.Playlist{}
	for categoryKey, categoryValue := range m3u.GetPlaylist().Categories {
		for channelKey, channelValue := range categoryValue.Channels {
			if strings.Contains(strings.ToLower(channelValue.Title), strings.ToLower(term)) {
				if searchResults.Categories == nil {
					searchResults.Categories = make(map[string]m3u.Category)
				}
				if _, ok := searchResults.Categories[categoryKey]; !ok {
					searchResults.Categories[categoryKey] = m3u.Category{
						Name:     categoryValue.Name,
						Channels: make(map[string]m3u.Channel),
					}
				}
				if _, ok := searchResults.Categories[categoryKey].Channels[channelKey]; !ok {
					searchResults.Categories[categoryKey].Channels[channelKey] = channelValue
				}
			}
		}
	}
	generateXML(w, r, "templates/search-results.xml", searchResults)
}

// PlayerHandler handles requests to the TV Channel player
// https://appletv.redbull.tv/player.xml?category=..&channel=..
func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	channel := r.URL.Query().Get("channel")
	selectedChannel := m3u.GetPlaylist().Categories[category].Channels[channel]
	generateXML(w, r, "templates/player.xml", selectedChannel)

	m3u.GetPlaylist().PrependRecentChannel(selectedChannel)
}
