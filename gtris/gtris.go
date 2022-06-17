package gtris

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X int
	Y int
}

func (p *Position) Add(other Position) {
	p.X += other.X
	p.Y += other.Y
}

func createImage(imgData []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}
