package gtris

import (
	_ "embed"
	_ "image/png"
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

var allPieces = []*Piece{
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
}
