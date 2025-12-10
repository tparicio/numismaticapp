package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/h2non/bimg"
)

type VipsProcessor struct{}

func NewVipsProcessor() *VipsProcessor {
	return &VipsProcessor{}
}

func (p *VipsProcessor) ToPNG(image []byte) ([]byte, error) {
	newImage, err := bimg.NewImage(image).Convert(bimg.PNG)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PNG: %w", err)
	}
	return newImage, nil
}

func (p *VipsProcessor) Rotate(image []byte, angle float64) ([]byte, error) {
	if angle == 0 {
		return image, nil
	}
	newImage, err := bimg.NewImage(image).Rotate(bimg.Angle(angle))
	if err != nil {
		return nil, fmt.Errorf("failed to rotate image: %w", err)
	}
	return newImage, nil
}

func (p *VipsProcessor) CropToContent(data []byte) ([]byte, error) {
	// Use standard image library to avoid bimg/libvips issues with AutoCrop
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	found := false

	// Scan for non-transparent pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				found = true
			}
		}
	}

	if !found {
		// Return original if empty
		return data, nil
	}

	// Crop rect
	cropRect := image.Rect(minX, minY, maxX+1, maxY+1)
	width := cropRect.Dx()
	height := cropRect.Dy()

	// Calculate square size with padding
	maxDim := width
	if height > maxDim {
		maxDim = height
	}
	padding := int(float64(maxDim) * 0.05)
	finalSize := maxDim + (padding * 2)

	// Create new square image
	newImg := image.NewRGBA(image.Rect(0, 0, finalSize, finalSize))

	// Calculate offset to center
	offsetX := (finalSize - width) / 2
	offsetY := (finalSize - height) / 2

	// Draw cropped part onto new image
	draw.Draw(newImg, image.Rect(offsetX, offsetY, offsetX+width, offsetY+height), img, cropRect.Min, draw.Over)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, newImg); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}

	return buf.Bytes(), nil
}

var _ domain.ImageProcessor = (*VipsProcessor)(nil)
