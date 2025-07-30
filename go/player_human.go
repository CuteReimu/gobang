package main

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math"
	"sync"
)

type humanPlayer struct {
	sync.Mutex
	board     [][]playerColor
	isTurn    bool
	p         point
	pColor    playerColor
	nextPoint chan point
}

func (h *humanPlayer) Update() error {
	if h.isTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		x -= 17
		y -= 17
		if x-x/35*35-18 < 10 && y-y/35*35-18 < 10 {
			x /= 35
			y /= 35
			if h.board[x][y] == colorEmpty {
				h.isTurn = false
				h.p = point{y, x}
				h.nextPoint <- h.p
				h.board[x][y] = h.pColor
			}
		}
	}
	return nil
}

func (h *humanPlayer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	img0 := ebiten.NewImage(35*(maxLen+1), 35*(maxLen+1))
	img := ebiten.NewImage(35*(maxLen-1), 1)
	img.Fill(color.Black)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(35, 35)
	for range maxLen {
		img0.DrawImage(img, opt)
		opt.GeoM.Translate(0, 35)
	}
	opt = &ebiten.DrawImageOptions{}
	screen.DrawImage(img0, opt)
	opt.GeoM.Translate(-35-maxLen/2*35, -35-maxLen/2*35)
	opt.GeoM.Rotate(math.Pi / 2)
	opt.GeoM.Translate(35+maxLen/2*35, 35+maxLen/2*35)
	screen.DrawImage(img0, opt)
	for i, row := range h.board {
		for j, color := range row {
			if color != colorEmpty {
				img := pieceBlack
				if color == colorWhite {
					img = pieceWhite
				}
				if i == h.p.y && j == h.p.x {
					img = pieceBlack2
					if color == colorWhite {
						img = pieceWhite2
					}
				}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(18+35*i), float64(18+35*j))
				screen.DrawImage(img, op)
			}
		}
	}
}

func (h *humanPlayer) Layout(int, int) (screenWidth int, screenHeight int) {
	return 35 * (maxLen + 1), 35 * (maxLen + 1)
}

func newHumanPlayer(color playerColor) *humanPlayer {
	hp := &humanPlayer{
		board:     make([][]playerColor, maxLen),
		pColor:    color,
		nextPoint: make(chan point),
	}
	for i := 0; i < maxLen; i++ {
		hp.board[i] = make([]playerColor, maxLen)
	}
	return hp
}

func (h *humanPlayer) color() playerColor {
	return h.pColor
}

func (h *humanPlayer) play() (point, error) {
	h.isTurn = true
	return <-h.nextPoint, nil
}

func (h *humanPlayer) display(p point) error {
	if h.board[p.y][p.x] != 0 {
		return errors.New(fmt.Sprintf("illegal argument: %s%s\n", p, h.board[p.y][p.x]))
	}
	color := h.pColor.conversion()
	h.board[p.y][p.x] = color
	h.p = p
	return nil
}

type humanWatcher struct {
	board [][]playerColor
	p     point
}

func newHumanWatcher() *humanWatcher {
	hp := &humanWatcher{
		board: make([][]playerColor, maxLen),
	}
	for i := 0; i < maxLen; i++ {
		hp.board[i] = make([]playerColor, maxLen)
	}
	return hp
}

func (h *humanWatcher) Update() error {
	return nil
}

func (h *humanWatcher) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	img0 := ebiten.NewImage(35*(maxLen+1), 35*(maxLen+1))
	img := ebiten.NewImage(35*(maxLen-1), 1)
	img.Fill(color.Black)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(35, 35)
	for range maxLen {
		img0.DrawImage(img, opt)
		opt.GeoM.Translate(0, 35)
	}
	opt = &ebiten.DrawImageOptions{}
	screen.DrawImage(img0, opt)
	opt.GeoM.Translate(-35-maxLen/2*35, -35-maxLen/2*35)
	opt.GeoM.Rotate(math.Pi / 2)
	opt.GeoM.Translate(35+maxLen/2*35, 35+maxLen/2*35)
	screen.DrawImage(img0, opt)
	for i, row := range h.board {
		for j, color := range row {
			if color != colorEmpty {
				img := pieceBlack
				if color == colorWhite {
					img = pieceWhite
				}
				if i == h.p.y && j == h.p.x {
					img = pieceBlack2
					if color == colorWhite {
						img = pieceWhite2
					}
				}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(18+35*i), float64(18+35*j))
				screen.DrawImage(img, op)
			}
		}
	}
}

func (h *humanWatcher) Layout(int, int) (screenWidth int, screenHeight int) {
	return 35 * (maxLen + 1), 35 * (maxLen + 1)
}

func (h *humanWatcher) display(color playerColor, p point) error {
	h.board[p.y][p.x] = color
	h.p = p
	return nil
}

var pieceWhite = ebiten.NewImage(33, 33)
var pieceBlack = ebiten.NewImage(33, 33)
var pieceWhite2 = ebiten.NewImage(33, 33)
var pieceBlack2 = ebiten.NewImage(33, 33)

func init() {
	pieceWhite.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	pieceBlack.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	pieceWhite2.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	pieceBlack2.Fill(color.RGBA{R: 0xee, G: 0xd2, B: 0x5c, A: 0xff})
	for i := range 33 {
		for j := range 33 {
			diff := math.Sqrt(float64((i-16)*(i-16) + (j-16)*(j-16)))
			if diff <= 16 {
				pieceBlack.Set(i, j, color.Black)
				pieceBlack2.Set(i, j, color.Black)
				if diff > 14.5 {
					pieceWhite.Set(i, j, color.Black)
					pieceWhite2.Set(i, j, color.Black)
				} else if diff > 13 {
					pieceWhite.Set(i, j, color.White)
					pieceWhite2.Set(i, j, color.White)
					pieceBlack2.Set(i, j, color.White)
				} else if diff > 11.5 {
					pieceWhite.Set(i, j, color.White)
					pieceWhite2.Set(i, j, color.Black)
				} else {
					pieceWhite.Set(i, j, color.White)
					pieceWhite2.Set(i, j, color.White)
				}
			}
		}
	}
}
