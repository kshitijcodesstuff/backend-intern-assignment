package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"testing"
)

func TestCalculatePerimeter(t *testing.T) {
	// Create a 100x200 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 200))
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)

	perimeter, err := CalculatePerimeter(&buf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedPerimeter := 2 * (100 + 200)
	if perimeter != expectedPerimeter {
		t.Errorf("Expected perimeter %d, got %d", expectedPerimeter, perimeter)
	}
}

func TestCalculatePerimeterInvalidImage(t *testing.T) {
	// Provide invalid image data
	invalidData := []byte("invalid image data")
	_, err := CalculatePerimeter(bytes.NewReader(invalidData))
	if err == nil {
		t.Errorf("Expected an error for invalid image data")
	}
}
