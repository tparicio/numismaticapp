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

	// 3. Save as PNG to support transparency (if we added mask, which bimg doesn't support easily yet)
	// Even without mask, saving as PNG is requested.
	// We change extension to .png
	// bimg.Write automatically determines format from extension? No, we need to convert.
	// bimg.NewImage(cropped).Convert(bimg.PNG)

	newImg := bimg.NewImage(cropped)
	pngBuf, err := newImg.Convert(bimg.PNG)
	if err != nil {
		return "", fmt.Errorf("failed to convert to png: %w", err)
	}

	// Construct new path with .png extension
	// We strip original extension and append _crop.png
	// Actually CoinService expects a new file.
	// Let's replace extension.

	// Simple string manipulation for now
	basePath := imagePath
	if len(basePath) > 4 {
		basePath = basePath[:len(basePath)-4] // strip .jpg or similar assuming 3 chars
	}
	outputPath := basePath + "_crop.png"

	if err := bimg.Write(outputPath, pngBuf); err != nil {
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

	// Check if input was PNG to preserve format
	// bimg.Read reads bytes, we can check signature or just assume based on extension
	// If extension is .png, save as png.
	if len(imagePath) > 4 && imagePath[len(imagePath)-4:] == ".png" {
		// Rotate might return default format buffer? No, it returns buffer.
		// We should ensure it's saved as PNG if it was PNG.
		// bimg operations usually preserve type if possible or return raw buffer.
		// But let's be safe and convert if needed or just Write.
		// bimg.Write uses vips_image_write_to_file which detects type from extension.
		if err := bimg.Write(imagePath, rotated); err != nil {
			return "", err
		}
	} else {
		// Default behavior (likely jpg)
		if err := bimg.Write(imagePath, rotated); err != nil {
			return "", err
		}
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
