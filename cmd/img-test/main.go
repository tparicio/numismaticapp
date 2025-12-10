package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	infraImage "github.com/antonioparicio/numismaticapp/internal/infrastructure/image"
	infraStorage "github.com/antonioparicio/numismaticapp/internal/infrastructure/storage"
)

// ImageServiceImpl implements domain.ImageManager
type ImageServiceImpl struct {
	storage   domain.ImageStorage
	rembg     domain.BackgroundRemover
	processor domain.ImageProcessor
}

func NewImageService(storage domain.ImageStorage, rembg domain.BackgroundRemover, processor domain.ImageProcessor) *ImageServiceImpl {
	return &ImageServiceImpl{
		storage:   storage,
		rembg:     rembg,
		processor: processor,
	}
}

func (s *ImageServiceImpl) ProcessAndSave(ctx context.Context, coinID string, frontImg, backImg []byte) (*domain.ImagePaths, error) {
	// 1. Validation (simple check for now)
	if len(frontImg) == 0 || len(backImg) == 0 {
		return nil, fmt.Errorf("front and back images are required")
	}

	// 2. Save Originals
	if _, err := s.storage.Save(coinID, "original_front.jpg", frontImg); err != nil {
		return nil, err
	}
	if _, err := s.storage.Save(coinID, "original_back.jpg", backImg); err != nil {
		return nil, err
	}

	// 3. Process Images (Concurrent)
	type result struct {
		side string
		data []byte
		err  error
	}
	results := make(chan result, 2)

	process := func(side string, img []byte) {
		// Remove background
		cropped, err := s.rembg.RemoveBackground(ctx, img)
		if err != nil {
			results <- result{side: side, err: fmt.Errorf("rembg failed: %w", err)}
			return
		}
		results <- result{side: side, data: cropped}
	}

	go process("front", frontImg)
	go process("back", backImg)

	var croppedFront, croppedBack []byte
	for i := 0; i < 2; i++ {
		res := <-results
		if res.err != nil {
			return nil, res.err
		}
		if res.side == "front" {
			croppedFront = res.data
		} else {
			croppedBack = res.data
		}
	}

	// 4. Save Cropped
	if _, err := s.storage.Save(coinID, "cropped_front.png", croppedFront); err != nil {
		return nil, err
	}
	if _, err := s.storage.Save(coinID, "cropped_back.png", croppedBack); err != nil {
		return nil, err
	}

	// 5. Initial Final (Copy of Cropped, 0 rotation)
	if _, err := s.storage.Save(coinID, "final_front.png", croppedFront); err != nil {
		return nil, err
	}
	if _, err := s.storage.Save(coinID, "final_back.png", croppedBack); err != nil {
		return nil, err
	}

	return &domain.ImagePaths{
		OriginalFront: s.storage.GetPath(coinID, "original_front.jpg"),
		OriginalBack:  s.storage.GetPath(coinID, "original_back.jpg"),
		CroppedFront:  s.storage.GetPath(coinID, "cropped_front.png"),
		CroppedBack:   s.storage.GetPath(coinID, "cropped_back.png"),
		FinalFront:    s.storage.GetPath(coinID, "final_front.png"),
		FinalBack:     s.storage.GetPath(coinID, "final_back.png"),
	}, nil
}

func (s *ImageServiceImpl) ApplyRotation(ctx context.Context, coinID string, angleFront, angleBack float64) error {
	// Front
	if angleFront != 0 {
		croppedFront, err := s.storage.Load(coinID, "cropped_front.png")
		if err != nil {
			return err
		}
		rotatedFront, err := s.processor.Rotate(croppedFront, angleFront)
		if err != nil {
			return err
		}
		if _, err := s.storage.Save(coinID, "final_front.png", rotatedFront); err != nil {
			return err
		}
	}

	// Back
	if angleBack != 0 {
		croppedBack, err := s.storage.Load(coinID, "cropped_back.png")
		if err != nil {
			return err
		}
		rotatedBack, err := s.processor.Rotate(croppedBack, angleBack)
		if err != nil {
			return err
		}
		if _, err := s.storage.Save(coinID, "final_back.png", rotatedBack); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Setup
	cwd, _ := os.Getwd()
	storage := infraStorage.NewLocalStorage(filepath.Join(cwd, "storage"))
	rembgURL := os.Getenv("REMBG_URL")
	if rembgURL == "" {
		rembgURL = "http://localhost:5000"
	}
	rembg := infraImage.NewRembgClient(rembgURL)
	processor := infraImage.NewVipsProcessor()
	service := NewImageService(storage, rembg, processor)

	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
		coinID := uuid.New().String()

		fileFront, err := c.FormFile("front")
		if err != nil {
			return c.Status(400).SendString("front image missing")
		}
		fileBack, err := c.FormFile("back")
		if err != nil {
			return c.Status(400).SendString("back image missing")
		}

		read := func(fh *multipart.FileHeader) ([]byte, error) {
			f, err := fh.Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return io.ReadAll(f)
		}

		frontBytes, err := read(fileFront)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		backBytes, err := read(fileBack)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		paths, err := service.ProcessAndSave(c.Context(), coinID, frontBytes, backBytes)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(fiber.Map{
			"coin_id": coinID,
			"paths":   paths,
		})
	})

	app.Post("/rotate/:coinID", func(c *fiber.Ctx) error {
		coinID := c.Params("coinID")
		var body struct {
			AngleFront float64 `json:"angle_front"`
			AngleBack  float64 `json:"angle_back"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		if err := service.ApplyRotation(c.Context(), coinID, body.AngleFront, body.AngleBack); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.SendString("Rotation applied")
	})

	log.Fatal(app.Listen(":3000"))
}
