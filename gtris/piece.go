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

func CreateImage(imgData []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

type Piece struct {
	Blocks []string
	Image  *ebiten.Image
}

func NewPiece(blocks []string, imgData []byte) *Piece {
	return &Piece{
		Blocks: blocks,
		Image:  CreateImage(imgData),
	}
}

func (p *Piece) Draw(screen *ebiten.Image, gameZonePos *Position, piecePos *Position) {
	w, h := p.Image.Size()

	for dy, row := range p.Blocks {
		for dx, value := range row {
			if value == 'X' {
				screenPos := &Position{
					X: gameZonePos.X + (piecePos.X+dx)*w,
					Y: gameZonePos.Y + (piecePos.Y+dy)*h,
				}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
				screen.DrawImage(p.Image, op)
			}
		}
	}
}
