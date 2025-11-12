package assets

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/Mx437_IBM_VGA_8x16.ttf
var MainFont []byte

func LoadFont(font []byte, size float64) text.Face {
	fnt, err := opentype.Parse(font)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	return text.NewGoXFace(face)
}
