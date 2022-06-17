package gtris

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const keyDown = ebiten.KeyDown

var inputKeys = []ebiten.Key{keyDown, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeySpace}

type Input interface {
	Read() *ebiten.Key
}

type KeyboardInput struct {
}

func (*KeyboardInput) Read() *ebiten.Key {
	for _, key := range inputKeys {
		if inpututil.IsKeyJustPressed(key) {
			return &key
		}
	}

	return nil
}

type AttractModeInput struct {
	keyPressed chan ebiten.Key
}

func (input *AttractModeInput) Read() *ebiten.Key {
	select {
	case key := <-input.keyPressed:
		return &key
	default:
		return nil
	}
}

func NewAttractModeInput() *AttractModeInput {
	input := &AttractModeInput{
		keyPressed: make(chan ebiten.Key),
	}
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				key := keyDown
				if rand.Float32() < 0.5 {
					key = inputKeys[rand.Intn(len(inputKeys))]
				}
				input.keyPressed <- key
			}
		}
	}()

	return input
}
