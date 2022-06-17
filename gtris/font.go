package gtris

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewFont() font.Face {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	var gfont font.Face
	gfont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 8,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	return gfont
}
