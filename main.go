package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"path/filepath"
	"strings"
	"time"

	"image"
	"log"
	"os"

	"github.com/nfnt/resize"
)

var inputFile string
var treshold int
var imgWidth int
var inverted bool

func main() {
	flag.StringVar(&inputFile, "i", "", "Input file")

	flag.IntVar(&treshold, "t", 0, "Treshold of density")
	flag.IntVar(&imgWidth, "w", 48, "Width of output image")
	flag.BoolVar(&inverted, "invert", false, "Invert colors")
	flag.Parse()
	if len(inputFile) == 0 {
		log.Fatal("Input file is not specified")
	}
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var img image.Image
	ext := filepath.Ext(inputFile)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	case ".png":
		img, err = png.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	}
	img = resize.Resize(uint(imgWidth), 0, img, resize.Lanczos3)
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	resultName := fmt.Sprintf("%s_ascii_%d", strings.Split(f.Name(), ".")[0], +time.Now().Unix())
	imgString, err := os.Create(resultName + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	var density = " _.,-=+:;cba!?0123456789$W#@Ñ"
	if inverted {
		density = "Ñ@#W$9876543210?!abc;:+=-,._ "
	}
	for y := 0; y < treshold; y++ {
		if inverted {
			density = " " + density
			continue
		}
		density = density + " "
	}
	defer imgString.Close()
	for i := 0; i < height; i++ {
		text := ""
		for j := 0; j < width; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			gray := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			at := math.Floor(float64(int(gray)*(len(density)-1)) / 255)
			text = text + string(density[int(at)]) + ""
		}
		io.WriteString(imgString, text+"\n")
	}
	fmt.Printf("success, check your file %s\n", resultName+".txt")

}
