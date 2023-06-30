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

// Converts image.Image to a pair of .c/.h files.
func imgToC(name string, img image.Image) (bufC []byte, bufH []byte) {
	var (
		bufferH bytes.Buffer
		bufferC bytes.Buffer
	)

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

	headerFileName := strings.TrimSuffix(name, filepath.Ext(name)) + ".h"

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	bufferH.WriteString("// Generated with img2rgb565\n")
	bufferH.WriteString(fmt.Sprintf("// Original file name: %s\n", name))
	bufferH.WriteString(fmt.Sprintf("// Image size: %dx%d pixels (%d bytes)\n", width, height, width*height*2))
	bufferH.WriteString("\n")
	bufferH.WriteString("#ifndef " + guard + "\n")
	bufferH.WriteString("#define " + guard + "\n")
	bufferH.WriteString("\n")
	bufferH.WriteString("#include <stdint.h>\n")
	bufferH.WriteString("\n")
	bufferH.WriteString("extern uint16_t " + variable + "[];\n")
	bufferH.WriteString("\n")
	bufferH.WriteString("#endif // " + guard + "\n")

	bufferC.WriteString("// Generated with img2rgb565\n")
	bufferC.WriteString(fmt.Sprintf("// Original file name: %s\n", name))
	bufferC.WriteString(fmt.Sprintf("// Image size: %dx%d pixels (%d bytes)\n", width, height, width*height*2))
	bufferC.WriteString("\n")
	bufferC.WriteString("#include \"" + headerFileName + "\"\n")
	bufferC.WriteString("\n")
	bufferC.WriteString("uint16_t " + variable + "[] = {\n")

	w := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ { // Columns
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ { // Rows
			if w == 0 {
				bufferC.WriteString("\t")
			}

			r, g, b, _ := img.At(x, y).RGBA()

			// Scale 16-bit to 8-bit
			r /= 0x101
			g /= 0x101
			b /= 0x101

			// Convert RGB888 to RGB565
			rgb565 := uint16((b >> 3) | ((g >> 2) << 5) | ((r >> 3) << 11))

			bufferC.WriteString(fmt.Sprintf("0x%04x,", rgb565))

			w++
			if w == lineWidth {
				w = 0
				bufferC.WriteString("\n")
			} else {
				bufferC.WriteString(" ")
			}
		}
	}

	// Remove trailing ', '
	bufferC.Truncate(bufferC.Len() - 2)

	bufferC.WriteString("\n")
	bufferC.WriteString("};\n")

	return bufferC.Bytes(), bufferH.Bytes()
}

// Program entry point.
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Println("Usage: img2rgb565 IMAGE_FILE")
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

	imgC, imgH := imgToC(fileName, img)

	outNameC := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".c"
	err = os.WriteFile(outNameC, imgC, 0600)
	if err != nil {
		log.Fatal(err)
	}

	outNameH := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".h"
	err = os.WriteFile(outNameH, imgH, 0600)
	if err != nil {
		log.Fatal(err)
	}
}
