package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
)

const (
	lineWidth = 16 // pixels
)

// Converts image.Image to C header file.
func imgToC(name string, img image.Image) []byte {
	var buf bytes.Buffer

	guard := "_"
	for _, c := range strings.ToUpper(strings.TrimSuffix(name, filepath.Ext(name))) {
		if c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			guard += string(c)
			continue
		}
		guard += "_"
	}
	variable := "img" + strings.ToLower(guard)
	guard += "_H"

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	buf.WriteString("// Generated with img2rgb565\n")
	buf.WriteString(fmt.Sprintf("// Original file name: %s\n", name))
	buf.WriteString(fmt.Sprintf("// Image size: %dx%d pixels (%d bytes)\n", width, height, width*height*2))
	buf.WriteString("\n")
	buf.WriteString("#ifndef " + guard + "\n")
	buf.WriteString("#define " + guard + "\n")
	buf.WriteString("\n")
	buf.WriteString("#include <stdint.h>\n")
	buf.WriteString("\n")
	buf.WriteString("uint16_t " + variable + "[] = {\n")

	w := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ { // Columns
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ { // Rows
			if w == 0 {
				buf.WriteString("\t")
			}

			r, g, b, _ := img.At(x, y).RGBA()

			// Scale 16-bit to 8-bit
			r /= 0x101
			g /= 0x101
			b /= 0x101

			// Convert RGB888 to RGB565
			rgb565 := uint16((b >> 3) | ((g >> 2) << 5) | ((r >> 3) << 11))

			buf.WriteString(fmt.Sprintf("0x%04x,", rgb565))

			w++
			if w == lineWidth {
				w = 0
				buf.WriteString("\n")
			} else {
				buf.WriteString(" ")
			}
		}
	}

	// Remove trailing ', '
	buf.Truncate(buf.Len() - 2)

	buf.WriteString("\n")
	buf.WriteString("};\n")
	buf.WriteString("\n")
	buf.WriteString("#endif // " + guard + "\n")

	return buf.Bytes()
}

// Program entry point.
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Println("Usage: img2rgb565 FILE [OUTPUT]")
		log.Println("Supported image formats: bmp, png, jpeg, gif")
		return
	}

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(contents))
	if err != nil {
		log.Fatal(err)
	}

	fileName := filepath.Base(os.Args[1])

	imgC := imgToC(fileName, img)

	outName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".h"
	if len(os.Args) > 2 {
		outName = os.Args[2]
	}

	err = os.WriteFile(outName, imgC, 0600)
	if err != nil {
		log.Fatal(err)
	}
}
