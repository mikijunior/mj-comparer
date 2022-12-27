package diffimage

import (
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mjcomparer/internal/app/response"
	"mjcomparer/internal/app/store"
	"net/http"

	"github.com/vitali-fedulov/images4"
	_ "golang.org/x/image/webp"
)

func CompareImagesByUrls(src, dst string) (response.CompareResponse, error) {
	var result response.CompareResponse
	srcExt, err := GetImageExt(src)
	if err != nil {
		return result, err
	}

	dstExt, err := GetImageExt(dst)
	if err != nil {
		return result, err
	}

	if srcExt == dstExt {
		percent, err := CompareImageUrls(src, dst)
		if err != nil {
			return result, err
		}

		result.Is_same = percent < 10.00
		result.Percent = percent
		result.Algorithm = "Pixel compare"
	} else {
		isSimilar, err := CompareImagesED(src, dst)
		if err != nil {
			return result, err
		}
		result.Is_same = isSimilar
		result.Percent = 0.00
		result.Algorithm = "Euclidean distance"
	}

	return result, nil
}

func GetImageExt(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	buff := make([]byte, 512)
	_, err = response.Body.Read(buff)

	if err != nil {
		return "", errors.New("Can't load a content of url")
	}

	return http.DetectContentType(buff), nil
}

func CompareImageUrls(src, dst string) (percent float64, err error) {
	first, err := GetImage(src)
	if err != nil {
		return 0, err
	}
	second, err := GetImage(dst)
	if err != nil {
		return 0, err
	}
	return CompareImagesByPixels(first, second)
}

func GetImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Can't load a content of url: " + url)
	}

	img, _, err := image.Decode(response.Body)

	return img, err
}

func CompareImagesByPixels(src, dst image.Image) (percent float64, err error) {
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()
	if !boundsMatch(srcBounds, dstBounds) {
		return 0.0, store.ErrImagesWithDifferentSize
	}

	diffImage := image.NewRGBA(image.Rect(0, 0, srcBounds.Max.X, srcBounds.Max.Y))

	var differentPixels float64
	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()
			diffImage.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 64})

			if !isEqualColor(src.At(x, y), dst.At(x, y)) {
				differentPixels++
			}
		}
	}

	return (differentPixels / float64(srcBounds.Max.X*srcBounds.Max.Y) * 100), nil
}

func isEqualColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()

	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func boundsMatch(a, b image.Rectangle) bool {
	return a.Min.X == b.Min.X && a.Min.Y == b.Min.Y && a.Max.X == b.Max.X && a.Max.Y == b.Max.Y
}

// Euclidean distance algorithm. Read more: https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html
func CompareImagesED(src, dst string) (bool, error) {
	img1, err := GetImage(src)
	if err != nil {
		return false, err
	}
	img2, err := GetImage(dst)
	if err != nil {
		return false, err
	}

	// Icons are compact image representations (image "hashes").
	icon1 := images4.Icon(img1)
	icon2 := images4.Icon(img2)

	return images4.Similar(icon1, icon2), nil
}
