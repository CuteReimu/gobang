package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

type robotPlayer struct {
	boardStatus
	boardCache
	pColor            playerColor
	maxLevelCount     int
	maxCountEachLevel int
	maxCheckmateCount int
}

func newRobotPlayer(color playerColor) player {
	rp := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            color,
		maxLevelCount:     6,
		maxCountEachLevel: 16,
		maxCheckmateCount: 12,
	}
	rp.initBoardStatus()
	return rp
}

func (r *robotPlayer) color() playerColor {
	return r.pColor
}

func (r *robotPlayer) play() (point, error) {
	if r.count == 0 {
		p := point{maxLen / 2, maxLen / 2}
		r.set(p, r.pColor)
		return p, nil
	}
	p1, ok := r.findForm5(r.pColor)
	if ok {
		r.set(p1, r.pColor)
		return p1, nil
	}
	p1, ok = r.stop4(r.pColor)
	if ok {
		r.set(p1, r.pColor)
		return p1, nil
	}
	for i := 2; i <= r.maxCheckmateCount; i += 2 {
		if p, ok := r.calculateKill(r.pColor, true, i); ok {
			return p, nil
		}
	}
	result := r.max(r.maxLevelCount, 100000000)
	if result == nil {
		return point{}, errors.New("algorithm error")
	}
	r.set(result.p, r.pColor)
	return result.p, nil
}

func (r *robotPlayer) calculateKill(color playerColor, aggressive bool, step int) (point, bool) {
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == 0 {
				r.set(p, color)
				if !r.exists4(color.conversion()) && (!aggressive || r.exists4(color)) {
					if _, ok := r.calculateKill(color.conversion(), !aggressive, step-1); !ok {
						r.set(p, 0)
						return p, true
					}
				}
				r.set(p, 0)
			}
		}
	}
	return p, false
}

func (r *robotPlayer) stop4(color playerColor) (point, bool) {
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == colorEmpty {
				for _, dir := range fourDirections {
					leftCount, rightCount := 0, 0
					for k := -1; k >= -4; k-- {
						if p1 := p.move(dir, k); p1.checkRange() && r.get(p1) == color.conversion() {
							leftCount++
						} else {
							break
						}
					}
					for k := 1; k <= 4; k++ {
						if p1 := p.move(dir, k); p1.checkRange() && r.get(p1) == color.conversion() {
							rightCount++
						} else {
							break
						}
					}
					if leftCount+rightCount >= 4 {
						return p, true
					}
				}
			}
		}
	}
	return p, false
}

func (r *robotPlayer) exists4(color playerColor) bool {
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == color || r.get(p) == colorEmpty {
				for _, dir := range fourDirections {
					count0, count1 := 0, 0
					for k := 0; k <= 4; k++ {
						pk := p.move(dir, k)
						if pk.checkRange() {
							kColor := r.get(pk)
							if kColor == 0 {
								count0++
							} else if kColor == color {
								count1++
							}
						}
					}
					if count0 == 1 && count1 == 4 {
						return true
					}
				}
			}
		}
	}
	return false
}

func (r *robotPlayer) findForm5(color playerColor) (point, bool) {
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == colorEmpty {
				for _, dir := range fourDirections {
					leftCount, rightCount := 0, 0
					for k := -1; k >= -4; k-- {
						if pk := p.move(dir, k); pk.checkRange() && r.get(pk) == color {
							leftCount++
						} else {
							break
						}
					}
					for k := 1; k <= 4; k++ {
						if pk := p.move(dir, k); pk.checkRange() && r.get(pk) == color {
							rightCount++
						} else {
							break
						}
					}
					if leftCount+rightCount >= 4 {
						return p, true
					}
				}
			}
		}
	}
	return p, false
}

