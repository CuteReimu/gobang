package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

// EvaluationParams holds configurable evaluation parameters
type EvaluationParams struct {
	// Pattern values for evaluatePoint2
	LiveFour             int // 活四
	DeadFourA            int // 死四A
	DeadFourB            int // 死四B
	DeadFourC            int // 死四C
	LiveThreeNear        int // 活三 近3位置
	LiveThreeBonus       int // 活三额外奖励
	LiveThreeFar         int // 活三 远3位置
	DeadThree            int // 死三
	DeadThreeBonus       int // 死三额外奖励
	TwoCount2            int // 活二×2的奖励
	TwoCount1            int // 活二×1的奖励
	ScatterMultiplier    int // 散棋乘数
	OpponentPenalty      int // 对手惩罚
	OpponentMinorPenalty int // 对手小惩罚

	// Pattern values for evaluateBoard
	FiveInRow          int            // 五连珠
	FourInRowOpen      int            // 活四
	FourInRowClosed    int            // 死四
	ThreeInRowVariants map[string]int // 活三的各种变体
}

// getDefaultEvaluationParams returns the default evaluation parameters
func getDefaultEvaluationParams() *EvaluationParams {
	return &EvaluationParams{
		LiveFour:             300000,
		DeadFourA:            250000,
		DeadFourB:            240000,
		DeadFourC:            230000,
		LiveThreeNear:        1450,
		LiveThreeBonus:       6000,
		LiveThreeFar:         350,
		DeadThree:            700,
		DeadThreeBonus:       6700,
		TwoCount2:            3000,
		TwoCount1:            2725,
		ScatterMultiplier:    5,
		OpponentPenalty:      500,
		OpponentMinorPenalty: 300,
		FiveInRow:            1000000,
		FourInRowOpen:        300000,
		FourInRowClosed:      25000,
		ThreeInRowVariants: map[string]int{
			"open":   22000,
			"semi":   500,
			"closed": 26000,
			"gap":    800,
			"basic":  650,
			"corner": 150,
		},
	}
}

type robotPlayer struct {
	boardStatus
	boardCache
	pColor            playerColor
	maxLevelCount     int
	maxCountEachLevel int
	maxCheckmateCount int
	evalParams        *EvaluationParams
}

func newRobotPlayer(color playerColor) player {
	rp := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            color,
		maxLevelCount:     6,
		maxCountEachLevel: 16,
		maxCheckmateCount: 12,
		evalParams:        getDefaultEvaluationParams(),
	}
	rp.initBoardStatus()
	return rp
}

// newOptimizedRobotPlayer creates a robot player with optimized parameters
func newOptimizedRobotPlayer(color playerColor) player {
	rp := &optimizedRobotPlayer{
		robotPlayer: robotPlayer{
			boardCache:        make(boardCache),
			pColor:            color,
			maxLevelCount:     6,  // Maximum depth for iterative deepening
			maxCountEachLevel: 16, // Balanced candidate count
			maxCheckmateCount: 12, // Full checkmate search for tactical strength
			evalParams:        getImprovedOptimizedEvaluationParams(),
		},
		evalCache: make(map[uint64]int),
	}
	rp.initBoardStatus()
	return rp
}

// newBalancedRobotPlayer creates a robot player with balanced speed and strength
func newBalancedRobotPlayer(color playerColor) player {
	rp := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            color,
		maxLevelCount:     4,  // Even depth for proper minimax evaluation
		maxCountEachLevel: 16, // More candidates to compensate for reduced depth
		maxCheckmateCount: 12, // Full checkmate search
		evalParams:        getBalancedEvaluationParams(),
	}
	rp.initBoardStatus()
	return rp
}

// optimizedRobotPlayer - improved optimized AI with better balance of speed and strength
type optimizedRobotPlayer struct {
	robotPlayer
	evalCache map[uint64]int // Cache for position evaluations
	nodeCount int            // For debugging
}

// getImprovedOptimizedEvaluationParams returns improved optimized evaluation parameters with better tactical awareness
func getImprovedOptimizedEvaluationParams() *EvaluationParams {
	return &EvaluationParams{
		LiveFour:             450000,  // Much higher priority for winning moves
		DeadFourA:            380000,  // Enhanced threat detection
		DeadFourB:            360000,  // Enhanced threat detection
		DeadFourC:            340000,  // Enhanced threat detection
		LiveThreeNear:        3000,    // Much better three-in-a-row evaluation
		LiveThreeBonus:       10000,   // Stronger tactical evaluation
		LiveThreeFar:         750,     // Better distant threat recognition
		DeadThree:            1200,    // Enhanced defensive evaluation
		DeadThreeBonus:       9000,    // Strong defensive bonus
		TwoCount2:            5000,    // Better two-count evaluation
		TwoCount1:            4200,    // Better single-two evaluation
		ScatterMultiplier:    9,       // Enhanced position evaluation
		OpponentPenalty:      750,     // Stronger opponent threat response
		OpponentMinorPenalty: 450,     // Better minor threat response
		FiveInRow:            1500000, // Maximum priority for wins
		FourInRowOpen:        450000,  // Maximum priority for winning threats
		FourInRowClosed:      45000,   // Better closed-four evaluation
		ThreeInRowVariants: map[string]int{
			"open":   45000, // Much stronger open three evaluation
			"semi":   1000,  // Better semi-open evaluation
			"closed": 50000, // Much stronger closed three
			"gap":    1500,  // Better gap pattern recognition
			"basic":  1200,  // Enhanced basic patterns
			"corner": 300,   // Better corner evaluation
		},
	}
}

