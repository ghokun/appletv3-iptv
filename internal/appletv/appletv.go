package appletv

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/ghokun/appletv3-iptv/internal/config"
	"github.com/ghokun/appletv3-iptv/internal/logging"
	"github.com/ghokun/appletv3-iptv/internal/m3u"
)

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	logging.Warn("Error at " + r.RequestURI + ". With details: " + err.Error())
	GenerateErrorXML(w, r, ErrorData{
		Title:       "Error",
		Description: err.Error(),
	})
}

func unsupportedOperationHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, errors.New("Unsupported operation"))
}

// MainHandler https://appletv.redbull.tv
func MainHandler(w http.ResponseWriter, r *http.Request) {
	logging.CheckLogRotationAndRotate()
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/main.xml", m3u.GetPlaylist().GetChannelsCount())
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ChannelsHandler https://appletv.redbull.tv/channels.xml
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/channels.xml", m3u.GetPlaylist())
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ChannelOptionsHandler https://appletv.redbull.tv/channel-options.xml?category=..&channel=..
func ChannelOptionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		category := r.URL.Query().Get("category")
		channel := r.URL.Query().Get("channel")
		value, err := m3u.GetPlaylist().GetChannel(category, channel)
		if err != nil {
			errorHandler(w, r, err)
		} else {
			GenerateXML(w, r, "templates/channel-options.xml", value)
		}
	default:
		unsupportedOperationHandler(w, r)
	}
}

// RecentHandler https://appletv.redbull.tv/recent.xml
func RecentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/recent.xml", m3u.GetPlaylist())
	default:
		unsupportedOperationHandler(w, r)
	}
}

// FavoritesHandler https://appletv.redbull.tv/favorites.xml
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/favorites.xml", m3u.GetPlaylist())
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ToggleFavoriteHandler https://appletv.redbull.tv/toggle-favorite.xml?category=..&channel=..
func ToggleFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		category := r.URL.Query().Get("category")
		channel := r.URL.Query().Get("channel")
		err := m3u.GetPlaylist().ToggleFavoriteChannel(category, channel)
		if err != nil {
			errorHandler(w, r, err)
		}
	default:
		unsupportedOperationHandler(w, r)
	}
}

// CategoryHandler https://appletv.redbull.tv/category.xml?category=..
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		category := r.URL.Query().Get("category")
		value, err := m3u.GetPlaylist().GetCategory(category)
		if err != nil {
			errorHandler(w, r, err)
		} else {
			GenerateXML(w, r, "templates/category.xml", value)
		}
	default:
		unsupportedOperationHandler(w, r)
	}
}

// PlayerHandler https://appletv.redbull.tv/player.xml?category=..&channel=..
func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		category := r.URL.Query().Get("category")
		channel := r.URL.Query().Get("channel")
		selectedChannel, err := m3u.GetPlaylist().GetChannel(category, channel)
		if err != nil {
			errorHandler(w, r, err)
		} else {
			m3u.GetPlaylist().SetRecentChannel(selectedChannel)
			GenerateXML(w, r, "templates/player.xml", selectedChannel)
		}
	default:
		unsupportedOperationHandler(w, r)
	}
}

// SearchHandler https://appletv.redbull.tv/search.xml
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/search.xml", nil)
	default:
		unsupportedOperationHandler(w, r)
	}
}

// SearchResultsHandler https://appletv.redbull.tv/search-results.xml?term=..
func SearchResultsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		term := r.URL.Query().Get("term")
		GenerateXML(w, r, "templates/search-results.xml", m3u.GetPlaylist().SearchChannels(term))
	default:
		unsupportedOperationHandler(w, r)
	}
}

// SettingsHandler https://appletv.redbull.tv/settings.xml
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/settings.xml", GetSettingsData())
	default:
		unsupportedOperationHandler(w, r)
	}
}

// SetM3UHandler https://appletv.redbull.tv/set-m3u.xml
func SetM3UHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		newM3UPath := r.URL.Query().Get("m3u")
		err := config.Current.SaveM3UPath(newM3UPath)
		if err != nil {
			logging.Warn("Error while setting M3U address: " + err.Error())
		} else {
			logging.Info("Setting M3U address to: " + newM3UPath)
		}
		http.Redirect(w, r, "/settings.xml", http.StatusSeeOther)
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ReloadChannelsHandler https://appletv.redbull.tv/reload-channels.xml
func ReloadChannelsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GenerateXML(w, r, "templates/reload-channels.xml", nil)
	case "POST":
		logging.Info("Reloading channels...")
		logging.Info("Previous channel count is: " + strconv.Itoa(m3u.GetPlaylist().GetChannelsCount()))
		logging.Info("Previous recent channel count is: " + strconv.Itoa(m3u.GetPlaylist().GetRecentChannelsCount()))
		logging.Info("Previous favorite count is: " + strconv.Itoa(m3u.GetPlaylist().GetFavoriteChannelsCount()))
		recentChannels := m3u.GetPlaylist().GetRecentChannels()
		favoriteChannels := m3u.GetPlaylist().GetFavoriteChannels()
		err := m3u.GeneratePlaylist()
		for _, recent := range recentChannels {
			channel, err := m3u.GetPlaylist().GetChannel(recent.CategoryID, recent.ID)
			if err == nil {
				channel.IsRecent = true
				m3u.GetPlaylist().SetRecentChannel(channel)
			}
		}
		for _, favorite := range favoriteChannels {
			_, err := m3u.GetPlaylist().GetChannel(favorite.CategoryID, favorite.ID)
			if err == nil {
				m3u.GetPlaylist().ToggleFavoriteChannel(favorite.CategoryID, favorite.ID)
			}
		}
		if err != nil {
			errorHandler(w, r, err)
		}
		logging.Info("Reloaded channels.")
		logging.Info("Channel count after reload is: " + strconv.Itoa(m3u.GetPlaylist().GetChannelsCount()))
		logging.Info("Recent Channel count after reload is: " + strconv.Itoa(m3u.GetPlaylist().GetRecentChannelsCount()))
		logging.Info("Favorite Channel count after reload is: " + strconv.Itoa(m3u.GetPlaylist().GetFavoriteChannelsCount()))
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ClearRecentHandler https://appletv.redbull.tv/clear-recent.xml
func ClearRecentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		m3u.GetPlaylist().ClearRecentChannels()
		logging.Info("Cleared recently watched channels.")
		http.Redirect(w, r, "/settings.xml", http.StatusSeeOther)
	default:
		unsupportedOperationHandler(w, r)
	}
}

// ClearFavoritesHandler https://appletv.redbull.tv/clear-favorites.xml
func ClearFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		m3u.GetPlaylist().ClearFavoriteChannels()
		logging.Info("Cleared favorite channels.")
		http.Redirect(w, r, "/settings.xml", http.StatusSeeOther)
	default:
		unsupportedOperationHandler(w, r)
	}
}

// LogsHandler https://appletv.redbull.tv/logs.xml
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := time.Now().Format("2006-01-02")
		logFilePath := path.Join(config.Current.LoggingPath, t+".log")
		logs, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			errorHandler(w, r, err)
		} else {
			GenerateXML(w, r, "templates/logs.xml", string(logs))
		}
	default:
		unsupportedOperationHandler(w, r)
	}
}
