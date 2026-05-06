package nasa

import (
	"crypto/rand"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
}

type Generator struct {
	UploadDir string
	OutputDir string
	mu        sync.RWMutex
	cache     map[string]image.Image
}

func NewGenerator(uploadDir, outputDir string) *Generator {
	if outputDir == "" {
		outputDir = os.TempDir()
	}

	g := &Generator{
		UploadDir: uploadDir,
		OutputDir: outputDir,
		cache:     make(map[string]image.Image),
	}

	// Pre-warm cache for faster first request
	g.prewarmCache()

	return g
}

// prewarmCache loads all letter variants into memory at startup
func (g *Generator) prewarmCache() {
	for _, letter := range letters {
		for v := 1; v <= 3; v++ {
			path := filepath.Join(g.UploadDir, fmt.Sprintf("%c%d.jpg", letter, v))
			if _, err := os.Stat(path); err == nil {
				file, err := os.Open(path)
				if err != nil {
					continue
				}
				img, _, err := image.Decode(file)
				file.Close()
				if err == nil {
					g.storeCachedImage(path, img)
				}
			}
		}
	}
}

func (g *Generator) GenerateNameImage(name string) (string, error) {
	var images []image.Image

	for _, ch := range strings.ToLower(name) {
		if !strings.ContainsRune(letters, ch) {
			continue
		}

		img, err := g.randomLetterImage(ch)
		if err != nil {
			return "", err
		}
		if img != nil {
			images = append(images, img)
		}
	}

	if len(images) == 0 {
		return "", fmt.Errorf("image generation failed")
	}

	totalWidth := 0
	maxHeight := 0
	for _, img := range images {
		bounds := img.Bounds()
		totalWidth += bounds.Dx()
		if bounds.Dy() > maxHeight {
			maxHeight = bounds.Dy()
		}
	}

	final := image.NewRGBA(image.Rect(0, 0, totalWidth, maxHeight))
	x := 0
	for _, img := range images {
		bounds := img.Bounds()
		target := image.Rect(x, 0, x+bounds.Dx(), bounds.Dy())
		draw.Draw(final, target, img, bounds.Min, draw.Src)
		x += bounds.Dx()
	}

	if err := os.MkdirAll(g.OutputDir, 0o755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	outputPath := filepath.Join(g.OutputDir, fmt.Sprintf("nasa_name_%d.jpg", time.Now().UnixNano()))
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, final, &jpeg.Options{Quality: 82}); err != nil {
		return "", fmt.Errorf("encode jpeg: %w", err)
	}

	return outputPath, nil
}

func (g *Generator) randomLetterImage(letter rune) (image.Image, error) {
	// Pick random variant 1, 2, or 3
	variant, err := randomInt(1, 4) // 1-3
	if err != nil {
		return nil, err
	}
	path := filepath.Join(g.UploadDir, fmt.Sprintf("%c%d.jpg", letter, variant))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	if img := g.cachedImage(path); img != nil {
		return img, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	g.storeCachedImage(path, img)
	return img, nil
}

func (g *Generator) cachedImage(path string) image.Image {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.cache[path]
}

func (g *Generator) storeCachedImage(path string, img image.Image) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache[path] = img
}

func randomInt(min, max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + min, nil
}

