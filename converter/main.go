package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type HeaderDTO struct {
	NumPixels uint
	Pg        []PixelGrid
	Gifs      []Gif
}

type PixelGrid struct {
	Name      string
	ColorData []uint32
}

type Gif struct {
	Name   string
	Delay  []int
	Frames [][]uint32
}

var width int
var height int
var print bool

func init() {
	//Parse flags
	w := flag.Int("w", 14, "width of the board")
	h := flag.Int("h", 12, "height of the board")
	p := flag.Bool("p", false, "weather you should print to stdout")
	flag.Parse()
	width = *w
	height = *h
	print = *p
}

func normalizeImage(img image.Image, name string) image.Image {
	if img.Bounds().Max.X != width || img.Bounds().Max.Y != height {
		log.Printf("[WARN] %v is not the right size, resizing using nearestneighbor\n", name)
		img = imaging.Fit(img, width, height, imaging.NearestNeighbor)
		dst := image.NewNRGBA(img.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				r, g, b = r>>8, g>>8, b>>8
				dst.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}
		ctx := gg.NewContext(width, height)
		ctx.SetColor(color.Black)
		ctx.DrawRectangle(0, 0, float64(width), float64(height))
		ctx.Fill()
		ctx.DrawImageAnchored(dst, width/2, height/2, 0.5, 0.5)
		img = ctx.Image()
	}
	return img
}

func extractColors(img image.Image) []uint32 {
	colors := []uint32{}
	writeColor := func(r, g, b, a uint32) {
		rgb := r & 255
		rgb = (rgb << 8) + g&255
		rgb = (rgb << 8) + b&255
		colors = append(colors, rgb)
	}
	for h := height - 1; h >= 0; h-- {
		//If even go from left to right else go right to left
		if h%2 != 0 {
			for w := 0; w < width; w++ {
				c := img.At(w, h)
				writeColor(c.RGBA())
			}
		} else {
			for w := width - 1; w >= 0; w-- {
				c := img.At(w, h)
				writeColor(c.RGBA())
			}
		}
	}
	return colors
}

func staticImageWorker(imagePathChan <-chan string, pixelChan chan<- PixelGrid) {
	for path := range imagePathChan {
		imgFile, _ := os.Open(path)
		img, _, err := image.Decode(imgFile)
		if err != nil {
			panic("Failed to decode image")
		}
		img = normalizeImage(img, path)
		colors := extractColors(img)
		pixelChan <- PixelGrid{ColorData: colors, Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))}
		imgFile.Close()
	}
}

func gifWorker(imagePathChan <-chan string, gifDataChan chan<- Gif) {
	for path := range imagePathChan {
		f, _ := os.Open(path)
		gif, _ := gif.DecodeAll(f)
		gifDto := Gif{}
		gifDto.Delay = gif.Delay
		gifDto.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		for i, frame := range gif.Image {
			img := normalizeImage(frame, fmt.Sprintf("%v_%v", path, i))
			colors := extractColors(img)
			gifDto.Frames = append(gifDto.Frames, colors)
		}
		f.Close()
		gifDataChan <- gifDto
	}
}

func main() {
	//Load header output template
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	t, err := template.New("").Funcs(funcMap).ParseFiles("headers.tmpl")
	if err != nil {
		panic("Failed to parse header template: " + err.Error())
	}

	//Establish channels
	imagePathChan := make(chan string, 100)
	pixelDataChan := make(chan PixelGrid, 100)
	gifPathChan := make(chan string, 100)
	gifDataChan := make(chan Gif, 1000)

	//Walk through images directory
	pg := []PixelGrid{}
	gifs := []Gif{}
	for i := 0; i < 5; i++ {
		go staticImageWorker(imagePathChan, pixelDataChan)
		go gifWorker(gifPathChan, gifDataChan)
	}
	numImages := 0
	numGifs := 0
	filepath.Walk("../img", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpg" {
			imagePathChan <- path
			numImages++
		}
		if filepath.Ext(path) == ".gif" {
			gifPathChan <- path
			numGifs++
		}
		return nil
	})
	close(imagePathChan)
	close(gifPathChan)
	for a := 0; a < numImages; a++ {
		pg = append(pg, <-pixelDataChan)
	}
	close(pixelDataChan)
	for a := 0; a < numGifs; a++ {
		gifs = append(gifs, <-gifDataChan)
	}
	close(gifDataChan)

	outFile, _ := os.OpenFile("../src/art.h", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	defer outFile.Close()
	dto := HeaderDTO{
		NumPixels: uint(width * height),
		Pg:        pg,
		Gifs:      gifs,
	}
	t.ExecuteTemplate(outFile, "headers.tmpl", dto)
	if print {
		t.ExecuteTemplate(os.Stdout, "headers.tmpl", dto)
	}
}
