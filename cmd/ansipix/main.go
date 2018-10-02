package main

// ------------------------------------------------------------------------- //
// Imports //
// ------------------------------------------------------------------------- //

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"image/color"
	"image/draw"
	_ "image/gif"  // initialize decoder
	_ "image/jpeg" // initialize decoder
	_ "image/png"  // initialize decoder

	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
	_ "golang.org/x/image/bmp"  // initialize decoder
	_ "golang.org/x/image/tiff" // initialize decoder
	_ "golang.org/x/image/webp" // initialize decoder
)

// ------------------------------------------------------------------------- //
// Constants //
// ------------------------------------------------------------------------- //

const (
	BlockSizeX int = 4
	BlockSizeY int = 8
)

// ------------------------------------------------------------------------- //
// Main Entrypoint //
// ------------------------------------------------------------------------- //

func main() {
	//gopath := os.Getenv("GOPATH")
	path := filepath.FromSlash(os.Args[1])
	fmt.Println(path)

	// black background
	bg := color.Gray16{0}

	imgOriginal, err := loadImageFromFile(path)
	if err != nil {
		panic(err)
	}

	imgScaled := scaleImage(imgOriginal)
	ansiImage := createAnsiImage(imgScaled, bg)

	clearTerminal()
	drawAnsiImage(ansiImage)
}

// ------------------------------------------------------------------------- //
// Types //
// ------------------------------------------------------------------------- //

type ansiPixel struct {
	Brightness uint8
	R, G, B    uint8
	source     *ansiImage
}

type ansiImage struct {
	w, h   int
	bgR    uint8
	bgG    uint8
	bgB    uint8
	pixels [][]*ansiPixel
}

// ------------------------------------------------------------------------- //
// Package Methods //
// ------------------------------------------------------------------------- //

func isTerminal() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

func getTerminalSize() (width, height int, err error) {
	// VT100 terminal size
	width = 80
	height = 24

	if isTerminal() {
		fd := os.Stdout.Fd()
		width, height, err = terminal.GetSize(int(fd))
	}

	if width > 200 {
		width = 120
	}

	if height < 70 || height > 120 {
		height = 100
	}

	return
}

// ClearTerminal clears current terminal buffer using ANSI escape code.
// (Nice info for ANSI escape codes - http://unix.stackexchange.com/questions/124762/how-does-clear-command-work)
func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func loadImageFromFile(name string) (image.Image, error) {
	reader, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return loadImageFromReader(reader)
}

func loadImageFromReader(reader io.Reader) (image.Image, error) {
	//  (JPEG, PNG, GIF, BMP, TIFF)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func scaleImage(img image.Image) image.Image {
	sfy, sfx := BlockSizeY, BlockSizeX // 8x4 --> with dithering
	tx, ty, err := getTerminalSize()
	if err != nil {
		panic(err)
	}

	scaledImage := imaging.Resize(img, sfy*ty, sfx*tx, imaging.Lanczos)
	return scaledImage
}

func createAnsiImage(img image.Image, bg color.Color) *ansiImage {
	bounds := img.Bounds()

	yMax, xMax := bounds.Max.Y, bounds.Max.X
	yMax = yMax / BlockSizeY // always sets 1 ANSIPixel block...
	xMax = xMax / BlockSizeX // per 8x4 real pixels --> with dithering

	rgbaOut := composeImage(img, bg)

	ansiimage := newAnsiImage(xMax, yMax, bg)

	newAnsiPixels(ansiimage, rgbaOut)

	return ansiimage
}

func composeImage(img image.Image, bg color.Color) *image.RGBA {
	// http://stackoverflow.com/questions/36595687/transparent-pixel-color-go-lang-image
	var rgbaOut *image.RGBA
	bounds := img.Bounds()

	if _, _, _, a := bg.RGBA(); a >= 0xffff {
		rgbaOut = image.NewRGBA(bounds)
		draw.Draw(rgbaOut, bounds, image.NewUniform(bg), image.ZP, draw.Src)
		draw.Draw(rgbaOut, bounds, img, image.ZP, draw.Over)
	} else {
		if v, ok := img.(*image.RGBA); ok {
			rgbaOut = v
		} else {
			rgbaOut = image.NewRGBA(bounds)
			draw.Draw(rgbaOut, bounds, img, image.ZP, draw.Src)
		}
	}

	return rgbaOut
}

