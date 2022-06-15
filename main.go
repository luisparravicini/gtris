package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 256
	screenHeight = 240
)

//go:embed images/block.png
var imgBlock []byte

type Game struct {
	blockImage *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Nearest Filter (default) VS Linear Filter")

	w, h := g.blockImage.Size()
	for x := 0; x < 10; x++ {
		for y := 0; y < 24; y++ {
			op := &ebiten.DrawImageOptions{}
			// op.GeoM.Scale(4, 4)
			op.GeoM.Translate(float64(x*w), float64(y*h))
			screen.DrawImage(g.blockImage, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(imgBlock))
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		blockImage: ebiten.NewImageFromImage(img),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Filter (Ebiten Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
