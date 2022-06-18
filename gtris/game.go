package gtris

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	ScreenWidth  = 240
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
	dropTicks   uint
	elapsedDrop uint

	score       int
	state       GameState
	attractMode bool
	pieces      []*Piece

	nextPiece     *Piece
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

func (g *Game) Start() {
	g.state = GameStatePlaying
	g.score = 0
	g.attractMode = true
	g.input = g.inputAttractMode
	g.elapsedDrop = 0

	g.gameZone = make([][]*ebiten.Image, g.gameZoneSize.Height)
	for y := range g.gameZone {
		g.gameZone[y] = make([]*ebiten.Image, g.gameZoneSize.Width)
	}

	g.fetchNextPiece()
}

func (g *Game) StartPlay() {
	g.Start()
	g.attractMode = false
	g.input = g.inputKeyboard
}

func (g *Game) Update() error {
	g.elapsedDrop += 1

	switch g.state {
	case GameStatePlaying:
		if g.elapsedDrop > g.dropTicks {
			g.processInput(keyDown)
			g.elapsedDrop = 0
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

func (g *Game) processPiece() bool {
	g.transferPieceToGameZone()
	linesRemoved := g.checkForLines()
	g.updateScore(linesRemoved)
	g.fetchNextPiece()

	stopProcess := false
	deltaPos := Position{}
	if !g.insideGameZone(deltaPos) {
		g.state = GameStateGameOver
		stopProcess = true
	}

	return stopProcess
}

func (g *Game) processInput(key ebiten.Key) {
	if key == ebiten.KeyDown {
		deltaPos := Position{X: 0, Y: 1}
		if g.insideGameZone(deltaPos) {
			g.piecePosition.Add(deltaPos)
		} else {
			stopProcess := g.processPiece()
			if stopProcess {
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

	if key == ebiten.KeyUp {
		newPiece := g.rotatePiece()
		if g.pieceInsideGameZone(newPiece, *g.piecePosition) {
			g.currentPiece = newPiece
		}
	}
}

func (g *Game) drawText(screen *ebiten.Image, gameZonePos *Position) {
	boardBlockWidth, _ := g.bgBlockImage.Size()
	boardWidth := int(g.gameZoneSize.Width) * boardBlockWidth
	text.Draw(screen, "SCORE", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2, color.White)
	text.Draw(screen, fmt.Sprintf("%08d", g.score), g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+8, color.White)

	if g.state == GameStateGameOver {
		dy := 122
		text.Draw(screen, "GAME OVER", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy, color.White)
		text.Draw(screen, "space to start", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy+8, color.White)
	}

	if g.attractMode {
		dy := 148
		text.Draw(screen, "press space", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy, color.White)
		text.Draw(screen, "  to play", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy+8, color.White)
	}

	dy := 48
	text.Draw(screen, "NEXT", g.txtFont, boardWidth+gameZonePos.X*2, gameZonePos.Y*2+dy, color.White)
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
		g.currentPiece.Draw(screen, gameZonePos, g.piecePosition)
	}

	if g.nextPiece != nil {
		nextPos := &Position{X: int(math.Round(ScreenWidth * 0.5)), Y: int(math.Round(ScreenHeight * .37))}
		g.nextPiece.Draw(screen, nextPos, &Position{})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	ebiten.SetMaxTPS(18)

	game := &Game{
		txtFont:          NewFont(),
		inputAttractMode: NewAttractModeInput(),
		inputKeyboard:    &KeyboardInput{},
		dropTicks:        4,
		pieces:           allPieces,
		gameZoneSize:     Size{Width: 10, Height: 24},
		bgBlockImage:     createImage(imgBlockBG),
	}

	game.Start()

	return game
}
