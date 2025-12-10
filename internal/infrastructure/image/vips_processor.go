package image

import (
	"fmt"

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

var _ domain.ImageProcessor = (*VipsProcessor)(nil)
