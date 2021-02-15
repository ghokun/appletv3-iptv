package m3u

import (
	"bufio"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
)

// Playlist struct defines a M3U playlist. M3U playlist starts with #EXTM3U line.
type Playlist struct {
	Categories     map[string]Category
	RecentChannels []Channel
}

// PrependRecentChannel prepends a recently selected channel to the recent channels slice.
func (playlist *Playlist) PrependRecentChannel(channel Channel) {
	recentChannels := append([]Channel{channel}, playlist.RecentChannels...)
	playlist.RecentChannels = recentChannels
}

// Category in a M3U playlist, group-title attribute.
// Default value is Uncategorized.
type Category struct {
	Name     string
	Channels map[string]Channel
}

// Channel is a TV channel in an M3U playlist. Starts with #EXTINF:- prefix.
// #EXTINF:-1 tvg-id="" tvg-name="" tvg-country="" tvg-language="" tvg-logo="" tvg-url="" group-title="",Channel Name
// https://channel.url/stream.m3u8
type Channel struct {
	ID          string // tvg-id or .Title. Spaces are replaced with underscore.
	Title       string // Channel title, after comma
	MediaURL    string // Second line after #EXTINF:-...
	Logo        string // tvg-logo or placeholder. 16x9 aspect ratio
	Description string // TODO will be used for EPG implementation
	Category    string // group-title or Uncategorized
}

var singleton *Playlist

// GeneratePlaylist takes an m3u playlist and creates Playlists.
func GeneratePlaylist(fileName string) (err error) {
	playlist, err := ParseM3U(fileName)
	singleton = &playlist
	return
}

// GetPlaylist returns singleton
func GetPlaylist() *Playlist {
	return singleton
}

func parseAttributes(attributes string, title string) (category string, id string, logo string, description string) {
	tagsRegExp, _ := regexp.Compile("([a-zA-Z0-9-]+?)=\"([^\"]+)\"")
	tagList := tagsRegExp.FindAllString(attributes, -1)
	id = title
	category = "Uncategorized"
	description = ""

	for i := range tagList {
		tagInfo := strings.Split(tagList[i], "=")
		tagKey := tagInfo[0]
		tagValue := strings.Replace(tagInfo[1], "\"", "", -1)
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
	id = strings.ReplaceAll(id, " ", "_")
	id = base64.StdEncoding.EncodeToString([]byte(id))
	// Cache logo
	_ = os.Mkdir(".cache", os.ModePerm)
	_ = os.Mkdir(".cache/logo", os.ModePerm)

	// Get
	if logo != "" {

		response, err := http.Get(logo)
		if err != nil {
			log.Println(err)
		}
		if response != nil && response.Body != nil {
			defer response.Body.Close()
		}
		// Create file
		logo = ".cache/logo/" + id + ".png"
		file, err := os.Create(logo)
		if err != nil {
			log.Fatal(err)
		}
		logo = "https://appletv.redbull.tv/logo/" + id + ".png"
		defer file.Close()

		image, _, err := image.Decode(response.Body)
		if err == nil {
			newImage := resize.Resize(160, 90, image, resize.Lanczos3)
			err = jpeg.Encode(file, newImage, nil)
		}
	}
	return
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
			}

			if playlist.Categories == nil {
				playlist.Categories = make(map[string]Category)
			}
			if _, ok := playlist.Categories[category]; !ok {
				playlist.Categories[category] = Category{
					Name:     category,
					Channels: make(map[string]Channel),
				}
			}
			if _, ok := playlist.Categories[category].Channels[channel.ID]; !ok {
				playlist.Categories[category].Channels[channel.ID] = channel
			}
		}
	}
	return playlist, err
}
