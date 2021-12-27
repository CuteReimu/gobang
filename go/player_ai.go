package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
)

type AIPlayer struct {
	board   [][]playerColor
	pColor  playerColor
	me      [][]float64
	them    [][]float64
	notMe   [][]float64
	notThem [][]float64
}

func newAIPlayer(color playerColor) *AIPlayer {
	p := &AIPlayer{
		board:  make([][]playerColor, maxLen),
		pColor: color,
	}
	for i := 0; i < len(p.board); i++ {
		p.board[i] = make([]playerColor, maxLen)
	}
	p.load()
	return p
}

func (A *AIPlayer) copy(color playerColor) *AIPlayer {
	p := &AIPlayer{
		board:   make([][]playerColor, maxLen),
		pColor:  color,
		me:      make([][]float64, maxLen),
		them:    make([][]float64, maxLen),
		notMe:   make([][]float64, maxLen),
		notThem: make([][]float64, maxLen),
	}
	for i := 0; i < len(p.board); i++ {
		p.board[i] = make([]playerColor, maxLen)
		p.me[i] = make([]float64, maxLen)
		copy(p.me[i], A.me[i])
		p.them[i] = make([]float64, maxLen)
		copy(p.them[i], A.them[i])
		p.notMe[i] = make([]float64, maxLen)
		copy(p.notMe[i], A.notMe[i])
		p.notThem[i] = make([]float64, maxLen)
		copy(p.notThem[i], A.notThem[i])
	}
	return p
}

func (A *AIPlayer) color() playerColor {
	return A.pColor
}

func (A *AIPlayer) play() (point, error) {
	maxP := point{-1, -1}
	maxV := float64(-1000000000000000000000000000000)
	f := func(p point) float64 {
		var result float64
		for x := 0; x < maxLen; x++ {
			for y := 0; y < maxLen; y++ {
				p2 := point{p.x + x - maxLen/2, p.y + y - maxLen/2}
				if !p2.checkRange() {
					result += A.notThem[y][x]
					result += A.notMe[y][x]
				} else {
					if A.board[p2.y][p2.x] == A.pColor {
						result += A.me[y][x]
						result += A.notThem[y][x]
					} else if A.board[p2.y][p2.x] != 0 {
						result += A.them[y][x]
						result += A.notMe[y][x]
					}
				}
			}
		}
		return result
	}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			if A.board[i][j] == colorEmpty {
				v := f(point{j, i})
				if v > maxV {
					maxV = v
					maxP = point{j, i}
				}
			}
		}
	}
	A.board[maxP.y][maxP.x] = A.pColor
	return maxP, nil
}

func (A *AIPlayer) display(p point) error {
	if A.board[p.y][p.x] != 0 {
		return fmt.Errorf("illegal argument: %s%s\n", p, A.board[p.y][p.x])
	}
	A.board[p.y][p.x] = A.pColor.conversion()
	return nil
}

func (A *AIPlayer) rand() {
	f := func(b [][]float64) {
		for i := 0; i < len(b); i++ {
			for j := 0; j < len(b[i]); j++ {
				b[i][j] += rand.NormFloat64()
			}
		}
	}
	f(A.me)
	f(A.them)
	f(A.notMe)
	f(A.notThem)
}

func (A *AIPlayer) load() {
	A.me = make([][]float64, maxLen)
	A.them = make([][]float64, maxLen)
	A.notMe = make([][]float64, maxLen)
	A.notThem = make([][]float64, maxLen)
	for i := 0; i < maxLen; i++ {
		A.me[i] = make([]float64, maxLen)
		A.them[i] = make([]float64, maxLen)
		A.notMe[i] = make([]float64, maxLen)
		A.notThem[i] = make([]float64, maxLen)
	}
	buf, err := ioutil.ReadFile("data.dat")
	if err != nil {
		log.Println(err.Error())
		A.rand()
		val := float64(0)
		for i := 0; i < maxLen; i++ {
			for j := 0; j < maxLen; j++ {
				val += A.me[i][j] + A.them[i][j] + A.notMe[i][j] + A.notThem[i][j]
			}
		}
		fmt.Println(val)
		return
	}
	f := func(b [][]float64, buf []byte) {
		for i := 0; i < maxLen; i++ {
			for j := 0; j < maxLen; j++ {
				u := binary.LittleEndian.Uint64(buf[i*j*8 : i*j*8+8])
				b[i][j] = math.Float64frombits(u)
			}
		}
	}
	f(A.me, buf[:maxLen*maxLen*8])
	f(A.them, buf[maxLen*maxLen*8:maxLen*maxLen*16])
	f(A.notMe, buf[maxLen*maxLen*16:maxLen*maxLen*24])
	f(A.notThem, buf[maxLen*maxLen*24:])
}

func (A *AIPlayer) save() {
	f := func(b [][]float64) []byte {
		var ret []byte
		for i := 0; i < len(b); i++ {
			for j := 0; j < len(b[i]); j++ {
				u := math.Float64bits(b[i][j])
				buf := make([]byte, 8)
				binary.LittleEndian.PutUint64(buf, u)
				ret = append(ret, buf...)
			}
		}
		return ret
	}
	b1 := append(f(A.me), f(A.them)...)
	b2 := append(b1, f(A.notMe)...)
	b3 := append(b2, f(A.notThem)...)
	err := ioutil.WriteFile("data.dat", b3, 0644)
	if err != nil {
		log.Println(err.Error())
	}
}
