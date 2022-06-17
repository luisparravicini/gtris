package gtris

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const pieceBlockMarker = 1

type Piece struct {
	Blocks [][]int
	Image  *ebiten.Image
}

func NewPiece(blocks [][]int, imgData []byte) *Piece {
	return &Piece{
		Blocks: blocks,
		Image:  createImage(imgData),
	}
}

func (p *Piece) Draw(screen *ebiten.Image, gameZonePos *Position, piecePos *Position) {
	w, h := p.Image.Size()

	for dy, row := range p.Blocks {
		for dx, value := range row {
			if value == pieceBlockMarker {
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
