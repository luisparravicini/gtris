package gtris

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) transferPieceToGameZone() {
	piece := g.currentPiece
	piecePos := g.piecePosition
	for dy, row := range piece.Blocks {
		for dx, value := range row {
			if value != pieceBlockMarker {
				continue
			}

			gameZonePos := &Position{
				X: piecePos.X + dx,
				Y: piecePos.Y + dy,
			}

			g.gameZone[gameZonePos.Y][gameZonePos.X] = piece.Image
		}
	}
}

func (g *Game) insideGameZone(deltaPos Position) bool {
	piecePos := *g.piecePosition
	piecePos.Add(deltaPos)

	for dy, row := range g.currentPiece.Blocks {
		for dx, value := range row {
			if value == pieceBlockMarker {
				screenPos := &Position{
					X: piecePos.X + dx,
					Y: piecePos.Y + dy,
				}

				if screenPos.X < 0 || screenPos.X >= int(g.gameZoneSize.Width) {
					return false
				}
				if screenPos.Y < 0 || screenPos.Y >= int(g.gameZoneSize.Height) {
					return false
				}

				if g.gameZone[screenPos.Y][screenPos.X] != nil {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) nextPiece() {
	g.currentPiece = g.pieces[rand.Intn(len(g.pieces))]
	g.piecePosition = &Position{X: int(g.gameZoneSize.Width)/2 - 1, Y: 0}
}

func (g *Game) checkForLines() int {
	lines := []int{}
	for y, row := range g.gameZone {
		var full = true
		for _, cellImage := range row {
			if cellImage == nil {
				full = false
				break
			}
		}
		if full {
			lines = append(lines, y)
		}
	}

	for _, y := range lines {
		emptyRow := [][]*ebiten.Image{
			make([]*ebiten.Image, g.gameZoneSize.Width),
		}
		g.gameZone = append(append(emptyRow, g.gameZone[0:y]...), g.gameZone[(y+1):]...)
	}

	return len(lines)
}
