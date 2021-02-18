package appletv

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
	"golang.org/x/text/language"
)

var version = "0.1.0"

// TemplateData is used in all xml templates under templates folder.
type TemplateData struct {
	BasePath     string
	BodyID       string
	Data         interface{}
	Translations map[string]string
}

// ErrorDialog is shown in error xml.
type ErrorDialog struct {
	Title       string
	Description string
}

type Settings struct {
	Version              string
	M3U                  string
	ReloadChannelsActive bool
	ChannelCount         int
	RecentCount          int
	FavoritesCount       int
	LogsActive           bool
}

var matcher = language.NewMatcher([]language.Tag{
	language.AmericanEnglish,
	language.Turkish,
})

func generateXML(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	template, err := template.ParseFiles("templates/base.xml", templateName)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

func generateErrorXML(w http.ResponseWriter, r *http.Request, errorDialog ErrorDialog) {
	generateXML(w, r, "templates/error.xml", errorDialog)
}

// MainHandler handles requests to the main navigation page.
// https://appletv.redbull.tv
func MainHandler(w http.ResponseWriter, r *http.Request) {
	logging.CheckLogRotationAndRotate()
	generateXML(w, r, "templates/main.xml", m3u.GetPlaylist())
}

// ChannelsHandler handles requests to the channels page.
// https://appletv.redbull.tv/channels.xml
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/channels.xml", m3u.GetPlaylist())
}

func ChannelOptionsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	channel := r.URL.Query().Get("channel")
	selectedChannel := m3u.GetPlaylist().Categories[category].Channels[channel]
	generateXML(w, r, "templates/channel-options.xml", selectedChannel)
}

// RecentHandler handles requests to the recent tab on channels page.
// https://appletv.redbull.tv/recent.xml
func RecentHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/recent.xml", m3u.GetPlaylist())
}

// RecentHandler handles requests to the recent tab on channels page.
// https://appletv.redbull.tv/favorites.xml
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/favorites.xml", m3u.GetPlaylist())
}

func ToggleFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	channel := r.URL.Query().Get("channel")
	selectedChannel := m3u.GetPlaylist().Categories[category].Channels[channel]
	selectedChannel.IsFavorite = !selectedChannel.IsFavorite
	m3u.GetPlaylist().Categories[category].Channels[channel] = selectedChannel
}

// CategoryHandler handles requests to the navigated tab on channels page.
// https://appletv.redbull.tv/category.xml?category=..
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	generateXML(w, r, "templates/category.xml", m3u.GetPlaylist().Categories[category])
}

// PlayerHandler handles requests to the TV Channel player
// https://appletv.redbull.tv/player.xml?category=..&channel=..
func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	channel := r.URL.Query().Get("channel")
	selectedChannel := m3u.GetPlaylist().Categories[category].Channels[channel]
	m3u.GetPlaylist().SetRecentChannel(selectedChannel)
	generateXML(w, r, "templates/player.xml", selectedChannel)
}

// SearchHandler handles request to the Search navigation pane.
// https://appletv.redbull.tv/search.xml
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	generateXML(w, r, "templates/search.xml", m3u.GetPlaylist())
}

// SearchResultsHandler handles search responses.
// https://appletv.redbull.tv/search-results.xml?term=..
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

// SettingsHandler handles request to the Settings navigation pane.
// POST https://appletv.redbull.tv/settings.xml
func SettingsHandler(w http.ResponseWriter, r *http.Request, enableLogs bool) {
	generateXML(w, r, "templates/settings.xml", Settings{
		Version:              version,
		M3U:                  *m3u.M3UFile,
		ReloadChannelsActive: *m3u.M3UFile != "",
		ChannelCount:         m3u.GetPlaylist().GetChannelCount(),
		RecentCount:          m3u.GetPlaylist().GetRecentCount(),
		FavoritesCount:       m3u.GetPlaylist().GetFavoritesCount(),
		LogsActive:           enableLogs,
	})
}

type m3uBody struct {
	NewFile string `json:"newFile"`
}

// SetM3UHandler sets M3U url.
// POST https://appletv.redbull.tv/set-m3u
func SetM3UHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var body m3uBody

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m3u.SetM3UFile(body.NewFile)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body.NewFile))
}

// ReloadChannelsHandler reloads playlist from given M3U list.
// POST https://appletv.redbull.tv/reload-channels
func ReloadChannelsHandler(w http.ResponseWriter, r *http.Request) {
	//m3u.ReloadChannels()
	println(r.URL.Query().Get("id"))
}

// ReloadChannelsHandler reloads playlist from given M3U list.
// POST https://appletv.redbull.tv/reload-channels
func ClearRecentHandler(w http.ResponseWriter, r *http.Request) {
	m3u.GetPlaylist().ClearRecentChannels()
}

// ReloadChannelsHandler reloads playlist from given M3U list.
// POST https://appletv.redbull.tv/reload-channels
func ClearFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	m3u.GetPlaylist().ClearFavoriteChannels()
}

// LogsHandler handles requests to the logs menu in settings page. Only daily log is shown.
// GET https://appletv.redbull.tv/logs.xml
func LogsHandler(w http.ResponseWriter, r *http.Request, loggingPath string) {
	t := time.Now().Format("2006-01-02")
	logFilePath := path.Join(loggingPath, t+".log")
	logs, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		generateErrorXML(w, r, ErrorDialog{Title: "Error", Description: err.Error()})
	}
	generateXML(w, r, "templates/logs.xml", string(logs))
}
