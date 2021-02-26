package m3u

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/ghokun/appletv3-iptv/internal/config"
)

var (
	singleton *Playlist
)

// GeneratePlaylist takes an m3u playlist and creates Playlists.
func GeneratePlaylist() (err error) {
	playlist, err := ParseM3U(config.Current.M3UPath)
	singleton = &playlist
	return
}

// GetPlaylist returns singleton
func GetPlaylist() *Playlist {
	return singleton
}

// ParseM3U parses an m3u list.
// Modified code of https://github.com/jamesnetherton/m3u/blob/master/m3u.go
func ParseM3U(fileNameOrURL string) (playlist Playlist, err error) {

	var f io.ReadCloser
	var data *http.Response
	if strings.HasPrefix(fileNameOrURL, "http://") || strings.HasPrefix(fileNameOrURL, "https://") {
		data, err = http.Get(fileNameOrURL)
		f = data.Body
	} else {
		f, err = os.Open(fileNameOrURL)
	}

	if err != nil {
		err = errors.New("Unable to open playlist file")
		return
	}
	defer f.Close()

	onFirstLine := true
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if onFirstLine && !strings.HasPrefix(line, "#EXTM3U") {
			err = errors.New("Invalid m3u file format. Expected #EXTM3U file header")
			return
		}

		onFirstLine = false

		// Find #EXTINF prefixes
		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.Replace(line, "#EXTINF:", "", -1)
			channelInfo := strings.Split(line, ",")
			if len(channelInfo) < 2 {
				err = errors.New("Invalid m3u file format. Expected EXTINF metadata to contain tvg attributes and channel name")
				return
			}
			attributes := channelInfo[0]
			title := channelInfo[1]
			category, id, logo, description := parseAttributes(attributes, title)
			err = cacheChannelLogo(id, logo)
			categoryID := hex.EncodeToString([]byte(category))
			// Next line is m3u8 url
			scanner.Scan()
			mediaURL := scanner.Text()

			channel := Channel{
				ID:          id,
				Title:       title,
				MediaURL:    mediaURL,
				Logo:        logo,
				Description: description,
				Category:    category,
				CategoryID:  categoryID,
			}

			if playlist.Categories == nil {
				playlist.Categories = make(map[string]Category)
			}
			if _, ok := playlist.Categories[categoryID]; !ok {
				playlist.Categories[categoryID] = Category{
					ID:       categoryID,
					Name:     category,
					Channels: make(map[string]Channel),
				}
			}
			if _, ok := playlist.Categories[categoryID].Channels[channel.ID]; !ok {
				playlist.Categories[categoryID].Channels[channel.ID] = channel
			}
		}
	}
	return playlist, err
}

func parseAttributes(attributes string, title string) (category string, id string, logo string, description string) {
	tagsRegExp, _ := regexp.Compile("([a-zA-Z0-9-]+?)=\"([^\"]+)\"")
	tags := tagsRegExp.FindAllString(attributes, -1)
	id = title
	category = "Uncategorized"
	description = "TODO EPG"

	for i := range tags {
		tagParts := strings.Split(tags[i], "=")
		tagKey := tagParts[0]
		tagValue := strings.Replace(tagParts[1], "\"", "", -1)
		if tagKey == "group-title" {
			category = tagValue
		}
		if tagKey == "tvg-id" {
			id = tagValue
		}
		if tagKey == "tvg-logo" {
			logo = tagValue
		}
		if tagKey == "tvg-url" {
			description = tagValue
		}
	}
	id = hex.EncodeToString([]byte(id))
	return category, id, logo, description
}