// getOptimizedEvaluationParams returns optimized evaluation parameters
func getOptimizedEvaluationParams() *EvaluationParams {
	return &EvaluationParams{
		LiveFour:             320000,  // Slightly increased
		DeadFourA:            260000,  // Slightly increased
		DeadFourB:            245000,  // Slightly increased
		DeadFourC:            235000,  // Slightly increased
		LiveThreeNear:        1500,    // Slightly increased
		LiveThreeBonus:       6200,    // Slightly increased
		LiveThreeFar:         400,     // Slightly increased
		DeadThree:            750,     // Slightly increased
		DeadThreeBonus:       6800,    // Slightly increased
		TwoCount2:            3100,    // Slightly increased
		TwoCount1:            2800,    // Slightly increased
		ScatterMultiplier:    6,       // Slightly increased
		OpponentPenalty:      480,     // Slightly decreased for balance
		OpponentMinorPenalty: 280,     // Slightly decreased for balance
		FiveInRow:            1050000, // Increased for priority
		FourInRowOpen:        315000,  // Slightly increased
		FourInRowClosed:      26000,   // Slightly increased
		ThreeInRowVariants: map[string]int{
			"open":   23000, // Slightly increased
			"semi":   520,   // Slightly increased
			"closed": 27000, // Slightly increased
			"gap":    850,   // Slightly increased
			"basic":  680,   // Slightly increased
			"corner": 160,   // Slightly increased
		},
	}
}

// getBalancedEvaluationParams returns balanced evaluation parameters for stronger play
func getBalancedEvaluationParams() *EvaluationParams {
	return &EvaluationParams{
		LiveFour:             350000,  // Higher priority for winning moves
		DeadFourA:            280000,  // Higher threat detection
		DeadFourB:            265000,  // Higher threat detection
		DeadFourC:            250000,  // Higher threat detection
		LiveThreeNear:        2000,    // Improved three-in-a-row evaluation
		LiveThreeBonus:       8000,    // Stronger bonus for good positions
		LiveThreeFar:         500,     // Better distant threat recognition
		DeadThree:            900,     // Improved defensive evaluation
		DeadThreeBonus:       7500,    // Stronger defensive bonus
		TwoCount2:            3500,    // Better two-count evaluation
		TwoCount1:            3000,    // Better single-two evaluation
		ScatterMultiplier:    7,       // Improved position evaluation
		OpponentPenalty:      600,     // Stronger opponent threat response
		OpponentMinorPenalty: 350,     // Better minor threat response
		FiveInRow:            1200000, // Highest priority for wins
		FourInRowOpen:        350000,  // Higher priority for winning threats
		FourInRowClosed:      30000,   // Better closed-four evaluation
		ThreeInRowVariants: map[string]int{
			"open":   28000, // Stronger open three evaluation
			"semi":   650,   // Better semi-open evaluation
			"closed": 32000, // Stronger closed three
			"gap":    1000,  // Better gap pattern recognition
			"basic":  800,   // Improved basic patterns
			"corner": 200,   // Better corner evaluation
		},
	}
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

	// Use iterative deepening for better time management
	result := r.iterativeDeepening()
	if result == nil {
		return point{}, errors.New("algorithm error")
	}
	r.set(result.p, r.pColor)
	return result.p, nil
}

// iterativeDeepening implements iterative deepening for better time management
func (r *robotPlayer) iterativeDeepening() *pointAndValue {
	var bestResult *pointAndValue

	// Adaptive depth based on game phase and threats
	maxDepth := r.getAdaptiveDepth()

	// Start with shallow searches and progressively deepen
	for depth := 2; depth <= maxDepth; depth++ {
		result := r.max(depth, 100000000)
		if result != nil {
			bestResult = result
		}

		// Early termination for strong positions
		if bestResult != nil && bestResult.value > 800000 {
			break
		}

		// If we find a very good move early, don't spend more time
		if depth >= 4 && bestResult != nil && bestResult.value > 200000 {
			break
		}
	}

	return bestResult
}

// getAdaptiveDepth returns adaptive search depth based on game state (always even)
func (r *robotPlayer) getAdaptiveDepth() int {
	baseDepth := r.maxLevelCount

	// Check for immediate threats that require deeper analysis
	if r.hasImmediateThreats() {
		return baseDepth + 2 // Deeper search for tactical positions (maintains even depth)
	}

	// In opening, use slightly less depth for speed (ensure even number)
	if r.count < 8 {
		adjusted := baseDepth - 2
		if adjusted < 2 {
			adjusted = 2
		}
		return adjusted
	}

	// In middle game with many pieces, use standard depth
	if r.count >= 8 && r.count < 20 {
		return baseDepth
	}

	// In endgame, use deeper search (maintain even depth)
	return baseDepth + 2
}

