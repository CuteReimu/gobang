package main

import (
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"sync"
)

type humanPlayer struct {
	sync.Mutex
	board     [][]playerColor
	btn       [][]*walk.PushButton
	isTurn    bool
	p         point
	pNotNil   bool
	pColor    playerColor
	nextPoint chan point
}

func newHumanPlayer(color playerColor) player {
	hp := &humanPlayer{
		board:     make([][]playerColor, maxLen),
		btn:       make([][]*walk.PushButton, maxLen),
		pColor:    color,
		nextPoint: make(chan point),
	}
	for i := 0; i < maxLen; i++ {
		hp.board[i] = make([]playerColor, maxLen)
		hp.btn[i] = make([]*walk.PushButton, maxLen)
	}
	ch := make(chan bool)
	go func() {
		var mainWindow *walk.MainWindow
		mw := MainWindow{
			AssignTo: &mainWindow,
			Bounds:   Rectangle{X: 10, Y: 10, Width: 600, Height: 600},
			Font:     Font{Family: "宋体", PointSize: 18},
			Layout:   Grid{Columns: maxLen, MarginsZero: true, SpacingZero: true},
		}
		for i := 0; i < maxLen; i++ {
			for j := 0; j < maxLen; j++ {
				ii, jj := i, j
				mw.Children = append(mw.Children, PushButton{
					AssignTo: &hp.btn[ii][jj],
					MaxSize:  Size{Width: 40, Height: 40},
					Text:     "",
					OnClicked: func() {
						if hp.isTurn {
							hp.Lock()
							defer hp.Unlock()
							if hp.isTurn && hp.board[ii][jj] == 0 {
								hp.board[ii][jj] = hp.pColor
								if err := hp.btn[ii][jj].SetText(hp.pColor.getString1()); err != nil {
									log.Println(err.Error())
								}
								if hp.pNotNil {
									if err := hp.btn[hp.p.y][hp.p.x].SetText(hp.pColor.conversion().getString0()); err != nil {
										log.Println(err.Error())
									}
								}
								hp.p, hp.pNotNil = point{jj, ii}, true
								hp.isTurn = false
								hp.nextPoint <- hp.p
							}
						}
					},
				})
			}
		}
		if err := mw.Create(); err != nil {
			panic(err)
		}
		ch <- true
		code := mainWindow.Run()
		os.Exit(code)
	}()
	<-ch
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
	if err := h.btn[p.y][p.x].SetText(color.getString1()); err != nil {
		return err
	}
	if h.pNotNil {
		if err := h.btn[h.p.y][h.p.x].SetText(h.pColor.getString0()); err != nil {
			return err
		}
	}
	h.p, h.pNotNil = p, true
	return nil
}

type humanWatcher struct {
	btn     [][]*walk.PushButton
	p       point
	pNotNil bool
}

func newHumanWatcher() *humanWatcher {
	hp := &humanWatcher{
		btn: make([][]*walk.PushButton, maxLen),
	}
	for i := 0; i < maxLen; i++ {
		hp.btn[i] = make([]*walk.PushButton, maxLen)
	}
	ch := make(chan bool)
	go func() {
		var mainWindow *walk.MainWindow
		mw := MainWindow{
			AssignTo: &mainWindow,
			Bounds:   Rectangle{X: 10, Y: 10, Width: 600, Height: 600},
			Font:     Font{Family: "宋体", PointSize: 18},
			Layout:   Grid{Columns: maxLen, MarginsZero: true, SpacingZero: true},
		}
		for i := 0; i < maxLen; i++ {
			for j := 0; j < maxLen; j++ {
				ii, jj := i, j
				mw.Children = append(mw.Children, PushButton{
					AssignTo: &hp.btn[ii][jj],
					MaxSize:  Size{Width: 40, Height: 40},
					Text:     "",
				})
			}
		}
		if err := mw.Create(); err != nil {
			panic(err)
		}
		ch <- true
		code := mainWindow.Run()
		os.Exit(code)
	}()
	<-ch
	return hp
}

func (h *humanWatcher) display(color playerColor, p point) error {
	if err := h.btn[p.y][p.x].SetText(color.getString1()); err != nil {
		return err
	}
	//if brush, err := walk.NewSolidColorBrush(walk.RGB(255, 255, 0)); err != nil {
	//	return err
	//} else {
	//	h.btn[p.y][p.x].SetBackground(brush)
	//}
	if h.pNotNil {
		if err := h.btn[h.p.y][h.p.x].SetText(color.conversion().getString0()); err != nil {
			return err
		}
	}
	h.p, h.pNotNil = p, true
	return nil
}
