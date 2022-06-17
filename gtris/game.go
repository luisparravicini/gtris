package gtris

import (
	"fmt"
	"image/color"
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

type Size struct {
	Width  uint
	Height uint
}

type GameState int

const (
	GameStateGameOver GameState = iota
	GameStatePlaying
)

type Game struct {
	lastTime time.Time
	fallTime uint
	elapsed  uint

	score       int
	state       GameState
	attractMode bool
	pieces      []*Piece

	currentPiece  *Piece
	piecePosition *Position

	gameZoneSize Size
	gameZone     [][]*ebiten.Image
	bgBlockImage *ebiten.Image

	txtFont font.Face

	input            Input
	inputAttractMode Input
	inputKeyboard    Input
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

func (g *Game) Start() {
	g.state = GameStatePlaying
	g.score = 0
	g.attractMode = true
	g.input = g.inputAttractMode
	g.lastTime = time.Now()
	g.elapsed = 0

	g.gameZone = make([][]*ebiten.Image, g.gameZoneSize.Height)
	for y := range g.gameZone {
		g.gameZone[y] = make([]*ebiten.Image, g.gameZoneSize.Width)
	}

	g.nextPiece()
}

func (g *Game) StartPlay() {
	g.Start()
	g.attractMode = false
	g.input = g.inputKeyboard
}

func (g *Game) Update() error {
	g.elapsed += uint(time.Since(g.lastTime).Milliseconds())
	g.lastTime = time.Now()

	switch g.state {
	case GameStatePlaying:
		if g.elapsed > g.fallTime {
			g.processInput(keyDown)
			g.elapsed = 0
			return nil
		}

		key := g.input.Read()
		if key != nil {
			g.processInput(*key)
		}

		if g.attractMode && g.inputKeyboard.IsSpacePressed() {
			g.StartPlay()
		}
	case GameStateGameOver:
		key := g.input.Read()
		if key != nil {
			g.Start()
		}
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
			linesRemoved := g.checkForLines()
			g.updateScore(linesRemoved)
			g.nextPiece()

			deltaPos := Position{}
			if !g.insideGameZone(deltaPos) {
				g.state = GameStateGameOver
				return
			}
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

func (g *Game) drawText(screen *ebiten.Image, gameZonePos *Position) {
	boardBlockWidth, _ := g.bgBlockImage.Size()
	boardWidth := int(g.gameZoneSize.Width) * boardBlockWidth
	text.Draw(screen, "SCORE", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2, color.White)
	text.Draw(screen, fmt.Sprintf("%08d", g.score), g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+8, color.White)

	if g.state == GameStateGameOver {
		dy := 32
		text.Draw(screen, "GAME OVER", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy, color.White)
		text.Draw(screen, "space to start", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy+8, color.White)
	}

	if g.attractMode {
		dy := 96
		text.Draw(screen, "press space", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy, color.White)
		text.Draw(screen, "  to play", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy+8, color.White)
	}
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

func (g *Game) updateScore(lines int) {
	perLineScore := 100
	g.score += lines * perLineScore
	if lines > 1 {
		bonus := perLineScore / 2
		for i := 0; i < int(lines); i++ {
			g.score += bonus
			bonus *= 2
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	gameZonePos := &Position{X: 16, Y: 16}

	g.drawText(screen, gameZonePos)

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
		txtFont:          NewFont(),
		inputAttractMode: NewAttractModeInput(),
		inputKeyboard:    &KeyboardInput{},
		fallTime:         700,
		pieces:           allPieces,
		gameZoneSize:     Size{Width: 10, Height: 24},
		bgBlockImage:     createImage(imgBlockBG),
	}

	game.Start()

	return game
}
