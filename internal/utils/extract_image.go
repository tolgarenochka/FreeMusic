package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

// CreateWhiteQuestionImage ...
func CreateWhiteQuestionImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Рисуем знак вопроса (можете изменить на другой символ или изображение)
	questionMark := drawFont()
	draw.Draw(img, image.Rect(0, 0, width, height), &image.Uniform{color.Black}, image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(5, 5, width-5, height-5), questionMark, image.Point{}, draw.Over)

	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)

	return buf.Bytes()
}

// drawFont ...
func drawFont() *image.RGBA {
	width, height := 10, 20
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	black := color.RGBA{A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (x > 1 && x < width-2) || (y > 3 && y < height-4) {
				img.Set(x, y, white)
			} else {
				img.Set(x, y, black)
			}
		}
	}

	return img
}
