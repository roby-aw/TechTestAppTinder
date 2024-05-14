package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

func DecodeImage(data []byte, ext string) (image.Image, error) {
	switch ext {
	case ".jpg":
		return jpeg.Decode(bytes.NewReader(data))
	case ".jpeg":
		return jpeg.Decode(bytes.NewReader(data))
	case ".JPEG":
		return jpeg.Decode(bytes.NewReader(data))
	case ".JPG":
		return jpeg.Decode(bytes.NewReader(data))
	case ".png":
		return png.Decode(bytes.NewReader(data))
	default:
		return nil, errors.New("image not support")
	}
}
