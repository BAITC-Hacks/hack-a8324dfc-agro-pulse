package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// Если красных пикселей больше 20% кадра, считаем изображение браком.
const redShareLimit = 0.20

// isRed проверяет, «красный» ли пиксель.
// Каналы r, g, b лежат в диапазоне 0–65535.
func isRed(r, g, b uint32) bool {
	return r > 50000 && r > g*2 && r > b*2
}

// classify открывает файл и возвращает "OK" или "DEFECT".
func classify(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	total := bounds.Dx() * bounds.Dy()
	if total == 0 {
		return "", fmt.Errorf("пустое изображение")
	}

	red := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if isRed(r, g, b) {
				red++
			}
		}
	}

	if float64(red)/float64(total) > redShareLimit {
		return "DEFECT", nil
	}
	return "OK", nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <image>")
		os.Exit(1)
	}

	result, err := classify(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
