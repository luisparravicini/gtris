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

//go:embed images/block-a.png
var imgBlockA []byte

//go:embed images/block-b.png
var imgBlockB []byte

//go:embed images/block-c.png
var imgBlockC []byte

//go:embed images/block-d.png
var imgBlockD []byte

//go:embed images/block-e.png
var imgBlockE []byte

//go:embed images/block-f.png
var imgBlockF []byte

//go:embed images/block-g.png
var imgBlockG []byte

type Game struct {
	pieces []*Piece
}

func (g *Game) Update() error {
	return nil
}

type Position struct {
	X int
	Y int
}

type Piece struct {
	Blocks []string
	Image  *ebiten.Image
}

func NewPiece(blocks []string, imgData []byte) *Piece {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}

	return &Piece{
		Blocks: blocks,
		Image:  ebiten.NewImageFromImage(img),
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Nearest Filter (default) VS Linear Filter")

	piece := g.pieces[2]

	w, h := piece.Image.Size()
	piecePos := &Position{}
	gameZonePos := &Position{X: 16, Y: 16}

	for dy, row := range piece.Blocks {
		for dx, value := range row {
			if value == 'X' {
				op := &ebiten.DrawImageOptions{}
				// op.ColorM.Translate(1, 0, 0, 1)
				screenPos := &Position{
					X: gameZonePos.X + (piecePos.X+dx)*w,
					Y: gameZonePos.Y + (piecePos.Y+dy)*h,
				}
				op.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
				screen.DrawImage(piece.Image, op)
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
	game := &Game{
		pieces: []*Piece{
			NewPiece([]string{
				"XXXX",
				"    ",
			}, imgBlockA),
			NewPiece([]string{
				"X   ",
				"XXXX",
			}, imgBlockB),
			NewPiece([]string{
				"   X",
				"XXXX",
			}, imgBlockC),

			NewPiece([]string{
				"XX  ",
				"XX  ",
			}, imgBlockD),

			NewPiece([]string{
				" XX ",
				"XX  ",
			}, imgBlockE),
			NewPiece([]string{
				" X  ",
				"XXX ",
			}, imgBlockF),
			NewPiece([]string{
				"XX  ",
				" XX ",
			}, imgBlockG),
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("gtris")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
