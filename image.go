package chidoc

import (
	"image"
	"image/png"
	"net/http"

	"golang.org/x/image/draw"
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

// ImageFromURL return handler from a URL
// URL must be returns a png
func ImageFromURLScaled(url string, scale float32) HandlerImage {
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

		rw := int(float32(img.Bounds().Max.X) * scale)
		rh := int(float32(img.Bounds().Max.Y) * scale)

		dst := image.NewRGBA(image.Rect(0, 0, rw, rh))

		draw.ApproxBiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
		return dst
	}
}
