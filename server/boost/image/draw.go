package image

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

func DrawRectangles(inputPath, outputPath string, rectColor color.Color, thickness int, rects ...image.Rectangle) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer input.Close()

	inputImage, _, err := image.Decode(input)
	if err != nil {
		return err
	}

	bounds := inputImage.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, inputImage, bounds.Min, draw.Src)

	for _, rect := range rects {
		for t := 0; t < thickness; t++ {
			r := image.Rect(rect.Min.X-t, rect.Min.Y-t, rect.Max.X+t, rect.Max.Y+t)
			for x := r.Min.X; x < r.Max.X; x++ {
				if x >= bounds.Min.X && x < bounds.Max.X {
					if r.Min.Y >= bounds.Min.Y {
						rgba.Set(x, r.Min.Y, rectColor)
					}
					if r.Max.Y-1 < bounds.Max.Y {
						rgba.Set(x, r.Max.Y-1, rectColor)
					}
				}
			}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				if y >= bounds.Min.Y && y < bounds.Max.Y {
					if r.Min.X >= bounds.Min.X {
						rgba.Set(r.Min.X, y, rectColor)
					}
					if r.Max.X-1 < bounds.Max.X {
						rgba.Set(r.Max.X-1, y, rectColor)
					}
				}
			}
		}
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer output.Close()

	if err = jpeg.Encode(output, rgba, nil); err != nil {
		return err
	}
	return nil
}
