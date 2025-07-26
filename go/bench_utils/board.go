package main

import (
	"log"
	"math/rand"
)

type boardStatus struct {
	blackHash [][]uint64
	whiteHash [][]uint64
	board     [][]playerColor
	hash      uint64
	count     int
}

func (b *boardStatus) initBoardStatus() {
	b.blackHash = make([][]uint64, maxLen)
	b.whiteHash = make([][]uint64, maxLen)
	b.board = make([][]playerColor, maxLen)
	r := rand.New(rand.NewSource(1551980916123)) //rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < maxLen; i++ {
		b.blackHash[i] = make([]uint64, maxLen)
		b.whiteHash[i] = make([]uint64, maxLen)
		b.board[i] = make([]playerColor, maxLen)
		for j := 0; j < maxLen; j++ {
			b.blackHash[i][j] = r.Uint64()
			b.whiteHash[i][j] = r.Uint64()
		}
	}
}

func (b *boardStatus) setIfEmpty(p point, color playerColor) bool {
	if b.board[p.y][p.x] != colorEmpty {
		return false
	}
	switch color {
	case colorEmpty:
		return true
	case colorBlack:
		b.hash ^= b.blackHash[p.y][p.x]
	case colorWhite:
		b.hash ^= b.whiteHash[p.y][p.x]
	default:
		log.Printf("illegal argument: %s%s\n", p, color)
		return false
	}
	b.board[p.y][p.x] = color
	b.count++
	return true
}

func (b *boardStatus) set(p point, color playerColor) {
	if b.board[p.y][p.x] == color {
		return
	}
	switch color {
	case colorEmpty:
	case colorBlack:
		b.hash ^= b.blackHash[p.y][p.x]
		b.count++
	case colorWhite:
		b.hash ^= b.whiteHash[p.y][p.x]
		b.count++
	default:
		log.Printf("illegal argument: %s%s\n", p, color)
		return
	}
	switch b.board[p.y][p.x] {
	case colorBlack:
		b.hash ^= b.blackHash[p.y][p.x]
		b.count--
	case colorWhite:
		b.hash ^= b.whiteHash[p.y][p.x]
		b.count--
	}
	b.board[p.y][p.x] = color
}

func (b *boardStatus) get(p point) playerColor {
	return b.board[p.y][p.x]
}

func (b *boardStatus) isNeighbor(p point) bool {
	if !p.checkRange() {
		return false
	}
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			p2 := point{p.x + j, p.y + i}
			if p2.checkRange() && b.get(p2) > colorEmpty {
				return true
			}
		}
	}
	return false
}

type boardCache map[uint64]map[int]*pointAndValue

func (b boardCache) putIntoCache(key uint64, deep int, val *pointAndValue) {
	m, ok := b[key]
	if !ok {
		m = make(map[int]*pointAndValue)
		b[key] = m
	}
	m[deep] = val
}

func (b boardCache) getFromCache(key uint64, deep int) *pointAndValue {
	m, ok := b[key]
	if !ok {
		return nil
	}
	v, _ := m[deep]
	return v
}
