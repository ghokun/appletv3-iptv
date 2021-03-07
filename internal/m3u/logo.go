package m3u

import (
	"errors"
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

const (
	cacheFolder = ".cache/logo"
	missing     = "/assets/images/missing_logo.png"
)

func missingResponse(err error) (string, error) {
	return missing, err
}

func fileExists(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

func computeChannelLogo(id string, rawLogo string) (logo string, err error) {

	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		return missingResponse(err)
	}

	if rawLogo == "" {
		return missingResponse(nil)
	}

	logoFilename := cacheFolder + "/" + id + ".png"
	logo = "/logo/" + id + ".png"

	if fileExists(logoFilename) {
		return logo, nil
	}

	response, err := http.Get(rawLogo)
	if err != nil {
		return missingResponse(err)
	}

	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		image, _, err := image.Decode(response.Body)
		if err != nil {
			return missingResponse(err)
		}
		resizedImage := resize.Resize(320, 180, image, resize.Lanczos3)

		file, err := os.Create(logoFilename)
		if err != nil {
			return missingResponse(err)
		}
		err = jpeg.Encode(file, resizedImage, nil)
		if err != nil {
			return missingResponse(err)
		}
		return logo, nil
	}
	return missingResponse(errors.New("Error while fetching channel logo. Status code: " + response.Status))
}