func newAnsiImage(w, h int, bg color.Color) *ansiImage {
	// create instance
	r, g, b, _ := bg.RGBA()
	ansimage := &ansiImage{
		w:      w,
		h:      h,
		bgR:    uint8(r),
		bgG:    uint8(g),
		bgB:    uint8(b),
		pixels: nil,
	}

	// initialize pixels to (0,0,0,0)
	ansimage.pixels = func() [][]*ansiPixel {
		v := make([][]*ansiPixel, h)
		for y := 0; y < h; y++ {
			v[y] = make([]*ansiPixel, w)
			for x := 0; x < w; x++ {
				v[y][x] = &ansiPixel{
					R:          0,
					G:          0,
					B:          0,
					Brightness: 0,
					source:     ansimage,
				}
			}
		}
		return v
	}()

	return ansimage
}

func newAnsiPixels(img *ansiImage, rgba *image.RGBA) {
	// calculate brightness

	pixelCount := BlockSizeY * BlockSizeX
	for y := 0; y < img.h; y++ {
		for x := 0; x < img.w; x++ {

			var sumR, sumG, sumB, sumBri float64
			for dy := 0; dy < BlockSizeY; dy++ {
				py := BlockSizeY*y + dy

				for dx := 0; dx < BlockSizeX; dx++ {
					px := BlockSizeX*x + dx

					pixel := rgba.At(px, py)
					color, _ := colorful.MakeColor(pixel)
					_, _, v := color.Hsv()
					sumR += color.R
					sumG += color.G
					sumB += color.B
					sumBri += v
				}
			}

			r := uint8(sumR/float64(pixelCount)*255.0 + 0.5)
			g := uint8(sumG/float64(pixelCount)*255.0 + 0.5)
			b := uint8(sumB/float64(pixelCount)*255.0 + 0.5)
			brightness := uint8(sumBri/float64(pixelCount)*255.0 + 0.5)

			img.pixels[y][x].R = r
			img.pixels[y][x].G = g
			img.pixels[y][x].B = b
			img.pixels[y][x].Brightness = brightness
		}
	}
}

func drawAnsiImage(img *ansiImage) {
	type renderData struct {
		row    int
		render string
	}

	maxprocs := runtime.NumCPU()
	rows := make([]string, img.h)

	for y := 0; y < img.h; y += maxprocs {
		ch := make(chan renderData, maxprocs)
		for n, r := 0, y; (n <= maxprocs) && (r+1 < img.h); n, r = n+1, y+n+1 {
			go func(y int) {
				var str string
				for x := 0; x < img.w; x++ {
					pixel := img.pixels[y][x]
					str += drawAnsiPixel(pixel)
				}
				str += "\033[0m\n" // reset ansi style
				ch <- renderData{row: y, render: str}
			}(r)
		}
		for n, r := 0, y; (n <= maxprocs) && (r+1 < img.h); n, r = n+1, y+n+1 {
			data := <-ch
			rows[data.row] = data.render
		}
	}

	fmt.Println(strings.Join(rows, ""))
}

func drawAnsiPixel(px *ansiPixel) string {
	block := " "
	switch bri := px.Brightness; {
	case bri > 230:
		block = "#"
	case bri > 207:
		block = "&"
	case bri > 184:
		block = "$"
	case bri > 161:
		block = "X"
	case bri > 138:
		block = "x"
	case bri > 115:
		block = "="
	case bri > 92:
		block = "+"
	case bri > 69:
		block = ";"
	case bri > 46:
		block = ":"
	case bri > 23:
		block = "."
	}

	return fmt.Sprintf(
		"\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm%s",
		px.source.bgR, px.source.bgG, px.source.bgB,
		px.R, px.G, px.B,
		block,
	)
}