func (r *robotPlayer) checkForm5ByPoint(p point, color playerColor) bool {
	if r.get(p) != 0 {
		return false
	}
	r.set(p, color)
	count := 0
	for _, dir := range fourDirections {
		count = 0
		for i := -4; i <= 4; i++ {
			p2 := p.move(dir, i)
			if p2.checkRange() && r.get(p2) == color {
				count++
			} else {
				count = 0
			}
			if count <= i || count == 5 {
				break
			}
		}
		if count == 5 {
			break
		}
	}
	r.set(p, colorEmpty)
	return count == 5
}

func (r *robotPlayer) display(p point) error {
	if r.get(p) != 0 {
		return errors.New(fmt.Sprintf("illegal argument: %s%s", p, r.get(p)))
	}
	r.set(p, r.pColor.conversion())
	return nil
}

func (r *robotPlayer) max(step int, foundminVal int) *pointAndValue {
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}
	var queue pointAndValueSlice
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == 0 && r.isNeighbor(p) {
				evathis := r.evaluatePoint(p, r.pColor)
				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)
	if step == 1 {
		if len(queue) == 0 {
			log.Println("algorithm error")
			return nil
		}
		p = queue[0].p
		r.setIfEmpty(p, r.pColor)
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.putIntoCache(r.hash, step, result)
		return result
	}
	maxPoint := point{}
	maxVal := -100000000
	i := 0
	for _, obj := range queue {
		i++
		if i > r.maxCountEachLevel {
			break
		}
		p = obj.p
		r.set(p, r.pColor)
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal > 800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}
		evathis := r.min(step-1, maxVal).value //最大值最小值法
		if evathis >= foundminVal {
			r.set(p, 0)
			result := &pointAndValue{p, evathis}
			r.putIntoCache(r.hash, step, result)
			return result
		}
		if evathis > maxVal || evathis == maxVal && p.nearMidThan(maxPoint) {
			maxVal = evathis
			maxPoint = p
		}
		r.set(p, 0)
	}
	if maxVal < -99999999 {
		return nil
	}
	result := &pointAndValue{maxPoint, maxVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

func (r *robotPlayer) min(step int, foundmaxVal int) *pointAndValue {
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}
	var queue pointAndValueSlice
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == 0 && r.isNeighbor(p) {
				evathis := r.evaluatePoint(p, r.pColor.conversion())
				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)
	if step == 1 {
		if len(queue) == 0 {
			log.Println("algorithm error")
			return nil
		}
		p := queue[0].p
		r.setIfEmpty(p, r.pColor.conversion())
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, 0)
		result := &pointAndValue{p, val}
		r.putIntoCache(r.hash, step, result)
		return result
	}
	var minPoint point
	minVal := 100000000
	i := 0
	for _, obj := range queue {
		i++
		if i > r.maxCountEachLevel {
			break
		}
		p = obj.p
		r.set(p, r.pColor.conversion())
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal < -800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}
		evathis := r.max(step-1, minVal).value //最大值最小值法
		if evathis <= foundmaxVal {
			r.set(p, 0)
			result := &pointAndValue{p, evathis}
			r.putIntoCache(r.hash, step, result)
			return result
		}
		if evathis < minVal || evathis == minVal && p.nearMidThan(minPoint) {
			minVal = evathis
			minPoint = p
		}
		r.set(p, 0)
	}
	if minVal > 99999999 {
		return nil
	}
	result := &pointAndValue{minPoint, minVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

func (r *robotPlayer) evaluatePoint(p point, color playerColor) int {
	return r.evaluatePoint2(p, color, colorBlack) + r.evaluatePoint2(p, color, colorWhite)
}

func (r *robotPlayer) evaluatePoint2(p point, me playerColor, plyer playerColor) (value int) {
	numoftwo := 0
	getLine := func(p point, dir direction, j int) playerColor {
		p2 := p.move(dir, j)
		if p2.checkRange() {
			return r.get(p2)
		}
		return -1
	}
	for _, dir := range eightDirections { // 8个方向
		// 活四 01111* *代表当前空位置 0代表其他空位置 下同
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer && getLine(p, dir, -5) == 0 {
			value += 300000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四A 21111*
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer && (getLine(p, dir, -5) == plyer.conversion() || getLine(p, dir, -5) == -1) {
			value += 250000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四B 111*1
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, 1) == plyer {
			value += 240000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 死四C 11*11
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer && getLine(p, dir, 2) == plyer {
			value += 230000
			if me != plyer {
				value -= 500
			}
			continue
		}
		// 活三 近3位置 111*0
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer {
			if getLine(p, dir, 1) == 0 {
				value += 1450
				if getLine(p, dir, -4) == 0 {
					value += 6000
					if me != plyer {
						value -= 300
					}
				}
			}
			if (getLine(p, dir, 1) == plyer.conversion() || getLine(p, dir, 1) == -1) && getLine(p, dir, -4) == 0 {
				value += 500
			}
			if (getLine(p, dir, -4) == plyer.conversion() || getLine(p, dir, -4) == -1) && getLine(p, dir, 1) == 0 {
				value += 500
			}
			continue
		}
		// 活三 远3位置 1110*
		if getLine(p, dir, -1) == 0 && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer {
			value += 350
			continue
		}
		// 死三 11*1
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer {
			value += 700
			if getLine(p, dir, -3) == 0 && getLine(p, dir, 2) == 0 {
				value += 6700
				continue
			}
			if (getLine(p, dir, -3) == plyer.conversion() || getLine(p, dir, -3) == -1) && (getLine(p, dir, 2) == plyer.conversion() || getLine(p, dir, 2) == -1) {
				value -= 700
				continue
			} else {
				value += 800
				continue
			}
		}
		// 活二的个数（因为会算2次，就2倍）
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == 0 && getLine(p, dir, 1) == 0 {
			if getLine(p, dir, 2) == 0 || getLine(p, dir, -4) == 0 {
				numoftwo += 2
			} else {
				value += 250
			}
		}
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == 0 && getLine(p, dir, 2) == plyer && getLine(p, dir, 1) == 0 && getLine(p, dir, 3) == 0 {
			numoftwo += 2
		}
		if getLine(p, dir, -1) == 0 && getLine(p, dir, 4) == 0 && getLine(p, dir, 3) == plyer && (getLine(p, dir, 2) == plyer && getLine(p, dir, 1) == 0 || getLine(p, dir, 1) == plyer && getLine(p, dir, 2) == 0) {
			numoftwo += 2
		}
		if getLine(p, dir, -1) == plyer && getLine(p, dir, 1) == plyer && getLine(p, dir, -2) == 0 && getLine(p, dir, 2) == 0 {
			if getLine(p, dir, 3) == 0 || getLine(p, dir, -3) == 0 {
				numoftwo++
			} else {
				value += 125
			}
		}
		// 其余散棋
		numOfplyer := 0
		for k := -4; k <= 0; k++ { // ++++* +++*+ ++*++ +*+++ *++++
			temp := 0
			for l := 0; l <= 4; l++ {
				if getLine(p, dir, k+l) == plyer {
					temp += 5 - abs(k+l)
				} else if getLine(p, dir, k+l) == plyer.conversion() || getLine(p, dir, k+l) == -1 {
					temp = 0
					break
				}
			}
			numOfplyer += temp
		}
		value += numOfplyer * 5
	}
	numoftwo /= 2
	if numoftwo >= 2 {
		value += 3000
		if me != plyer {
			value -= 100
		}
	} else if numoftwo == 1 {
		value += 2725
		if me != plyer {
			value -= 10
		}
	}
	return
}

func (r *robotPlayer) evaluateBoard(color playerColor) (values int) {
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) != color {
				continue
			}
			for _, dir := range eightDirections {
				colors := make([]playerColor, 9)
				for k := 0; k < 9; k++ {
					pk := p.move(dir, k-4)
					if pk.checkRange() {
						colors[k] = r.get(pk)
					} else {
						colors[k] = playerColor(-1)
					}
				}
				if colors[5] == color && colors[6] == color && colors[7] == color && colors[8] == color {
					values += 1000000
					continue
				}
				if colors[5] == color && colors[6] == color && colors[7] == color && colors[3] == 0 {
					if colors[8] == 0 { //?AAAA?
						values += 300000 / 2
					} else if colors[8] != color { //AAAA?
						values += 25000
					}
					continue
				}
				if colors[5] == color && colors[6] == color {
					if colors[7] == 0 && colors[8] == color { //AAA?A
						values += 30000
						continue
					}
					if colors[3] == 0 && colors[7] == 0 {
						if colors[2] == 0 && colors[8] != color || colors[8] == 0 && colors[2] != color { //??AAA??
							values += 22000 / 2
						} else if colors[2] != color && colors[2] != 0 && colors[8] != color && colors[8] != 0 { //?AAA?
							values += 500 / 2
						}
						continue
					}
					if colors[3] != 0 && colors[3] != color && colors[7] == 0 && colors[8] == 0 { //AAA??
						values += 500
						continue
					}
				}
				if colors[5] == color && colors[6] == 0 && colors[7] == color && colors[8] == color { //AA?AA
					values += 26000 / 2
					continue
				}
				if colors[5] == 0 && colors[6] == color && colors[7] == color {
					if colors[3] == 0 && colors[8] == 0 { //?A?AA?
						values += 22000
					} else if (colors[3] != 0 && colors[3] != color && colors[8] == 0) || (colors[8] != 0 && colors[8] != color && colors[3] == 0) { //A?AA? ?A?AA
						values += 800
					}
					continue
				}
				if colors[5] == 0 && colors[8] == color {
					if colors[6] == 0 && colors[7] == color { //A??AA
						values += 600
					} else if colors[6] == color && colors[7] == 0 { //A?A?A
						values += 550 / 2
					}
					continue
				}
				if colors[5] == color {
					if colors[3] == 0 && colors[6] == 0 {
						if colors[1] == 0 && colors[2] == 0 && colors[7] != 0 && colors[7] != color || colors[8] == 0 && colors[7] == 0 && colors[2] != 0 && colors[2] != color { //??AA??
							values += 650 / 2
						} else if colors[2] != 0 && colors[2] != color && colors[7] == 0 && colors[8] != 0 && colors[8] != color { //?AA??
							values += 150
						}
					} else if colors[3] != 0 && colors[3] != color && colors[6] == 0 && colors[7] == 0 && colors[8] == 0 { //AA???
						values += 150
					}
					continue
				}
				if colors[5] == 0 && colors[6] == color {
					if colors[3] == 0 && colors[7] == 0 {
						if colors[2] != 0 && colors[2] != color && colors[8] == 0 || colors[2] == 0 && colors[8] != 0 && colors[8] != color { //??A?A??
							values += 250 / 2
						}
						if colors[2] != 0 && colors[2] != color && colors[8] != 0 && colors[8] != color { //?A?A?
							values += 150 / 2
						}
					} else if colors[3] != 0 && colors[3] != color && colors[7] == 0 && colors[8] == 0 { //A?A??
						values += 150
					}
					continue
				}
				if colors[5] == 0 && colors[6] == 0 && colors[7] == color {
					if colors[3] == 0 && colors[8] == 0 { //?A??A?
						values += 200 / 2
						continue
					}
					if colors[3] != 0 && colors[3] != color && colors[8] == 0 { //A??A?
						p5 := p.move(dir, 5)
						if p5.checkRange() {
							color5 := r.get(p5)
							if color5 == 0 {
								values += 200
							} else if color5 != color {
								values += 150
							}
						}
					}
					continue
				}
			}
		}
	}
	return values
}

type pointAndValue struct {
	p     point
	value int
}

type pointAndValueSlice []*pointAndValue

func (s pointAndValueSlice) Len() int {
	return len(s)
}

func (s pointAndValueSlice) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s pointAndValueSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
