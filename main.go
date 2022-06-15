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

type Position struct {
	X int
	Y int
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Nearest Filter (default) VS Linear Filter")

	// piece := []string{
	// 	"XXXX",
	// 	"    ",
	// }
	// piece := []string{
	// 	"X   ",
	// 	"XXXX",
	// }
	// piece := []string{
	// 	"   X",
	// 	"XXXX",
	// }
	// piece := []string{
	// 	"XX  ",
	// 	"XX  ",
	// }
	// piece := []string{
	// 	" XX ",
	// 	"XX  ",
	// }
	piece := []string{
		" X  ",
		"XXX ",
	}
	// piece := []string{
	// 	"XX  ",
	// 	" XX ",
	// }
	w, h := g.blockImage.Size()
	piecePos := &Position{}
	gameZonePos := &Position{X: 16, Y: 16}

	for dy, row := range piece {
		for dx, value := range row {
			if value == 'X' {
				op := &ebiten.DrawImageOptions{}
				screenPos := &Position{
					X: gameZonePos.X + (piecePos.X+dx)*w,
					Y: gameZonePos.Y + (piecePos.Y+dy)*h,
				}
				op.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
				screen.DrawImage(g.blockImage, op)
			}
		}
	}

	// for x := 0; x < 10; x++ {
	// 	for y := 0; y < 24; y++ {
	// 		op := &ebiten.DrawImageOptions{}
	// 		// op.GeoM.Scale(4, 4)
	// 		op.GeoM.Translate(float64(x*w), float64(y*h))
	// 		screen.DrawImage(g.blockImage, op)
	// 	}
	// }
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
	ebiten.SetWindowTitle("gtris")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
