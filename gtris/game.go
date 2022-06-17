package gtris

import (
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 240
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

//go:embed images/block-bg.png
var imgBlockBG []byte

type Size struct {
	Width  uint
	Height uint
}

type Game struct {
	lastTime      uint
	fallTime      uint
	score         uint
	pieces        []*Piece
	currentPiece  *Piece
	piecePosition *Position
	gameZoneSize  Size
	gameZone      [][]*ebiten.Image
	bgBlockImage  *ebiten.Image
	txtFont       font.Face
	input         Input
}

func (g *Game) nextPiece() {
	g.currentPiece = g.pieces[rand.Intn(len(g.pieces))]
	g.piecePosition = &Position{X: int(g.gameZoneSize.Width)/2 - 1, Y: 0}
}

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

func (g *Game) Update() error {
	if g.lastTime == 0 {
		g.lastTime = uint(time.Now().UnixMilli())
	}

	key := g.input.Read()
	if key != nil {
		g.processInput(*key)
	}

	return nil
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

func (g *Game) processInput(key ebiten.Key) {
	if key == ebiten.KeyDown {
		deltaPos := Position{X: 0, Y: 1}
		if g.insideGameZone(deltaPos) {
			g.piecePosition.Add(deltaPos)
		} else {
			g.transferPieceToGameZone()
			g.nextPiece()
		}
	}

	if key == ebiten.KeyLeft {
		deltaPos := Position{X: -1, Y: 0}
		if g.insideGameZone(deltaPos) {
			g.piecePosition.Add(deltaPos)
		}
	}

	if key == ebiten.KeyRight {
		deltaPos := Position{X: 1, Y: 0}
		if g.insideGameZone(deltaPos) {
			g.piecePosition.Add(deltaPos)
		}
	}
}

func (g *Game) drawScore(screen *ebiten.Image, gameZonePos *Position) {
	boardBlockWidth, _ := g.bgBlockImage.Size()
	boardWidth := int(g.gameZoneSize.Width) * boardBlockWidth
	text.Draw(screen, "SCORE", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2, color.White)
	text.Draw(screen, fmt.Sprintf("%08d", g.score), g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+8, color.White)
}

func (g *Game) Draw(screen *ebiten.Image) {
	gameZonePos := &Position{X: 16, Y: 16}

	g.drawScore(screen, gameZonePos)

	gameZone := g.gameZone
	for y, row := range gameZone {
		for x, cellImage := range row {
			if cellImage == nil {
				cellImage = g.bgBlockImage
			}

			w, h := cellImage.Size()
			screenPos := &Position{
				X: gameZonePos.X + x*w,
				Y: gameZonePos.Y + y*h,
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
			screen.DrawImage(cellImage, op)
		}
	}

	if g.currentPiece != nil {
		piece := g.currentPiece
		piece.Draw(screen, gameZonePos, g.piecePosition)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	game := &Game{
		txtFont:  NewFont(),
		input:    NewAttractModeInput(),
		fallTime: 300,
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
		gameZoneSize: Size{Width: 10, Height: 24},
		bgBlockImage: createImage(imgBlockBG),
	}

	game.gameZone = make([][]*ebiten.Image, game.gameZoneSize.Height)
	for y := range game.gameZone {
		game.gameZone[y] = make([]*ebiten.Image, game.gameZoneSize.Width)
	}

	game.nextPiece()

	return game
}
