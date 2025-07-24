package main

import "fmt"

const maxLen = 15

type direction struct {
	x, y int
}

var fourDirections = []direction{{-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

var eightDirections = []direction{{-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}}

type point struct {
	x, y int
}

func (p point) move(dir direction, length int) point {
	if length == 0 {
		return p
	}
	return point{p.x + dir.x*length, p.y + dir.y*length}
}

func (p point) checkRange() bool {
	return p.x < maxLen && p.x >= 0 && p.y < maxLen && p.y >= 0
}

func (p point) nearMidThan(p2 point) bool {
	return max(abs(p.x-maxLen/2), abs(p.y-maxLen/2)) < max(abs(p2.x-maxLen/2), abs(p2.y-maxLen/2))
}

func (p point) hash() int {
	return p.y*maxLen + p.x
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func max(x ...int) int {
	m := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] > m {
			m = x[i]
		}
	}
	return m
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
