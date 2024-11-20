package utils

import (
	"image"
	"io"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// CalculatePerimeter calculates the perimeter of an image
func CalculatePerimeter(body io.Reader) (int, error) {
	img, _, err := image.Decode(body)
	if err != nil {
		return 0, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	return 2 * (width + height), nil
}
