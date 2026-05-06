package nasa

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateNameImage(t *testing.T) {
	uploadDir := t.TempDir()
	outputDir := t.TempDir()

	writeLetterImage(t, uploadDir, "0", "a", 11, 7)
	writeLetterImage(t, uploadDir, "1", "b", 13, 9)

	generator := NewGenerator(uploadDir, outputDir)
	path, err := generator.GenerateNameImage("ab!")
	if err != nil {
		t.Fatalf("GenerateNameImage returned error: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open generated image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("decode generated image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 24 {
		t.Fatalf("width = %d, want 24", bounds.Dx())
	}
	if bounds.Dy() != 9 {
		t.Fatalf("height = %d, want 9", bounds.Dy())
	}
}

func writeLetterImage(t *testing.T, uploadDir, folder, letter string, width, height int) {
	t.Helper()

	dir := filepath.Join(uploadDir, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create upload dir: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 50, B: 50, A: 255})
		}
	}

	file, err := os.Create(filepath.Join(dir, letter+".jpg"))
	if err != nil {
		t.Fatalf("create letter image: %v", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}
}