// hasImmediateThreats checks if there are immediate tactical threats on the board
func (r *robotPlayer) hasImmediateThreats() bool {
	// Check if opponent has 4 in a row (immediate win threat)
	if r.exists4(r.pColor.conversion()) {
		return true
	}

	// Check if we have 4 in a row (immediate win opportunity)
	if r.exists4(r.pColor) {
		return true
	}

	// Check for multiple threats
	threatsCount := r.countThreats(r.pColor) + r.countThreats(r.pColor.conversion())
	return threatsCount >= 2
}

// countThreats counts the number of three-in-a-row threats for a given color
func (r *robotPlayer) countThreats(color playerColor) int {
	threats := 0
	p := point{}

	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == colorEmpty {
				// Check if placing a piece here creates a threat
				r.set(p, color)

				for _, dir := range fourDirections {
					count := 1
					// Count in positive direction
					for k := 1; k < 5; k++ {
						pk := p.move(dir, k)
						if pk.checkRange() && r.get(pk) == color {
							count++
						} else {
							break
						}
					}
					// Count in negative direction
					for k := 1; k < 5; k++ {
						pk := p.move(dir, -k)
						if pk.checkRange() && r.get(pk) == color {
							count++
						} else {
							break
						}
					}

					if count >= 3 {
						threats++
						break // Only count once per position
					}
				}

				r.set(p, colorEmpty)
			}
		}
	}

	return threats
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
				evathis := r.evaluatePoint2(p, r.pColor, r.pColor)
				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)

	// Adaptive candidate count based on game phase
	maxCandidates := r.getAdaptiveCandidateCount(len(queue))

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
		if i > maxCandidates {
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
				evathis := r.evaluatePoint2(p, r.pColor.conversion(), r.pColor.conversion())
				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)

	// Adaptive candidate count based on game phase
	maxCandidates := r.getAdaptiveCandidateCount(len(queue))

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
		if i > maxCandidates {
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

// getAdaptiveCandidateCount returns adaptive candidate count based on game phase
func (r *robotPlayer) getAdaptiveCandidateCount(totalCandidates int) int {
	// In early game (fewer pieces), consider more candidates
	// In late game (more pieces), focus on fewer but better candidates
	if r.count < 10 {
		return min(r.maxCountEachLevel+4, totalCandidates)
	} else if r.count < 20 {
		return min(r.maxCountEachLevel, totalCandidates)
	} else {
		return min(r.maxCountEachLevel-2, totalCandidates)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Enhanced evaluatePoint for better move ordering and alpha-beta pruning efficiency
func (r *optimizedRobotPlayer) evaluatePoint(p point, color playerColor) int {
	// Start with base evaluation
	baseValue := r.robotPlayer.evaluatePoint2(p, color, color)

	// Add tactical bonuses for better move ordering
	tacticalBonus := 0

	// Simulate placing the piece
	r.set(p, color)

	// Check for immediate wins (highest priority)
	if r.checkForm5ByPoint(p, color) {
		tacticalBonus += 2000000
	}

	// Check for creating live fours (very high priority)
	if r.createsLiveFour(p, color) {
		tacticalBonus += 500000
	}

	// Check for blocking opponent's live fours (very high priority)
	if r.createsLiveFour(p, color.conversion()) {
		tacticalBonus += 450000
	}

	// Check for creating multiple threats (high priority)
	if r.createsMultipleThreats(p, color) {
		tacticalBonus += 300000
	}

	// Check for creating live threes (medium-high priority)
	liveThrees := r.countLiveThreesAt(p, color)
	tacticalBonus += liveThrees * 50000

	// Check for blocking opponent's live threes (medium priority)
	opponentLiveThrees := r.countLiveThreesAt(p, color.conversion())
	tacticalBonus += opponentLiveThrees * 30000

	// Bonus for center positions in early game
	if r.count < 10 {
		center := maxLen / 2
		distance := abs(p.x-center) + abs(p.y-center)
		tacticalBonus += (10 - distance) * 100
	}

	// Remove the piece
	r.set(p, colorEmpty)

	return baseValue + tacticalBonus
}

// Helper functions for tactical evaluation

func (r *optimizedRobotPlayer) createsLiveFour(p point, color playerColor) bool {
	// Check all four directions for live four patterns
	for _, dir := range fourDirections {
		count := 1 // The piece we're placing

		// Count pieces in positive direction
		for i := 1; i < 5; i++ {
			pos := p.move(dir, i)
			if !pos.checkRange() || r.get(pos) != color {
				break
			}
			count++
		}

		// Count pieces in negative direction
		for i := 1; i < 5; i++ {
			pos := p.move(dir, -i)
			if !pos.checkRange() || r.get(pos) != color {
				break
			}
			count++
		}

		// Check if it forms a live four (4 in a row with open ends)
		if count >= 4 {
			// Check if both ends are open
			leftEnd := p.move(dir, -(count - 1))
			rightEnd := p.move(dir, count)
			if leftEnd.checkRange() && rightEnd.checkRange() &&
				r.get(leftEnd) == colorEmpty && r.get(rightEnd) == colorEmpty {
				return true
			}
		}
	}
	return false
}

func (r *optimizedRobotPlayer) createsMultipleThreats(p point, color playerColor) bool {
	threatsCount := 0

	for _, dir := range fourDirections {
		if r.createsThreatenPattern(p, color, dir) {
			threatsCount++
		}
	}

	return threatsCount >= 2
}

func (r *optimizedRobotPlayer) createsThreatenPattern(p point, color playerColor, dir direction) bool {
	count := 1 // The piece we're placing

	// Count consecutive pieces in both directions
	for i := 1; i < 4; i++ {
		pos := p.move(dir, i)
		if !pos.checkRange() || r.get(pos) != color {
			break
		}
		count++
	}

	for i := 1; i < 4; i++ {
		pos := p.move(dir, -i)
		if !pos.checkRange() || r.get(pos) != color {
			break
		}
		count++
	}

	return count >= 3
}

func (r *optimizedRobotPlayer) countLiveThreesAt(p point, color playerColor) int {
	liveThrees := 0

	for _, dir := range fourDirections {
		if r.formsLiveThree(p, color, dir) {
			liveThrees++
		}
	}

	return liveThrees
}

func (r *optimizedRobotPlayer) formsLiveThree(p point, color playerColor, dir direction) bool {
	count := 1 // The piece we're placing

	// Simple live three check: exactly 3 pieces with open ends
	for i := 1; i < 3; i++ {
		pos := p.move(dir, i)
		if !pos.checkRange() || r.get(pos) != color {
			break
		}
		count++
	}

	for i := 1; i < 3; i++ {
		pos := p.move(dir, -i)
		if !pos.checkRange() || r.get(pos) != color {
			break
		}
		count++
	}

	if count == 3 {
		// Check if ends are open
		leftEnd := p.move(dir, -2)
		rightEnd := p.move(dir, 2)
		if leftEnd.checkRange() && rightEnd.checkRange() &&
			r.get(leftEnd) == colorEmpty && r.get(rightEnd) == colorEmpty {
			return true
		}
	}

	return false
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
			value += r.evalParams.LiveFour
			if me != plyer {
				value -= r.evalParams.OpponentPenalty
			}
			continue
		}
		// 死四A 21111*
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer && (getLine(p, dir, -5) == plyer.conversion() || getLine(p, dir, -5) == -1) {
			value += r.evalParams.DeadFourA
			if me != plyer {
				value -= r.evalParams.OpponentPenalty
			}
			continue
		}
		// 死四B 111*1
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, 1) == plyer {
			value += r.evalParams.DeadFourB
			if me != plyer {
				value -= r.evalParams.OpponentPenalty
			}
			continue
		}
		// 死四C 11*11
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer && getLine(p, dir, 2) == plyer {
			value += r.evalParams.DeadFourC
			if me != plyer {
				value -= r.evalParams.OpponentPenalty
			}
			continue
		}
		// 活三 近3位置 111*0
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer {
			if getLine(p, dir, 1) == 0 {
				value += r.evalParams.LiveThreeNear
				if getLine(p, dir, -4) == 0 {
					value += r.evalParams.LiveThreeBonus
					if me != plyer {
						value -= r.evalParams.OpponentMinorPenalty
					}
				}
			}
			if (getLine(p, dir, 1) == plyer.conversion() || getLine(p, dir, 1) == -1) && getLine(p, dir, -4) == 0 {
				value += r.evalParams.OpponentPenalty
			}
			if (getLine(p, dir, -4) == plyer.conversion() || getLine(p, dir, -4) == -1) && getLine(p, dir, 1) == 0 {
				value += r.evalParams.OpponentPenalty
			}
			continue
		}
		// 活三 远3位置 1110*
		if getLine(p, dir, -1) == 0 && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer {
			value += r.evalParams.LiveThreeFar
			continue
		}
		// 死三 11*1
		if getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer {
			value += r.evalParams.DeadThree
			if getLine(p, dir, -3) == 0 && getLine(p, dir, 2) == 0 {
				value += r.evalParams.DeadThreeBonus
				continue
			}
			if (getLine(p, dir, -3) == plyer.conversion() || getLine(p, dir, -3) == -1) && (getLine(p, dir, 2) == plyer.conversion() || getLine(p, dir, 2) == -1) {
				value -= r.evalParams.DeadThree
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
		value += numOfplyer * r.evalParams.ScatterMultiplier
	}
	numoftwo /= 2
	if numoftwo >= 2 {
		value += r.evalParams.TwoCount2
		if me != plyer {
			value -= 100
		}
	} else if numoftwo == 1 {
		value += r.evalParams.TwoCount1
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
					values += r.evalParams.FiveInRow
					continue
				}
				if colors[5] == color && colors[6] == color && colors[7] == color && colors[3] == 0 {
					if colors[8] == 0 { //?AAAA?
						values += r.evalParams.FourInRowOpen / 2
					} else if colors[8] != color { //AAAA?
						values += r.evalParams.FourInRowClosed
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
							values += r.evalParams.ThreeInRowVariants["open"] / 2
						} else if colors[2] != color && colors[2] != 0 && colors[8] != color && colors[8] != 0 { //?AAA?
							values += r.evalParams.ThreeInRowVariants["semi"] / 2
						}
						continue
					}
					if colors[3] != 0 && colors[3] != color && colors[7] == 0 && colors[8] == 0 { //AAA??
						values += r.evalParams.ThreeInRowVariants["semi"]
						continue
					}
				}
				if colors[5] == color && colors[6] == 0 && colors[7] == color && colors[8] == color { //AA?AA
					values += r.evalParams.ThreeInRowVariants["closed"] / 2
					continue
				}
				if colors[5] == 0 && colors[6] == color && colors[7] == color {
					if colors[3] == 0 && colors[8] == 0 { //?A?AA?
						values += r.evalParams.ThreeInRowVariants["open"]
					} else if (colors[3] != 0 && colors[3] != color && colors[8] == 0) || (colors[8] != 0 && colors[8] != color && colors[3] == 0) { //A?AA? ?A?AA
						values += r.evalParams.ThreeInRowVariants["gap"]
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
							values += r.evalParams.ThreeInRowVariants["basic"] / 2
						} else if colors[2] != 0 && colors[2] != color && colors[7] == 0 && colors[8] != 0 && colors[8] != color { //?AA??
							values += r.evalParams.ThreeInRowVariants["corner"]
						}
					} else if colors[3] != 0 && colors[3] != color && colors[6] == 0 && colors[7] == 0 && colors[8] == 0 { //AA???
						values += r.evalParams.ThreeInRowVariants["corner"]
					}
					continue
				}
				if colors[5] == 0 && colors[6] == color {
					if colors[3] == 0 && colors[7] == 0 {
						if colors[2] != 0 && colors[2] != color && colors[8] == 0 || colors[2] == 0 && colors[8] != 0 && colors[8] != color { //??A?A??
							values += 250 / 2
						}
						if colors[2] != 0 && colors[2] != color && colors[8] != 0 && colors[8] != color { //?A?A?
							values += r.evalParams.ThreeInRowVariants["corner"] / 2
						}
					} else if colors[3] != 0 && colors[3] != color && colors[7] == 0 && colors[8] == 0 { //A?A??
						values += r.evalParams.ThreeInRowVariants["corner"]
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
								values += r.evalParams.ThreeInRowVariants["corner"]
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

// SelfPlayResult holds the result of a self-play game
type SelfPlayResult struct {
	Winner   playerColor
	Moves    int
	Duration int // in milliseconds
}

// adjustParameters adjusts evaluation parameters based on game outcomes
func (r *robotPlayer) adjustParameters(results []SelfPlayResult) {
	if len(results) < 5 {
		return // Need at least 5 games for adjustment
	}

	winRate := r.calculateWinRate(results)
	avgMoves := r.calculateAverageMovesPerGame(results)

	// If win rate is too low, make AI more aggressive
	if winRate < 0.4 {
		r.evalParams.LiveFour += 10000
		r.evalParams.DeadFourA += 8000
		r.evalParams.LiveThreeNear += 100
		r.evalParams.LiveThreeBonus += 500
	}

	// If games are too long, prioritize quicker wins
	if avgMoves > 50 {
		r.evalParams.FiveInRow += 50000
		r.evalParams.FourInRowOpen += 15000
	}

	// If games are too short, encourage more strategic play
	if avgMoves < 25 {
		r.evalParams.ThreeInRowVariants["open"] += 1000
		r.evalParams.ThreeInRowVariants["semi"] += 500
	}
}

func (r *robotPlayer) calculateWinRate(results []SelfPlayResult) float64 {
	wins := 0
	for _, result := range results {
		if result.Winner == r.pColor {
			wins++
		}
	}
	return float64(wins) / float64(len(results))
}

func (r *robotPlayer) calculateAverageMovesPerGame(results []SelfPlayResult) float64 {
	totalMoves := 0
	for _, result := range results {
		totalMoves += result.Moves
	}
	return float64(totalMoves) / float64(len(results))
}

// createPlayerCopy creates a copy of the robot player with the same parameters
func (r *robotPlayer) createPlayerCopy(color playerColor) *robotPlayer {
	rp := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            color,
		maxLevelCount:     r.maxLevelCount,
		maxCountEachLevel: r.maxCountEachLevel,
		maxCheckmateCount: r.maxCheckmateCount,
		evalParams:        r.copyEvaluationParams(),
	}
	rp.initBoardStatus()
	return rp
}

func (r *robotPlayer) copyEvaluationParams() *EvaluationParams {
	copied := *r.evalParams
	// Deep copy the map
	copied.ThreeInRowVariants = make(map[string]int)
	for k, v := range r.evalParams.ThreeInRowVariants {
		copied.ThreeInRowVariants[k] = v
	}
	return &copied
}

func (s pointAndValueSlice) Len() int {
	return len(s)
}

func (s pointAndValueSlice) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s pointAndValueSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
// play method for optimized robot player - iterative deepening with time management
func (r *optimizedRobotPlayer) play() (point, error) {
	r.nodeCount = 0 // Reset node count

	if r.count == 0 {
		p := point{maxLen / 2, maxLen / 2}
		r.set(p, r.pColor)
		return p, nil
	}

	// Quick win/defense checks
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

	// Quick checkmate search (only up to 4 steps to maintain speed)
	for i := 2; i <= 4; i += 2 {
		if p, ok := r.calculateKill(r.pColor, true, i); ok {
			return p, nil
		}
	}

	// Use time-controlled iterative deepening
	result := r.timeControlledIterativeDeepening()
	if result == nil {
		return point{}, errors.New("algorithm error")
	}

	r.set(result.p, r.pColor)
	return result.p, nil
}

// timeControlledIterativeDeepening implements iterative deepening with time management
// It starts with depth 4, then tries depth 6, but stops if depth 6 takes more than 60 seconds
func (r *optimizedRobotPlayer) timeControlledIterativeDeepening() *pointAndValue {
	var bestResult *pointAndValue
	startTime := time.Now()

	fmt.Printf("开始AI思考... ")

	// Phase 1: Quick search at depth 4 (should be very fast)
	fmt.Printf("深度4搜索... ")
	depth4Start := time.Now()
	result4 := r.optimizedIterativeDeepening(4)
	depth4Duration := time.Since(depth4Start)
	fmt.Printf("完成(%.3fs) ", depth4Duration.Seconds())

	if result4 != nil {
		bestResult = result4

		// If we found a winning move at depth 4, return immediately
		if result4.value > 800000 {
			fmt.Printf("发现胜负手!\n")
			return bestResult
		}
	}

	// Phase 2: Try deeper search at depth 6, but with 60-second timeout
	fmt.Printf("深度6搜索... ")
	depth6Start := time.Now()

	// Use context for proper cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resultChan := make(chan *pointAndValue, 1)

	go func() {
		result6 := r.optimizedIterativeDeepeningWithContext(ctx, 6)
		resultChan <- result6
	}()

	// Wait for either completion or timeout
	select {
	case result6 := <-resultChan:
		depth6Duration := time.Since(depth6Start)
		fmt.Printf("完成(%.3fs) ", depth6Duration.Seconds())
		if result6 != nil {
			bestResult = result6
		}
	case <-ctx.Done():
		fmt.Printf("超时(60s) ")
		// Context cancelled, goroutine should stop gracefully
	}

	totalDuration := time.Since(startTime)
	fmt.Printf("总用时: %.3fs\n", totalDuration.Seconds())

	return bestResult
}

// optimizedIterativeDeepening performs search up to the specified maximum depth
func (r *optimizedRobotPlayer) optimizedIterativeDeepening(maxDepth int) *pointAndValue {
	var bestResult *pointAndValue

	// Start with shallow search and progressively deepen
	for depth := 2; depth <= maxDepth; depth += 2 {
		result := r.optimizedMax(depth, -1000000000, 1000000000)

		if result != nil {
			bestResult = result
		}

		// Early termination for extremely strong positions (near-win)
		if bestResult != nil && bestResult.value > 1200000 {
			break
		}

		// Conservative early termination for very strong tactical wins
		if depth >= 4 && bestResult != nil && bestResult.value > 900000 {
			break
		}
	}

	return bestResult
}

// optimizedIterativeDeepeningWithContext performs search up to the specified maximum depth with context cancellation
func (r *optimizedRobotPlayer) optimizedIterativeDeepeningWithContext(ctx context.Context, maxDepth int) *pointAndValue {
	var bestResult *pointAndValue

	// Start with shallow search and progressively deepen
	for depth := 2; depth <= maxDepth; depth += 2 {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			return bestResult // Return best result found so far
		default:
			// Continue
		}

		result := r.optimizedMaxWithContext(ctx, depth, -1000000000, 1000000000)

		if result != nil {
			bestResult = result
		}

		// Early termination for extremely strong positions (near-win)
		if bestResult != nil && bestResult.value > 1200000 {
			break
		}

		// Conservative early termination for very strong tactical wins
		if depth >= 4 && bestResult != nil && bestResult.value > 900000 {
			break
		}
	}

	return bestResult
}

// getImprovedAdaptiveDepth returns adaptive search depth with better tactical awareness
func (r *optimizedRobotPlayer) getImprovedAdaptiveDepth() int {
	baseDepth := r.maxLevelCount

	// Check for complex tactical positions that require deeper analysis
	if r.hasComplexThreats() {
		return min(baseDepth+2, 8) // Deeper search for complex tactical positions, max 8
	}

	// Check for immediate threats that require deeper analysis
	if r.hasImmediateThreats() {
		return min(baseDepth+2, 8) // Deeper search for tactical positions
	}

	// In opening, use full depth for better positioning
	if r.count < 8 {
		return baseDepth
	}

	// In middle game with many pieces, use deeper search for tactics
	if r.count >= 8 && r.count < 25 {
		return min(baseDepth+2, 8) // Critical tactical phase
	}

	// In endgame, use full depth
	return baseDepth
}

// hasComplexThreats checks for complex tactical threats requiring deeper analysis
func (r *optimizedRobotPlayer) hasComplexThreats() bool {
	// Check for multiple live threes (potential double threats)
	myThreats := r.countLiveThreats(r.pColor)
	opponentThreats := r.countLiveThreats(r.pColor.conversion())

	// Complex position if either player has multiple live threes
	if myThreats >= 2 || opponentThreats >= 2 {
		return true
	}

	// Check for mixed threats (combination of threes and fours)
	if (myThreats >= 1 && r.exists4(r.pColor)) ||
		(opponentThreats >= 1 && r.exists4(r.pColor.conversion())) {
		return true
	}

	return false
}

// countLiveThreats counts live three-in-a-row threats for tactical analysis
func (r *optimizedRobotPlayer) countLiveThreats(color playerColor) int {
	threats := 0
	p := point{}

	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == colorEmpty {
				// Check if placing a piece here creates a live three
				r.set(p, color)

				for _, dir := range fourDirections {
					count := 1
					blocked := 0

					// Count in positive direction
					for k := 1; k < 4; k++ {
						pk := p.move(dir, k)
						if pk.checkRange() && r.get(pk) == color {
							count++
						} else if pk.checkRange() && r.get(pk) == color.conversion() {
							blocked++
							break
						} else {
							break
						}
					}

					// Count in negative direction
					for k := 1; k < 4; k++ {
						pk := p.move(dir, -k)
						if pk.checkRange() && r.get(pk) == color {
							count++
						} else if pk.checkRange() && r.get(pk) == color.conversion() {
							blocked++
							break
						} else {
							break
						}
					}

					// Live three: exactly 3 pieces and not blocked on both ends
					if count == 3 && blocked == 0 {
						threats++
						break // Only count once per position
					}
				}

				r.set(p, colorEmpty)
			}
		}
	}

	return threats
}

