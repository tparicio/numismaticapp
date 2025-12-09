package image

import (
	"fmt"

	"github.com/h2non/bimg"
)

type VipsImageService struct{}

func NewVipsImageService() *VipsImageService {
	return &VipsImageService{}
}

func (s *VipsImageService) ProcessCoinImages(frontPath, backPath string, rotationAngle float64) (string, string, error) {
	// 1. Crop front to circle
	croppedFront, err := s.CropToCircle(frontPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to crop front image: %w", err)
	}

	// 2. Rotate front
	rotatedFront, err := s.Rotate(croppedFront, rotationAngle)
	if err != nil {
		return "", "", fmt.Errorf("failed to rotate front image: %w", err)
	}

	// 3. Crop back to circle
	croppedBack, err := s.CropToCircle(backPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to crop back image: %w", err)
	}

	// 4. Rotate back (assuming same angle correction needed, or 0 if independent?
	// Usually coins are flipped 180 or 0 degrees relative to each other (coin vs medal alignment).
	// For now, we apply the same correction assuming the user took photos roughly same orientation,
	// or we might need separate angles from Gemini. The prompt asked for "vertical_correction_angle" generally.
	// Let's apply it to front. For back, it might be different.
	// Ideally Gemini gives two angles. For MVP, let's apply to front and assume back is manual or same.)

	// Actually, let's just rotate front for now as per prompt "vertical_correction_angle".
	// If the user takes photos separately, they might have different rotations.
	// Let's apply the same rotation to back as a best guess, or 0.
	rotatedBack, err := s.Rotate(croppedBack, rotationAngle)
	if err != nil {
		return "", "", fmt.Errorf("failed to rotate back image: %w", err)
	}

	return rotatedFront, rotatedBack, nil
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

	// Simple "smart" crop - assuming coin is centered-ish.
	// Real "detect circle" is hard with just bimg (needs OpenCV).
	// For MVP, we'll do a center square crop which is often good enough if user centers the coin.
	// Or we can use bimg.SmartCrop.

	dimension := size.Width
	if size.Height < size.Width {
		dimension = size.Height
	}

	// Crop to square first
	cropped, err := img.Crop(dimension, dimension, bimg.GravityCentre)
	if err != nil {
		return "", err
	}

	// Create a new image from cropped buffer
	// Note: bimg doesn't have a direct "mask to circle" easily without SVG overlay or similar.
	// For now, we will just return the square cropped image.
	// CSS radius: 50% will handle the visual circle.
	// If we strictly need to make pixels transparent outside circle, it requires more complex vips operations.
	// Let's stick to Square Crop for MVP backend processing.

	outputPath := imagePath + "_processed.jpg"
	if err := bimg.Write(outputPath, cropped); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *VipsImageService) Rotate(imagePath string, angle float64) (string, error) {
	if angle == 0 {
		return imagePath, nil
	}

	buffer, err := bimg.Read(imagePath)
	if err != nil {
		return "", err
	}

	img := bimg.NewImage(buffer)
	rotated, err := img.Rotate(bimg.Angle(angle))
	if err != nil {
		return "", err
	}

	// Overwrite or new file? Let's overwrite the processed file
	if err := bimg.Write(imagePath, rotated); err != nil {
		return "", err
	}

	return imagePath, nil
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

	// Get file size
	size = int64(len(buffer))

	// Determine mime type (basic check based on type)
	// bimg can give us the type
	typeName := img.Type()
	mimeType = "image/" + typeName

	return dims.Width, dims.Height, size, mimeType, nil
}
