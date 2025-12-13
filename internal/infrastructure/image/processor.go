package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
	"strings"

	"github.com/disintegration/imaging"
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
	img, err := imaging.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image for trim: %w", err)
	}

	// Scan for non-transparent pixels to determine bounds
	bounds := img.Bounds()
	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 { // Simple threshold, can increase if noisy
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
		// Empty image, do nothing
		return nil
	}

	// Add slight padding to avoid cutting too close
	padding := 5
	minX = max(bounds.Min.X, minX-padding)
	minY = max(bounds.Min.Y, minY-padding)
	maxX = min(bounds.Max.X-1, maxX+padding)
	maxY = min(bounds.Max.Y-1, maxY+padding)

	rect := image.Rect(minX, minY, maxX+1, maxY+1)
	trimmedImg := imaging.Crop(img, rect)

	if err := imaging.Save(trimmedImg, imagePath); err != nil {
		return fmt.Errorf("failed to save trimmed image: %w", err)
	}

	slog.Debug("Trimmed image", "path", imagePath, "old_bounds", bounds, "new_bounds", rect)
	return nil
}

func (s *VipsImageService) CropToCircle(imagePath string) (string, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image for circle crop: %w", err)
	}

	bounds := img.Bounds()
	dim := bounds.Dx()
	if bounds.Dy() < dim {
		dim = bounds.Dy()
	}

	// 1. Crop to square center
	img = imaging.CropCenter(img, dim, dim)

	// 2. Create circular mask
	// Create a new RGBA image for the result
	dst := image.NewRGBA(image.Rect(0, 0, dim, dim))

	// Center and Radius
	cx, cy := float64(dim)/2, float64(dim)/2
	r := float64(dim) / 2

	// Iterate over pixels to apply circular mask
	for y := 0; y < dim; y++ {
		for x := 0; x < dim; x++ {
			dx, dy := float64(x)-cx, float64(y)-cy
			if dx*dx+dy*dy <= r*r {
				dst.Set(x, y, img.At(x, y))
			} else {
				dst.Set(x, y, color.Transparent)
			}
		}
	}

	// 3. Save as _crop.png
	basePath := imagePath
	if len(basePath) > 4 && strings.Contains(basePath, ".") {
		// strip extension
		idx := strings.LastIndex(basePath, ".")
		if idx > 0 {
			basePath = basePath[:idx]
		}
	}
	outputPath := basePath + "_crop.png"

	if err := imaging.Save(dst, outputPath); err != nil {
		return "", fmt.Errorf("failed to save circle crop: %w", err)
	}

	return outputPath, nil
}

func (s *VipsImageService) GetMetadata(imagePath string) (width, height int, size int64, mimeType string, err error) {
	// Get file info for size
	info, err := os.Stat(imagePath)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("failed to stat file: %w", err)
	}
	size = info.Size()

	// Decode config for dimensions and format
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close image file", "path", imagePath, "error", err)
		}
	}()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("failed to decode image config: %w", err)
	}

	mimeType = "image/" + format
	return config.Width, config.Height, size, mimeType, nil
}

func (s *VipsImageService) CropToContent(data []byte) ([]byte, error) {
	// Reusing the pure Go implementation I viewed earlier, but adapted to accept bytes and return bytes
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Reuse logic from Trim roughly?
	// Just return logic similar to Trim but for []byte
	// Scan bounds
	bounds := img.Bounds()
	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	found := false

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
		return data, nil
	}

	// Padding
	padding := int(float64(max(bounds.Dx(), bounds.Dy())) * 0.05)
	minX = max(bounds.Min.X, minX-padding)
	minY = max(bounds.Min.Y, minY-padding)
	maxX = min(bounds.Max.X-1, maxX+padding)
	maxY = min(bounds.Max.Y-1, maxY+padding)

	rect := image.Rect(minX, minY, maxX+1, maxY+1)

	// Crop
	trimmed := imaging.Crop(img, rect)

	// Encode to PNG buffer
	var buf bytes.Buffer
	if err := png.Encode(&buf, trimmed); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *VipsImageService) GenerateThumbnail(imagePath string, width int) (string, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image for thumbnail: %w", err)
	}

	// Resize (Fit width, maintain aspect ratio)
	newImg := imaging.Resize(img, width, 0, imaging.Lanczos)

	basePath := imagePath
	if idx := strings.LastIndex(basePath, "."); idx > 0 {
		basePath = basePath[:idx]
	}
	outputPath := basePath + "_thumb.png"

	if err := imaging.Save(newImg, outputPath); err != nil {
		return "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return outputPath, nil
}

func (s *VipsImageService) Rotate(imagePath string, angle float64) error {
	if angle == 0 {
		return nil
	}

	img, err := imaging.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open for rotation: %w", err)
	}

	slog.Info("Rotating image with imaging", "path", imagePath, "angle", angle)

	// Rotate (imaging rotates CCW)
	rotatedImg := imaging.Rotate(img, -angle, color.Transparent)

	if err := imaging.Save(rotatedImg, imagePath); err != nil {
		return fmt.Errorf("failed to save rotated image: %w", err)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