// optimizedMaxWithContext method for optimized robot player with context cancellation
func (r *optimizedRobotPlayer) optimizedMaxWithContext(ctx context.Context, step int, alpha, beta int) *pointAndValue {
	// Check if context was cancelled
	select {
	case <-ctx.Done():
		return nil
	default:
		// Continue
	}

	r.nodeCount++

	// Check cache first
	if cached := r.getCachedEvaluation(); cached != nil {
		return cached
	}

	candidates := r.getOptimizedCandidates(r.pColor)

	// Adaptive candidate pruning based on depth and game phase
	maxCandidates := r.getImprovedCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.set(p, r.pColor)
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.cacheEvaluation(result)
		return result
	}

	maxPoint := point{}
	maxVal := alpha

	for _, candidate := range candidates {
		// Check if context was cancelled during search
		select {
		case <-ctx.Done():
			// Timeout occurred, discard partial 6-layer result and use 4-layer result
			return nil
		default:
			// Continue
		}

		p := candidate.p
		r.set(p, r.pColor)

		// Quick evaluation for immediate wins
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal > 800000 {
			r.set(p, colorEmpty)
			result := &pointAndValue{p, boardVal}
			r.cacheEvaluation(result)
			return result
		}

		minResult := r.optimizedMinWithContext(ctx, step-1, maxVal, beta)
		if minResult == nil {
			r.set(p, colorEmpty)
			continue
		}
		evathis := minResult.value

		if evathis > maxVal {
			maxVal = evathis
			maxPoint = p
		}

		r.set(p, colorEmpty)

		// Alpha-beta pruning
		if maxVal >= beta {
			break
		}
	}

	result := &pointAndValue{maxPoint, maxVal}
	r.cacheEvaluation(result)
	return result
}

