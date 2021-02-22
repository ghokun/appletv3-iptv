package appletv

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"github.com/ghokun/appletv3-iptv/internal/config"
	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
	"golang.org/x/text/language"
)

const (
	basePath = "https://appletv.redbull.tv"
	baseXML  = "templates/base.xml"
	errorXML = "templates/error.xml"
)

var matcher = language.NewMatcher([]language.Tag{
	language.AmericanEnglish,
	language.Turkish,
})

// TemplateData struct is evaluated in all pages.
type TemplateData struct {
	BasePath     string
	BodyID       string
	Data         interface{}
	Translations map[string]string
}

// ErrorData struct is evaluated error pages.
type ErrorData struct {
	Title       string
	Description string
}

// SettingsData struct is evaluated in Setting pages.
type SettingsData struct {
	Version              string
	M3UPath              string
	ReloadChannelsActive bool
	ChannelCount         int
	RecentCount          int
	FavoritesCount       int
	LogsActive           bool
}

// GenerateXML : Parses base XML with given template
func GenerateXML(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	template, err := template.ParseFiles(baseXML, templateName)
	if err != nil {
		logging.Warn(err)
		GenerateErrorXML(w, r, ErrorData{
			Title:       "Template Parse Error",
			Description: err.Error(),
		})
		return
	}
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, accept)
	file, err := os.Open("templates/locales/" + tag.String() + ".json")
	if err != nil {
		logging.Warn(err)
		GenerateErrorXML(w, r, ErrorData{
			Title:       "Translation Load Error",
			Description: err.Error(),
		})
		return
	}
	defer file.Close()
	var translations map[string]string
	if err := json.NewDecoder(file).Decode(&translations); err != nil {
		logging.Warn(err)
		GenerateErrorXML(w, r, ErrorData{
			Title:       "Translation Decode Error",
			Description: err.Error(),
		})
		return
	}
	templateData := TemplateData{
		BasePath:     basePath,
		BodyID:       templateName,
		Data:         data,
		Translations: translations,
	}
	w.Header().Set("Content-Type", "application/xml")
	err = template.Execute(w, templateData)
	if err != nil {
		logging.Warn(err)
		GenerateErrorXML(w, r, ErrorData{
			Title:       "Template Execution Error",
			Description: err.Error(),
		})
		return
	}
}

// GenerateErrorXML : If this fails application should stop.
func GenerateErrorXML(w http.ResponseWriter, r *http.Request, errorData ErrorData) {
	template, err := template.ParseFiles(errorXML)
	if err != nil {
		logging.Fatal(err)
	}
	templateData := TemplateData{
		BasePath:     basePath,
		BodyID:       errorXML,
		Data:         errorData,
		Translations: nil,
	}
	w.Header().Set("Content-Type", "application/xml")
	err = template.Execute(w, templateData)
	if err != nil {
		logging.Fatal(err)
	}
}

// GetSettingsData provides data to Settings page.
func GetSettingsData() SettingsData {
	return SettingsData{
		Version:              config.Version,
		M3UPath:              config.Current.M3UPath,
		ReloadChannelsActive: config.Current.M3UPath != "",
		ChannelCount:         m3u.GetPlaylist().GetChannelsCount(),
		RecentCount:          m3u.GetPlaylist().GetRecentChannelsCount(),
		FavoritesCount:       m3u.GetPlaylist().GetFavoriteChannelsCount(),
		LogsActive:           config.Current.LogToFile,
	}
}
