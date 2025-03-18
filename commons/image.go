package commons

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func CalcImageSize(aspectRatio string, maxWidth, maxHeight int) (w int, h int, e error) {
	parts := strings.Split(aspectRatio, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid aspect ratio %s", aspectRatio)
	}

	widthRatio, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width in aspect ratio")
	}

	heightRatio, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height in aspect ratio")
	}

	widthScale := float64(maxWidth) / float64(widthRatio)
	heightScale := float64(maxHeight) / float64(heightRatio)

	scale := widthScale
	if heightScale < widthScale {
		scale = heightScale
	}

	// Calculate the final dimensions
	width := int(float64(widthRatio) * scale)
	height := int(float64(heightRatio) * scale)

	// Ensure width and height are divisible by 8
	width = width / 8 * 8
	height = height / 8 * 8

	return width, height, nil
}

func ResizeImage(input []byte, perc float64) ([]byte, error) {
	srcImage, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	var (
		buf       bytes.Buffer
		bounds    = srcImage.Bounds()
		newWidth  = float64(bounds.Dx()) * perc
		newHeight = float64(bounds.Dy()) * perc
	)

	dstImage := image.NewRGBA(image.Rect(0, 0, int(newWidth), int(newHeight)))
	draw.CatmullRom.Scale(dstImage, dstImage.Bounds(), srcImage, bounds, draw.Over, nil)

	if err := jpeg.Encode(&buf, dstImage, &jpeg.Options{Quality: 100}); err != nil {
		return nil, fmt.Errorf("error encoding image to JPEG: %v", err)
	}

	return buf.Bytes(), nil
}

func JoinImageHorizontally(images []image.Image) (image.Image, error) {
	var (
		newWidth  = 0
		newHeight = 0
		margin    = 5
	)

	// total size
	for _, img := range images {
		bounds := img.Bounds()
		newWidth += bounds.Dx() + margin
		if bounds.Dy() > newHeight {
			newHeight = bounds.Dy()
		}
	}

	if newHeight == 0 || newWidth == 0 {
		return nil, errors.New("no image provided or failed to determine pixels")
	}

	// canvas
	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(newImage, newImage.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// attach images
	xOffset := 0
	for _, img := range images {
		bounds := img.Bounds()
		draw.Draw(newImage, image.Rect(xOffset, 0, xOffset+bounds.Dx(), bounds.Dy()), img, bounds.Min, draw.Over)
		xOffset += bounds.Dx() + margin
	}

	return newImage, nil
}
