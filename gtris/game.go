package gtris

import (
	_ "embed"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	pieces        []*Piece
	currentPiece  *Piece
	piecePosition *Position
	gameZoneSize  Size
	gameZone      [][]*ebiten.Image
	bgBlockImage  *ebiten.Image
	input         Input
}

func (g *Game) nextPiece() {
	g.currentPiece = g.pieces[rand.Intn(len(g.pieces))]
	g.piecePosition = &Position{X: 4, Y: 0}
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
			if gameZonePos.X < 0 || gameZonePos.Y < 0 {
				continue
			}
			if gameZonePos.X >= int(g.gameZoneSize.Width) || gameZonePos.Y >= int(g.gameZoneSize.Height) {
				continue
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

func (g *Game) processInput(key ebiten.Key) {
	if key == ebiten.KeyDown {
		if g.piecePosition.Y < int(g.gameZoneSize.Height)-1 {
			g.piecePosition.Y++
		} else {
			g.transferPieceToGameZone()
			g.nextPiece()
		}
	}

	if key == ebiten.KeyLeft {
		if g.piecePosition.X > 0 {
			g.piecePosition.X--
		}
	}

	if key == ebiten.KeyRight {
		if g.piecePosition.X < int(g.gameZoneSize.Width) {
			g.piecePosition.X++
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Nearest Filter (default) VS Linear Filter")

	gameZonePos := &Position{X: 16, Y: 16}

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
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	game := &Game{
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