// optimizedMinWithContext method for optimized robot player with context cancellation
func (r *optimizedRobotPlayer) optimizedMinWithContext(ctx context.Context, step int, alpha, beta int) *pointAndValue {
	// Check if context was cancelled
	select {
	case <-ctx.Done():
		return nil
	default:
		// Continue
	}

	r.nodeCount++

	// Check cache first
	if cached := r.getCachedEvaluation(); cached != nil {
		return cached
	}

	candidates := r.getOptimizedCandidates(r.pColor.conversion())

	// Adaptive candidate pruning
	maxCandidates := r.getImprovedCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.set(p, r.pColor.conversion())
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.cacheEvaluation(result)
		return result
	}

	minPoint := point{}
	minVal := beta

	for _, candidate := range candidates {
		// Check if context was cancelled during search
		select {
		case <-ctx.Done():
			// Timeout occurred, discard partial 6-layer result and use 4-layer result
			return nil
		default:
			// Continue
		}

		p := candidate.p
		r.set(p, r.pColor.conversion())

		maxResult := r.optimizedMaxWithContext(ctx, step-1, alpha, minVal)
		if maxResult == nil {
			r.set(p, colorEmpty)
			continue
		}
		evathis := maxResult.value

		if evathis < minVal {
			minVal = evathis
			minPoint = p
		}

		r.set(p, colorEmpty)

		// Alpha-beta pruning
		if minVal <= alpha {
			break
		}
	}

	result := &pointAndValue{minPoint, minVal}
	r.cacheEvaluation(result)
	return result
}

