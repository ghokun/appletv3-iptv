package m3u

import (
	"errors"
	"strings"
)

// Playlist struct defines a M3U playlist. M3U playlist starts with #EXTM3U line.
type Playlist struct {
	Categories map[string]Category
}

// Category in a M3U playlist, group-title attribute.
type Category struct {
	ID       string
	Name     string
	Channels map[string]Channel
}

// Channel is a TV channel in an M3U playlist. Starts with #EXTINF:- prefix.
// #EXTINF:-1 tvg-id="" tvg-name="" tvg-country="" tvg-language="" tvg-logo="" tvg-url="" group-title="",Channel Name
// https://channel.url/stream.m3u8
type Channel struct {
	ID            string // tvg-id or .Title. Spaces are replaced with underscore.
	Title         string // Channel title, string that comes after comma
	MediaURL      string // Second line after #EXTINF:-...
	Logo          string // tvg-logo or placeholder. 16x9 aspect ratio
	Description   string // Unused for now, will be used for EPG implementation
	Category      string // group-title or Uncategorized if missing
	CategoryID    string // For link generation purposes
	IsRecent      bool   // Is channel recently watched?
	RecentOrdinal int    // Recent watch order
	IsFavorite    bool   // Is channel favorite?
}

// GetCategory - Gets Category and its children in current playlist.
func (playlist *Playlist) GetCategory(category string) (value Category, err error) {
	if value, ok := playlist.Categories[category]; ok {
		return value, nil
	}
	return value, errors.New("Category could not be found")
}

// GetChannel - Gets channel with given category and channel values.
func (playlist *Playlist) GetChannel(category string, channel string) (value Channel, err error) {
	cat, err := playlist.GetCategory(category)
	if err != nil {
		return value, err
	}
	if value, ok := cat.Channels[channel]; ok {
		return value, nil
	}
	return value, errors.New("Channel could not be found")
}

// GetChannelsCount - Gets count of all channels.
func (playlist *Playlist) GetChannelsCount() (count int) {
	count = 0
	if playlist == nil {
		return count
	}
	for _, category := range playlist.Categories {
		count += len(category.Channels)
	}
	return count
}

// GetRecentChannelsCount - Gets count of recently watched channels.
func (playlist *Playlist) GetRecentChannelsCount() (count int) {
	count = 0
	if playlist == nil {
		return count
	}
	for _, category := range playlist.Categories {
		for _, channel := range category.Channels {
			if channel.IsRecent {
				count++
			}
		}
	}
	return count
}

// GetFavoriteChannelsCount - Gets count of favorite channels.
func (playlist *Playlist) GetFavoriteChannelsCount() (count int) {
	count = 0
	if playlist == nil {
		return count
	}
	for _, category := range playlist.Categories {
		for _, channel := range category.Channels {
			if channel.IsFavorite {
				count++
			}
		}
	}
	return count
}

// GetRecentChannels - Gets recent channels and puts them in order.
func (playlist *Playlist) GetRecentChannels() (recentChannels []Channel) {
	recentCount := playlist.GetRecentChannelsCount()
	if recentCount == 0 {
		return nil
	}
	recentChannels = make([]Channel, recentCount)
	for _, category := range playlist.Categories {
		for _, channel := range category.Channels {
			if channel.IsRecent {
				recentChannels[channel.RecentOrdinal-1] = channel
			}
		}
	}
	return recentChannels
}

// SetRecentChannel - Sets selected channel as recent and updates order of other channels.
func (playlist *Playlist) SetRecentChannel(selectedChannel Channel) {
	for _, channel := range playlist.GetRecentChannels() {
		if selectedChannel.IsRecent {
			if channel.ID != selectedChannel.ID && selectedChannel.RecentOrdinal > channel.RecentOrdinal {
				channel.RecentOrdinal++
			}
		} else {
			channel.RecentOrdinal++
		}
		playlist.Categories[channel.CategoryID].Channels[channel.ID] = channel
	}
	selectedChannel.RecentOrdinal = 1
	selectedChannel.IsRecent = true
	playlist.Categories[selectedChannel.CategoryID].Channels[selectedChannel.ID] = selectedChannel
}

// ClearRecentChannels - Clears recent channel list.
func (playlist *Playlist) ClearRecentChannels() {
	for _, channel := range playlist.GetRecentChannels() {
		channel.IsRecent = false
		playlist.Categories[channel.CategoryID].Channels[channel.ID] = channel
	}
}

// GetFavoriteChannels - Gets favorite channels.
func (playlist *Playlist) GetFavoriteChannels() (favoriteChannels []Channel) {
	for _, category := range playlist.Categories {
		for _, channel := range category.Channels {
			if channel.IsFavorite {
				favoriteChannels = append(favoriteChannels, channel)
			}
		}
	}
	return favoriteChannels
}

// ToggleFavoriteChannel - Clears favorite channel list.
func (playlist *Playlist) ToggleFavoriteChannel(category string, channel string) (err error) {
	selectedChannel, err := playlist.GetChannel(category, channel)
	if err != nil {
		return err
	}
	selectedChannel.IsFavorite = !selectedChannel.IsFavorite
	playlist.Categories[category].Channels[channel] = selectedChannel
	return nil
}

// ClearFavoriteChannels - Clears favorite channel list.
func (playlist *Playlist) ClearFavoriteChannels() {
	for _, channel := range playlist.GetFavoriteChannels() {
		channel.IsFavorite = false
		playlist.Categories[channel.CategoryID].Channels[channel.ID] = channel
	}
}

// SearchChannels - Searches channel titles with the given term, case insensitive.
func (playlist *Playlist) SearchChannels(term string) (searchResults Playlist) {
	searchResults = Playlist{}
	for categoryKey, categoryValue := range playlist.Categories {
		for channelKey, channelValue := range categoryValue.Channels {
			if strings.Contains(strings.ToLower(channelValue.Title), strings.ToLower(term)) {
				if searchResults.Categories == nil {
					searchResults.Categories = make(map[string]Category)
				}
				if _, ok := searchResults.Categories[categoryKey]; !ok {
					searchResults.Categories[categoryKey] = Category{
						Name:     categoryValue.Name,
						Channels: make(map[string]Channel),
					}
				}
				if _, ok := searchResults.Categories[categoryKey].Channels[channelKey]; !ok {
					searchResults.Categories[categoryKey].Channels[channelKey] = channelValue
				}
			}
		}
	}
	return searchResults
}
