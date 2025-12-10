package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"github.com/h2non/bimg"
)

type VipsImageService struct{}

func NewVipsImageService() *VipsImageService {
	return &VipsImageService{}
}

func (s *VipsImageService) ProcessCoinImages(frontPath, backPath string) (string, string, error) {
	// 1. Crop front to circle
	croppedFront, err := s.CropToCircle(frontPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to crop front image: %w", err)
	}

	// 2. Trim front result
	if err := s.Trim(croppedFront); err != nil {
		return "", "", fmt.Errorf("failed to trim front image: %w", err)
	}

	// 3. Trim original front
	if err := s.Trim(frontPath); err != nil {
		return "", "", fmt.Errorf("failed to trim original front image: %w", err)
	}

	// 4. Crop back to circle
	croppedBack, err := s.CropToCircle(backPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to crop back image: %w", err)
	}

	// 5. Trim back result
	if err := s.Trim(croppedBack); err != nil {
		return "", "", fmt.Errorf("failed to trim back image: %w", err)
	}

	// 6. Trim original back
	if err := s.Trim(backPath); err != nil {
		return "", "", fmt.Errorf("failed to trim original back image: %w", err)
	}

	return croppedFront, croppedBack, nil
}

func (s *VipsImageService) Trim(imagePath string) error {
	buffer, err := bimg.Read(imagePath)
	if err != nil {
		return err
	}

	img := bimg.NewImage(buffer)
	size, _ := img.Size()
	fmt.Printf("DEBUG: Trimming %s. Original size: %dx%d\n", imagePath, size.Width, size.Height)

	// Use Process with explicit threshold to be more aggressive
	// Default is usually 10, we try 50 to catch shadows/noise
	options := bimg.Options{
		Trim:      true,
		Threshold: 50,
	}

	trimmed, err := img.Process(options)
	if err != nil {
		return err
	}

	trimmedImg := bimg.NewImage(trimmed)
	newSize, _ := trimmedImg.Size()
	fmt.Printf("DEBUG: Trimmed %s. New size: %dx%d\n", imagePath, newSize.Width, newSize.Height)

	if err := bimg.Write(imagePath, trimmed); err != nil {
		return err
	}

	return nil
}

func (s *VipsImageService) CropToCircle(imagePath string) (string, error) {
	buffer, err := bimg.Read(imagePath)
	if err != nil {
		return "", err
	}

	img := bimg.NewImage(buffer)
	size, err := img.Size()
	if err != nil {
		return "", err
	}

	// 1. Determine smallest dimension for square crop
	dimension := size.Width
	if size.Height < size.Width {
		dimension = size.Height
	}

	// 2. Crop to square centered
	cropped, err := img.Crop(dimension, dimension, bimg.GravityCentre)
	if err != nil {
		return "", err
	}

	// 3. Save as PNG
	newImg := bimg.NewImage(cropped)
	pngBuf, err := newImg.Convert(bimg.PNG)
	if err != nil {
		return "", fmt.Errorf("failed to convert to png: %w", err)
	}

	basePath := imagePath
	if len(basePath) > 4 {
		basePath = basePath[:len(basePath)-4]
	}
	outputPath := basePath + "_crop.png"

	if err := bimg.Write(outputPath, pngBuf); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *VipsImageService) GetMetadata(imagePath string) (width, height int, size int64, mimeType string, err error) {
	buffer, err := bimg.Read(imagePath)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("failed to read image: %w", err)
	}

	img := bimg.NewImage(buffer)
	dims, err := img.Size()
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("failed to get image size: %w", err)
	}

	size = int64(len(buffer))
	typeName := img.Type()
	mimeType = "image/" + typeName

	return dims.Width, dims.Height, size, mimeType, nil
}

func (s *VipsImageService) CropToContent(data []byte) ([]byte, error) {
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
	// maxX and maxY are inclusive in the loop, but Rect is exclusive at max
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