// optimizedMax method for optimized robot player with better pruning
func (r *optimizedRobotPlayer) optimizedMax(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if cached := r.getCachedEvaluation(); cached != nil {
		return cached
	}

	candidates := r.getOptimizedCandidates(r.pColor)

	// Adaptive candidate pruning based on depth and game phase
	maxCandidates := r.getImprovedCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.set(p, r.pColor)
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.cacheEvaluation(result)
		return result
	}

	maxPoint := point{}
	maxVal := alpha

	for _, candidate := range candidates {
		p := candidate.p
		r.set(p, r.pColor)

		// Quick evaluation for immediate wins
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal > 800000 {
			r.set(p, colorEmpty)
			result := &pointAndValue{p, boardVal}
			r.cacheEvaluation(result)
			return result
		}

		minResult := r.optimizedMin(step-1, maxVal, beta)
		evathis := minResult.value

		if evathis > maxVal {
			maxVal = evathis
			maxPoint = p
		}

		r.set(p, colorEmpty)

		// Alpha-beta pruning
		if maxVal >= beta {
			break
		}
	}

	result := &pointAndValue{maxPoint, maxVal}
	r.cacheEvaluation(result)
	return result
}

// optimizedMin method for optimized robot player
func (r *optimizedRobotPlayer) optimizedMin(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if cached := r.getCachedEvaluation(); cached != nil {
		return cached
	}

	candidates := r.getOptimizedCandidates(r.pColor.conversion())

	// Adaptive candidate pruning
	maxCandidates := r.getImprovedCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.set(p, r.pColor.conversion())
		val := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.cacheEvaluation(result)
		return result
	}

	minPoint := point{}
	minVal := beta

	for _, candidate := range candidates {
		p := candidate.p
		r.set(p, r.pColor.conversion())

		maxResult := r.optimizedMax(step-1, alpha, minVal)
		evathis := maxResult.value

		if evathis < minVal {
			minVal = evathis
			minPoint = p
		}

		r.set(p, colorEmpty)

		// Alpha-beta pruning
		if minVal <= alpha {
			break
		}
	}

	result := &pointAndValue{minPoint, minVal}
	r.cacheEvaluation(result)
	return result
}

