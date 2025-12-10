package domain

import "context"

// ImagePaths holds the paths for the processed images.
type ImagePaths struct {
	OriginalFront string
	OriginalBack  string
	CroppedFront  string
	CroppedBack   string
	FinalFront    string
	FinalBack     string
}

// ImageManager defines the interface for the image management module.
type ImageManager interface {
	// ProcessAndSave handles the full lifecycle of coin images:
	// validation, saving original, removing background, and saving initial final version.
	ProcessAndSave(ctx context.Context, coinID string, frontImg, backImg []byte) (*ImagePaths, error)

	// ApplyRotation applies a geometric rotation to the cropped images and updates the final images.
	ApplyRotation(ctx context.Context, coinID string, angleFront, angleBack float64) error
}

// ImageStorage defines the interface for file system operations.
type ImageStorage interface {
	Save(coinID, filename string, data []byte) (string, error)
	Load(coinID, filename string) ([]byte, error)
	Exists(coinID, filename string) bool
	GetPath(coinID, filename string) string
}

// BackgroundRemover defines the interface for the background removal service.
type BackgroundRemover interface {
	RemoveBackground(ctx context.Context, image []byte) ([]byte, error)
}

// ImageProcessor defines the interface for local image transformations.
type ImageProcessor interface {
	ToPNG(image []byte) ([]byte, error)
	Rotate(image []byte, angle float64) ([]byte, error)
}
