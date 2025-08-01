package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageService struct {
	thumbnailSizes []ThumbnailSize
	quality        int
}

type ThumbnailSize struct {
	Width  int
	Height int
	Name   string
}

type ProcessedImage struct {
	Data        []byte
	ContentType string
	Width       int
	Height      int
}

func NewImageService() *ImageService {
	return &ImageService{
		thumbnailSizes: []ThumbnailSize{
			{Width: 150, Height: 150, Name: "small"},
			{Width: 300, Height: 300, Name: "medium"},
			{Width: 600, Height: 600, Name: "large"},
		},
		quality: 85,
	}
}

func (is *ImageService) IsImageFile(contentType string) bool {
	imageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/bmp":  true,
	}
	return imageTypes[contentType]
}

func (is *ImageService) ProcessImage(reader io.Reader, contentType string) (*ProcessedImage, error) {
	// Read image data
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Decode image
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Optimize image if needed
	optimizedData, optimizedType, err := is.optimizeImage(img, format, contentType)
	if err != nil {
		// If optimization fails, return original
		return &ProcessedImage{
			Data:        data,
			ContentType: contentType,
			Width:       width,
			Height:      height,
		}, nil
	}

	return &ProcessedImage{
		Data:        optimizedData,
		ContentType: optimizedType,
		Width:       width,
		Height:      height,
	}, nil
}

func (is *ImageService) GenerateThumbnails(reader io.Reader, contentType string) (map[string]*ProcessedImage, error) {
	// Read and decode image
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	thumbnails := make(map[string]*ProcessedImage)

	for _, size := range is.thumbnailSizes {
		// Resize image maintaining aspect ratio
		thumbnail := imaging.Fit(img, size.Width, size.Height, imaging.Lanczos)

		// Encode thumbnail
		var buf bytes.Buffer
		var outputType string

		if strings.Contains(contentType, "png") {
			err = png.Encode(&buf, thumbnail)
			outputType = "image/png"
		} else {
			// Default to JPEG for all other formats
			err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: is.quality})
			outputType = "image/jpeg"
		}

		if err != nil {
			continue // Skip this thumbnail on error
		}

		bounds := thumbnail.Bounds()
		thumbnails[size.Name] = &ProcessedImage{
			Data:        buf.Bytes(),
			ContentType: outputType,
			Width:       bounds.Dx(),
			Height:      bounds.Dy(),
		}
	}

	return thumbnails, nil
}

func (is *ImageService) optimizeImage(img image.Image, format, contentType string) ([]byte, string, error) {
	var buf bytes.Buffer
	var outputType string

	// Convert to RGB if needed and optimize
	bounds := img.Bounds()
	optimized := imaging.Clone(img)

	// Apply optimization based on size
	maxDimension := 4096
	if bounds.Dx() > maxDimension || bounds.Dy() > maxDimension {
		optimized = imaging.Fit(img, maxDimension, maxDimension, imaging.Lanczos)
	}

	// Encode with optimization
	if strings.Contains(contentType, "png") && format == "png" {
		err := png.Encode(&buf, optimized)
		outputType = "image/png"
		if err != nil {
			return nil, "", err
		}
	} else {
		// Convert to JPEG for better compression
		err := jpeg.Encode(&buf, optimized, &jpeg.Options{Quality: is.quality})
		outputType = "image/jpeg"
		if err != nil {
			return nil, "", err
		}
	}

	return buf.Bytes(), outputType, nil
}

func (is *ImageService) ValidateImageSize(width, height int, maxSizeMB float64) error {
	const maxDimension = 8192 // 8K resolution
	
	if width > maxDimension || height > maxDimension {
		return fmt.Errorf("image dimensions too large: %dx%d (max: %dx%d)", 
			width, height, maxDimension, maxDimension)
	}

	// Estimate file size (rough calculation)
	estimatedSizeMB := float64(width*height*3) / (1024 * 1024) // Assuming 3 bytes per pixel
	if estimatedSizeMB > maxSizeMB {
		return fmt.Errorf("estimated image size too large: %.2fMB (max: %.2fMB)", 
			estimatedSizeMB, maxSizeMB)
	}

	return nil
}

func (is *ImageService) GetImageInfo(reader io.Reader) (int, int, string, error) {
	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to decode image config: %w", err)
	}

	return config.Width, config.Height, format, nil
}