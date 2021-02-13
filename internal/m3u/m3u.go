package m3u

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

// Channel is TV channel in an M3U playlist. Starts with #EXTINF:- prefix.
type Channel struct {
	ID          string
	Title       string
	MediaURL    string
	Logo        string
	Description string
}

// Parse  parses an m3u list.
func Parse(file string) (playlist Playlist, err error) {

	fox := Channel{
		ID:          "fox",
		Title:       "FOX TV",
		MediaURL:    "https://mn-nl.mncdn.com/blutv_foxtv/smil:foxtv_sd.smil/playlist.m3u8",
		Logo:        "http://fox.png",
		Description: "FOX TV ana haber",
	}
	trt := Channel{
		ID:          "trt",
		Title:       "TRT",
		MediaURL:    "http://tv-trt1.live.trt.com.tr/master_720.m3u8",
		Logo:        "http://trt.png",
		Description: "TRT 1",
	}
	newsChannels := make(map[string]Channel, 15)
	newsChannels["fox"] = fox
	newsChannels["trt"] = trt
	news := Category{
		ID:       "news",
		Name:     "News",
		Channels: newsChannels,
	}

	trtBelgesel := Channel{
		ID:          "trtbelgesel",
		Title:       "TRT Belgesel",
		MediaURL:    "http://tv-trtbelgesel.live.trt.com.tr/master_720.m3u8",
		Logo:        "http://rtbelgesel.png",
		Description: "TRT Belgesel",
	}
	documentaryChannels := make(map[string]Channel, 15)
	documentaryChannels["trtbelgesel"] = trtBelgesel
	documentary := Category{
		ID:       "documentary",
		Name:     "Documentary",
		Channels: documentaryChannels,
	}

	categories := make(map[string]Category, 15)
	categories["news"] = news
	categories["documentary"] = documentary

	playlist.Categories = categories
	return playlist, err
}

// #EXTINF:-1 tvg-id="" tvg-name="" tvg-country="INT" tvg-language="English" tvg-logo="https://i.imgur.com/mdNV8tQ.jpg" tvg-url="" group-title="News",Australia Channel (720p)
// https://austchannel-live.akamaized.net/hls/live/2002729/austchannel-news/master.m3u8