// getOptimizedCandidates gets candidate moves for optimized robot player with enhanced evaluation
func (r *optimizedRobotPlayer) getOptimizedCandidates(color playerColor) []*pointAndValue {
	var candidates []*pointAndValue
	p := point{}

	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == colorEmpty && r.isNeighbor(p) {
				// Use enhanced evaluation for better move ordering
				val := r.evaluatePoint(p, color)
				candidates = append(candidates, &pointAndValue{p, val})
			}
		}
	}

	// Sort candidates by value (best first) for better alpha-beta pruning
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].value > candidates[j].value
	})

	return candidates
}

// getImprovedCandidateLimit returns adaptive candidate limit with better tactical awareness
func (r *optimizedRobotPlayer) getImprovedCandidateLimit(totalCandidates, depth int) int {
	baseLimit := r.maxCountEachLevel

	// In tactical positions, consider more candidates to avoid missing key moves
	if r.hasComplexThreats() {
		baseLimit += 6
	}

	// For middle depths, use more candidates for better tactical analysis
	if depth >= 3 && depth <= 5 {
		baseLimit += 2
	}

	// Early game: check more positions for better opening play
	if r.count < 10 {
		baseLimit += 3
	}

	// Middle game: maintain high candidate count for tactical phases
	if r.count >= 10 && r.count < 25 {
		baseLimit += 2
	}

	// Ensure we don't exceed total candidates or go below minimum
	if baseLimit > totalCandidates {
		baseLimit = totalCandidates
	}
	if baseLimit < 12 {
		baseLimit = 12 // Higher minimum to avoid missing tactical moves
	}

	return baseLimit
}

// Simple caching methods for optimized robot player
func (r *optimizedRobotPlayer) getCachedEvaluation() *pointAndValue {
	// Disable broken cache that returns wrong point coordinates
	// The cache only stores evaluation values but not the actual best move points
	// This was causing the AI to always return (0,0)
	return nil
}

func (r *optimizedRobotPlayer) cacheEvaluation(result *pointAndValue) {
	if len(r.evalCache) > 8000 {
		// Clear cache when it gets too large
		r.evalCache = make(map[uint64]int)
	}
	r.evalCache[r.hash] = result.value
}

