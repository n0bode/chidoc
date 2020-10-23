package chidoc

import (
	"image"
	"image/png"
	"net/http"
)

// HandlerImage is a handle to get a png image
type HandlerImage func() image.Image

// ImageFromURL return handler from a URL
// URL must be returns a png
func ImageFromURL(url string) HandlerImage {
	return func() image.Image {
		resp, err := http.Get(url)
		if err != nil {
			return image.Black
		}
		defer resp.Body.Close()

		img, err := png.Decode(resp.Body)
		if err != nil {
			return image.Black
		}
		return img
	}
}
