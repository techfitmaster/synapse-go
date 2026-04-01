package imaging

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

// Format represents an image output format.
type Format int

const (
	// JPEG format with lossy compression.
	JPEG Format = iota
	// PNG format with lossless compression.
	PNG
)

// ResizeOption configures image resizing behavior.
type ResizeOption struct {
	MaxWidth  int // maximum width (0 = no limit)
	MaxHeight int // maximum height (0 = no limit)
	Quality   int // JPEG quality 1-100 (default 85, ignored for PNG)
	Format    Format
}

// Decode reads an image from the reader, auto-detecting format (JPEG/PNG).
func Decode(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

// Resize scales an image to fit within MaxWidth/MaxHeight while preserving aspect ratio.
// Returns the resized image bytes in the specified format.
func Resize(img image.Image, opt ResizeOption) ([]byte, error) {
	if opt.Quality == 0 {
		opt.Quality = 85
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dstW, dstH := fitDimensions(srcW, srcH, opt.MaxWidth, opt.MaxHeight)

	// Skip resize if image is already smaller
	if dstW >= srcW && dstH >= srcH {
		return encode(img, opt)
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return encode(dst, opt)
}

// Compress re-encodes an image at the given quality without resizing.
func Compress(img image.Image, quality int, format Format) ([]byte, error) {
	return encode(img, ResizeOption{Quality: quality, Format: format})
}

// Thumbnail generates a square thumbnail of the given size, center-cropped.
func Thumbnail(img image.Image, size int, quality int) ([]byte, error) {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Determine crop region (center square)
	cropSize := srcW
	if srcH < cropSize {
		cropSize = srcH
	}
	x0 := (srcW - cropSize) / 2
	y0 := (srcH - cropSize) / 2

	cropped := image.NewRGBA(image.Rect(0, 0, cropSize, cropSize))
	draw.Copy(cropped, image.Point{}, img, image.Rect(x0+bounds.Min.X, y0+bounds.Min.Y, x0+bounds.Min.X+cropSize, y0+bounds.Min.Y+cropSize), draw.Src, nil)

	// Scale to target size
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

	if quality == 0 {
		quality = 80
	}
	return encode(dst, ResizeOption{Quality: quality, Format: JPEG})
}

// fitDimensions calculates target dimensions preserving aspect ratio.
func fitDimensions(srcW, srcH, maxW, maxH int) (int, int) {
	if maxW == 0 && maxH == 0 {
		return srcW, srcH
	}

	ratioW := float64(1)
	ratioH := float64(1)

	if maxW > 0 && srcW > maxW {
		ratioW = float64(maxW) / float64(srcW)
	}
	if maxH > 0 && srcH > maxH {
		ratioH = float64(maxH) / float64(srcH)
	}

	ratio := ratioW
	if ratioH < ratio {
		ratio = ratioH
	}

	dstW := int(float64(srcW) * ratio)
	dstH := int(float64(srcH) * ratio)
	if dstW < 1 {
		dstW = 1
	}
	if dstH < 1 {
		dstH = 1
	}
	return dstW, dstH
}

func encode(img image.Image, opt ResizeOption) ([]byte, error) {
	var buf bytes.Buffer
	switch opt.Format {
	case PNG:
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("imaging: png encode: %w", err)
		}
	default: // JPEG
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: opt.Quality}); err != nil {
			return nil, fmt.Errorf("imaging: jpeg encode: %w", err)
		}
	}
	return buf.Bytes(), nil
}
