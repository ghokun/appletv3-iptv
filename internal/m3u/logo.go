package m3u

import (
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func cacheChannelLogo(id string, logo string) (err error) {

	// check if cache directory exist

	// Cache logos
	_ = os.Mkdir(".cache", os.ModePerm)
	_ = os.Mkdir(".cache/logo", os.ModePerm)

	// if logo text is empty put placeholder

	// if logo text is not empty

	// http get

	// success > create resized file

	// fail > put placeholder

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
		logo = "/logo/" + id + ".png"
		defer file.Close()

		image, _, err := image.Decode(response.Body)
		if err == nil {
			newImage := resize.Resize(160, 90, image, resize.Lanczos3)
			err = jpeg.Encode(file, newImage, nil)
		}
	} else {
		logo = "/assets/images/16x9.png"
	}
	return
}
