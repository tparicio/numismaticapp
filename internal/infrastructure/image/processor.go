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

	// 4. Rotate back
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

	if len(imagePath) > 4 && imagePath[len(imagePath)-4:] == ".png" {
		if err := bimg.Write(imagePath, rotated); err != nil {
			return "", err
		}
	} else {
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

	size = int64(len(buffer))
	typeName := img.Type()
	mimeType = "image/" + typeName

	return dims.Width, dims.Height, size, mimeType, nil
}
