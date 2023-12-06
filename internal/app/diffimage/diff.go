package diffimage

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"mjcomparer/internal/app/response"
	"net/http"

	"golang.org/x/image/tiff"
	_ "golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	_ "golang.org/x/image/webp"
)

// CompareImagesByUrls compares two images given their URLs and returns a comparison result.
func CompareImagesByUrls(src, dst string) (response.CompareResponse, error) {
	var result response.CompareResponse

	// Use SSI for Euclidean distance
	ssiPercent, err := CompareImagesBySSI(src, dst)
	if err != nil {
		return result, err
	}
	result.Is_same = math.Round(ssiPercent) > 95
	result.Percent = math.Round(ssiPercent*100) / 100
	result.Algorithm = "Structural Similarity Index"

	return result, nil
}

// CompareImagesBySSI compares two images using Structural Similarity Index.
func CompareImagesBySSI(src, dst string) (float64, error) {
	img1, err := GetImage(src)
	if err != nil {
		return 0, err
	}
	img2, err := GetImage(dst)
	if err != nil {
		return 0, err
	}

	return calculateSSI(img1, img2), nil
}

// GetImageExt retrieves the image extension from a given URL.
func GetImageExt(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	buff := make([]byte, 512)
	_, err = response.Body.Read(buff)
	if err != nil {
		return "", errors.New("can't load the content of URL")
	}

	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	case "image/tiff":
		return ".tiff", nil
	default:
		return "", errors.New("unsupported image format")
	}
}

// GetImage retrieves an image from a given URL.
func GetImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("can't load the content of URL: " + url)
	}

	ext, err := GetImageExt(url)
	if err != nil {
		return nil, err
	}

	var img image.Image
	switch ext {
	case ".jpeg", ".jpg", ".png", ".gif":
		img, _, err = image.Decode(response.Body)
	case ".webp":
		img, err = webp.Decode(response.Body)
	case ".tiff", ".tif":
		img, err = tiff.Decode(response.Body)
	default:
		return nil, errors.New("unsupported image format: " + ext)
	}

	if err != nil {
		return nil, errors.New("error decoding image URL \"" + url + "\". Error: " + err.Error())
	}

	return img, nil
}

//Structural Similarity Index
func calculateSSI(img1, img2 image.Image) float64 {
	c1 := 0.0001 // Constant to prevent division by zero

	muX := calculateMean(img1)
	muY := calculateMean(img2)

	sigmaX := calculateStdDev(img1, muX)
	sigmaY := calculateStdDev(img2, muY)
	sigmaXY := calculateCrossCovariance(img1, img2, muX, muY)

	l := (2*muX*muY + c1) / (muX*muX + muY*muY + c1)
	c := (2*sigmaX*sigmaY + c1) / (sigmaX*sigmaX + sigmaY*sigmaY + c1)
	s := (sigmaXY + c1/2) / (sigmaX*sigmaY + c1/2)

	ssi := l * c * s

	// Map SSI to the range [0%, 100%]
	ssiPercentage := (ssi + 1) / 2 * 100

	return ssiPercentage
}

func calculateMean(img image.Image) float64 {
	bounds := img.Bounds()
	total := 0.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			total += float64(r + g + b)
		}
	}

	total /= 3 * float64(bounds.Max.X*bounds.Max.Y)
	return total
}

func calculateStdDev(img image.Image, mean float64) float64 {
	bounds := img.Bounds()
	total := 0.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			total += (float64(r+g+b) - mean) * (float64(r+g+b) - mean)
		}
	}

	total /= 3 * float64(bounds.Max.X*bounds.Max.Y)
	return math.Sqrt(total)
}

func calculateCrossCovariance(img1, img2 image.Image, mean1, mean2 float64) float64 {
	bounds := img1.Bounds()
	total := 0.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			total += (float64(r1+g1+b1) - mean1) * (float64(r2+g2+b2) - mean2)
		}
	}

	total /= 3 * float64(bounds.Max.X*bounds.Max.Y)
	return total
}
