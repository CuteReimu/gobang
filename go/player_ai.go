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
	board  [][]playerColor
	pColor playerColor
	cache  []float64
}

func newAIPlayer() *AIPlayer {
	p := &AIPlayer{
		board:  make([][]playerColor, maxLen),
		pColor: colorBlack,
		cache:  make([]float64, 0x10000),
	}
	for i := 0; i < len(p.board); i++ {
		p.board[i] = make([]playerColor, maxLen)
	}
	p.load()
	return p
}

func (A *AIPlayer) calId(colors []playerColor) int {
	if len(colors) != 9 {
		panic(len(colors))
	}
	result := 0
	for i, c := range colors {
		if i != 4 {
			result = result*4 + int(c) + 1
		}
	}
	return result
}

func (A *AIPlayer) copy() *AIPlayer {
	p := &AIPlayer{
		board:  make([][]playerColor, maxLen),
		pColor: colorBlack,
		cache:  make([]float64, 0x10000),
	}
	for i := 0; i < len(p.board); i++ {
		p.board[i] = make([]playerColor, maxLen)
		copy(p.cache, A.cache)
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
		for _, dir := range fourDirections {
			line := make([]playerColor, 9)
			for i := -4; i <= 4; i++ {
				p2 := p.move(dir, i)
				if p2.checkRange() {
					line[i+4] = A.board[p2.y][p2.x]
				} else {
					line[i+4] = -1
				}
			}
			id := A.calId(line)
			result += A.cache[id]
		}
		return result
	}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			if A.board[i][j] == colorEmpty {
				p := point{j, i}
				v := f(p)
				if !maxP.checkRange() || v > maxV || v == maxV && p.nearMidThan(maxP) {
					maxV = v
					maxP = p
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
	f := func(b []float64) {
		for i := 0; i < len(b); i++ {
			b[i] += rand.NormFloat64()
		}
	}
	f(A.cache)
}

func (A *AIPlayer) load() {
	buf, err := ioutil.ReadFile("data.dat")
	if err != nil {
		log.Println(err.Error())
		A.rand()
		for i := 0; i < 0x10000; i++ {
			colors := A.parseId(i)
			if A.isFive(colors, A.pColor) {
				A.cache[i] += 50000
			} else if A.isFive(colors, A.pColor.conversion()) {
				A.cache[i] += 45000
			}
		}
		return
	}
	f := func(b []float64, buf []byte) {
		for i := 0; i < len(b); i++ {
			u := binary.LittleEndian.Uint64(buf[i*8 : i*8+8])
			b[i] = math.Float64frombits(u)
		}
	}
	f(A.cache, buf)
}

func (A *AIPlayer) save() {
	f := func(b []float64) []byte {
		var ret []byte
		for i := 0; i < len(b); i++ {
			u := math.Float64bits(b[i])
			buf := make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, u)
			ret = append(ret, buf...)
		}
		return ret
	}
	b1 := f(A.cache)
	err := ioutil.WriteFile("data.dat", b1, 0644)
	if err != nil {
		log.Println(err.Error())
	}
}

func (A *AIPlayer) parseId(id int) []playerColor {
	ret := make([]playerColor, 9)
	i := 8
	for id > 0 {
		if i == 4 {
			i = 3
		}
		ret[i] = playerColor(id%4 - 1)
		id /= 4
		i--
	}
	return ret
}

func (A *AIPlayer) isFive(colors []playerColor, targetColor playerColor) bool {
	count := 0
	for _, c := range colors {
		if c == targetColor {
			count++
			if count >= 5 {
				return true
			}
		} else {
			count = 0
		}
	}
	return false
}
