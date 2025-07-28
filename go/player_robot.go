package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
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
	rp := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            color,
		maxLevelCount:     4,  // Even depth for better minimax performance
		maxCountEachLevel: 12, // Reduced candidates for better performance
		maxCheckmateCount: 10, // Reduced checkmate search
		evalParams:        getOptimizedEvaluationParams(),
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

// newEnhancedRobotPlayer creates a robot player with 6-layer depth and algorithmic optimizations
func newEnhancedRobotPlayer(color playerColor) player {
	rp := &leanEnhancedRobotPlayer{
		robotPlayer: robotPlayer{
			boardCache:        make(boardCache),
			pColor:            color,
			maxLevelCount:     6,  // Maintain 6 layers as requested
			maxCountEachLevel: 12, // More aggressive candidate pruning
			maxCheckmateCount: 12, // Full checkmate search
			evalParams:        getEnhancedEvaluationParams(),
		},
		evalCache: make(map[uint64]int), // Simple evaluation cache
	}
	rp.initBoardStatus()
	return rp
}

// leanEnhancedRobotPlayer - simplified enhanced AI focused on real performance gains
type leanEnhancedRobotPlayer struct {
	robotPlayer
	evalCache map[uint64]int // Cache for position evaluations
	nodeCount int            // For debugging
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

// getEnhancedEvaluationParams returns enhanced evaluation parameters for 6-layer AI
func getEnhancedEvaluationParams() *EvaluationParams {
	return &EvaluationParams{
		LiveFour:             400000,  // Very high priority for winning moves
		DeadFourA:            320000,  // Strong threat detection
		DeadFourB:            300000,  // Strong threat detection
		DeadFourC:            280000,  // Strong threat detection
		LiveThreeNear:        2500,    // Enhanced three-in-a-row evaluation
		LiveThreeBonus:       9000,    // Strong bonus for good positions
		LiveThreeFar:         600,     // Better distant threat recognition
		DeadThree:            1100,    // Enhanced defensive evaluation
		DeadThreeBonus:       8500,    // Strong defensive bonus
		TwoCount2:            4000,    // Enhanced two-count evaluation
		TwoCount1:            3500,    // Enhanced single-two evaluation
		ScatterMultiplier:    8,       // Enhanced position evaluation
		OpponentPenalty:      700,     // Strong opponent threat response
		OpponentMinorPenalty: 400,     // Enhanced minor threat response
		FiveInRow:            1500000, // Maximum priority for wins
		FourInRowOpen:        400000,  // Maximum priority for winning threats
		FourInRowClosed:      35000,   // Enhanced closed-four evaluation
		ThreeInRowVariants: map[string]int{
			"open":   35000, // Enhanced open three evaluation
			"semi":   800,   // Enhanced semi-open evaluation
			"closed": 40000, // Enhanced closed three
			"gap":    1200,  // Enhanced gap pattern recognition
			"basic":  1000,  // Enhanced basic patterns
			"corner": 250,   // Enhanced corner evaluation
		},
	}
}

// enhancedRobotPlayer extends robotPlayer with advanced optimizations for 6-layer depth
type enhancedRobotPlayer struct {
	robotPlayer
	killerMoves      map[int][2]point // Killer moves for each depth level
	historyTable     map[point]int    // History heuristic table
	aspirationWindow int              // Aspiration window size
	nodeCount        int              // Node count for debugging
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
				evathis := r.evaluatePoint(p, r.pColor)
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
				evathis := r.evaluatePoint(p, r.pColor.conversion())
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

// Enhanced AI methods with algorithmic optimizations for 6-layer depth

// play method for enhanced robot player with optimizations
func (r *enhancedRobotPlayer) play() (point, error) {
	r.nodeCount = 0 // Reset node count

	if r.count == 0 {
		p := point{maxLen / 2, maxLen / 2}
		r.set(p, r.pColor)
		return p, nil
	}

	// Quick win/defense checks (same as original)
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

	// Enhanced checkmate calculation with reduced iterations for performance
	for i := 2; i <= r.maxCheckmateCount; i += 2 {
		if p, ok := r.calculateKill(r.pColor, true, i); ok {
			return p, nil
		}
	}

	// Use enhanced iterative deepening with aspiration windows
	result := r.enhancedIterativeDeepening()
	if result == nil {
		return point{}, errors.New("algorithm error")
	}

	// Update history table for this move
	r.historyTable[result.p]++

	r.set(result.p, r.pColor)
	return result.p, nil
}

// enhancedIterativeDeepening with aspiration windows and better time management
func (r *enhancedRobotPlayer) enhancedIterativeDeepening() *pointAndValue {
	var bestResult *pointAndValue
	alpha := -1000000000
	beta := 1000000000

	// Start with aspiration window around the previous best value
	if bestResult != nil {
		alpha = bestResult.value - r.aspirationWindow
		beta = bestResult.value + r.aspirationWindow
	}

	// Progressive deepening with aspiration windows
	for depth := 2; depth <= r.maxLevelCount; depth += 2 { // Even depths only
		result := r.enhancedMax(depth, alpha, beta)

		if result == nil {
			// Aspiration window failed, research with full window
			result = r.enhancedMax(depth, -1000000000, 1000000000)
		}

		if result != nil {
			bestResult = result

			// Adjust aspiration window for next iteration
			alpha = result.value - r.aspirationWindow
			beta = result.value + r.aspirationWindow
		}

		// Early termination for very strong positions
		if bestResult != nil && bestResult.value > 900000 {
			break
		}

		// If we find a good move and are beyond minimum depth, consider stopping
		if depth >= 4 && bestResult != nil && bestResult.value > 300000 {
			break
		}
	}

	return bestResult
}

// enhancedMax with improved alpha-beta pruning and move ordering
func (r *enhancedRobotPlayer) enhancedMax(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}

	// Generate and evaluate moves
	var queue pointAndValueSlice
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == 0 && r.isNeighbor(p) {
				evathis := r.evaluatePoint(p, r.pColor)

				// Apply history heuristic bonus
				if bonus, exists := r.historyTable[p]; exists {
					evathis += bonus * 10
				}

				// Apply killer move bonus
				if killers, exists := r.killerMoves[step]; exists {
					if p == killers[0] {
						evathis += 50000 // First killer
					} else if p == killers[1] {
						evathis += 25000 // Second killer
					}
				}

				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)

	// Adaptive candidate selection - more aggressive pruning
	maxCandidates := r.getEnhancedCandidateCount(len(queue), step)

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
	maxVal := alpha
	moveCount := 0

	for _, obj := range queue {
		if moveCount >= maxCandidates {
			break
		}
		moveCount++

		p = obj.p
		r.set(p, r.pColor)

		// Quick evaluation for early termination
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal > 800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		// Late Move Reduction (LMR) - search later moves with reduced depth
		if moveCount > 4 && step > 3 && boardVal < 100000 {
			// Search with reduced depth first
			reducedResult := r.enhancedMin(step-2, maxVal, beta)
			if reducedResult != nil && reducedResult.value <= maxVal {
				// If reduced search fails low, skip full search
				r.set(p, 0)
				continue
			}
		}

		evathisResult := r.enhancedMin(step-1, maxVal, beta)
		if evathisResult == nil {
			r.set(p, 0)
			continue
		}

		evalValue := evathisResult.value

		// Alpha-beta pruning
		if evalValue >= beta {
			r.set(p, 0)

			// Store killer move
			r.updateKillerMove(step, p)

			result := &pointAndValue{p, evalValue}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		if evalValue > maxVal {
			maxVal = evalValue
			maxPoint = p
		}

		r.set(p, 0)
	}

	if maxVal <= alpha {
		return nil
	}

	result := &pointAndValue{maxPoint, maxVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

// enhancedMin with improved alpha-beta pruning and move ordering
func (r *enhancedRobotPlayer) enhancedMin(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}

	// Generate and evaluate moves
	var queue pointAndValueSlice
	p := point{}
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p.x, p.y = j, i
			if r.get(p) == 0 && r.isNeighbor(p) {
				evathis := r.evaluatePoint(p, r.pColor.conversion())

				// Apply history heuristic bonus
				if bonus, exists := r.historyTable[p]; exists {
					evathis += bonus * 10
				}

				// Apply killer move bonus
				if killers, exists := r.killerMoves[step]; exists {
					if p == killers[0] {
						evathis += 50000 // First killer
					} else if p == killers[1] {
						evathis += 25000 // Second killer
					}
				}

				queue = append(queue, &pointAndValue{p, evathis})
			}
		}
	}
	sort.Sort(queue)

	// Adaptive candidate selection
	maxCandidates := r.getEnhancedCandidateCount(len(queue), step)

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
	minVal := beta
	moveCount := 0

	for _, obj := range queue {
		if moveCount >= maxCandidates {
			break
		}
		moveCount++

		p = obj.p
		r.set(p, r.pColor.conversion())

		// Quick evaluation for early termination
		boardVal := r.evaluateBoard(r.pColor) - r.evaluateBoard(r.pColor.conversion())
		if boardVal < -800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		// Late Move Reduction (LMR)
		if moveCount > 4 && step > 3 && boardVal > -100000 {
			// Search with reduced depth first
			reducedResult := r.enhancedMax(step-2, alpha, minVal)
			if reducedResult != nil && reducedResult.value >= minVal {
				// If reduced search fails high, skip full search
				r.set(p, 0)
				continue
			}
		}

		evathisResult := r.enhancedMax(step-1, alpha, minVal)
		if evathisResult == nil {
			r.set(p, 0)
			continue
		}

		evalValue := evathisResult.value

		// Alpha-beta pruning
		if evalValue <= alpha {
			r.set(p, 0)

			// Store killer move
			r.updateKillerMove(step, p)

			result := &pointAndValue{p, evalValue}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		if evalValue < minVal {
			minVal = evalValue
			minPoint = p
		}

		r.set(p, 0)
	}

	if minVal >= beta {
		return nil
	}

	result := &pointAndValue{minPoint, minVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

// getEnhancedCandidateCount returns more aggressive candidate pruning
func (r *enhancedRobotPlayer) getEnhancedCandidateCount(totalCandidates, depth int) int {
	baseCount := r.maxCountEachLevel

	// More aggressive pruning at deeper levels
	if depth <= 2 {
		baseCount += 4 // More candidates near leaf nodes
	} else if depth >= 4 {
		baseCount -= 2 // Fewer candidates at deeper levels
	}

	// Game phase adaptive adjustment
	if r.count < 8 {
		baseCount += 2 // Early game - more exploration
	} else if r.count > 20 {
		baseCount -= 3 // Late game - more focused search
	}

	return min(baseCount, totalCandidates)
}

// updateKillerMove updates the killer move table
func (r *enhancedRobotPlayer) updateKillerMove(depth int, move point) {
	killers := r.killerMoves[depth]

	// Shift killer moves
	if killers[0] != move {
		killers[1] = killers[0]
		killers[0] = move
		r.killerMoves[depth] = killers
	}
}

// Lean Enhanced AI methods with focused optimizations for 6-layer depth

// play method for lean enhanced robot player
func (r *leanEnhancedRobotPlayer) play() (point, error) {
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

	// Quick checkmate search with reduced iterations for performance
	for i := 2; i <= min(r.maxCheckmateCount, 8); i += 2 {
		if p, ok := r.calculateKill(r.pColor, true, i); ok {
			return p, nil
		}
	}

	// Use optimized iterative deepening
	result := r.optimizedIterativeDeepening()
	if result == nil {
		return point{}, errors.New("algorithm error")
	}

	r.set(result.p, r.pColor)
	return result.p, nil
}

// optimizedIterativeDeepening with better time management and early termination
func (r *leanEnhancedRobotPlayer) optimizedIterativeDeepening() *pointAndValue {
	var bestResult *pointAndValue

	// Start with shallow search and progressively deepen
	for depth := 2; depth <= r.maxLevelCount; depth += 2 {
		result := r.optimizedMax(depth, -1000000000, 1000000000)

		if result != nil {
			bestResult = result
		}

		// Early termination for very strong positions
		if bestResult != nil && bestResult.value > 800000 {
			break
		}

		// If we find a good move early and we're past minimum depth, consider stopping
		if depth >= 4 && bestResult != nil && bestResult.value > 300000 {
			break
		}

		// For mid-to-late game, if we have a decent move, don't spend too much time
		if r.count > 15 && depth >= 4 && bestResult != nil && bestResult.value > 50000 {
			break
		}
	}

	return bestResult
}

// optimizedMax with focused optimizations
func (r *leanEnhancedRobotPlayer) optimizedMax(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}

	// Generate candidates with aggressive pruning
	candidates := r.getOptimizedCandidates(r.pColor)

	// More aggressive candidate limit based on game phase
	maxCandidates := r.getAdaptiveCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.setIfEmpty(p, r.pColor)
		val := r.fastEvaluate(r.pColor)
		r.set(p, colorEmpty)
		result := &pointAndValue{p, val}
		r.putIntoCache(r.hash, step, result)
		return result
	}

	maxPoint := point{}
	maxVal := alpha

	for i, candidate := range candidates {
		if i >= maxCandidates {
			break
		}

		p := candidate.p
		r.set(p, r.pColor)

		// Quick evaluation for early termination
		boardVal := r.fastEvaluate(r.pColor)
		if boardVal > 800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		minResult := r.optimizedMin(step-1, maxVal, beta)
		if minResult == nil {
			r.set(p, 0)
			continue
		}

		evalValue := minResult.value

		// Alpha-beta pruning
		if evalValue >= beta {
			r.set(p, 0)
			result := &pointAndValue{p, evalValue}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		if evalValue > maxVal {
			maxVal = evalValue
			maxPoint = p
		}

		r.set(p, 0)
	}

	if maxVal <= alpha {
		return nil
	}

	result := &pointAndValue{maxPoint, maxVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

// optimizedMin with focused optimizations
func (r *leanEnhancedRobotPlayer) optimizedMin(step int, alpha, beta int) *pointAndValue {
	r.nodeCount++

	// Check cache first
	if v := r.getFromCache(r.hash, step); v != nil {
		return v
	}

	// Generate candidates with aggressive pruning
	candidates := r.getOptimizedCandidates(r.pColor.conversion())

	// More aggressive candidate limit
	maxCandidates := r.getAdaptiveCandidateLimit(len(candidates), step)
	if len(candidates) > maxCandidates {
		candidates = candidates[:maxCandidates]
	}

	if step == 1 {
		if len(candidates) == 0 {
			return nil
		}
		p := candidates[0].p
		r.setIfEmpty(p, r.pColor.conversion())
		val := r.fastEvaluate(r.pColor)
		r.set(p, 0)
		result := &pointAndValue{p, val}
		r.putIntoCache(r.hash, step, result)
		return result
	}

	var minPoint point
	minVal := beta

	for i, candidate := range candidates {
		if i >= maxCandidates {
			break
		}

		p := candidate.p
		r.set(p, r.pColor.conversion())

		// Quick evaluation for early termination
		boardVal := r.fastEvaluate(r.pColor)
		if boardVal < -800000 {
			r.set(p, 0)
			result := &pointAndValue{p, boardVal}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		maxResult := r.optimizedMax(step-1, alpha, minVal)
		if maxResult == nil {
			r.set(p, 0)
			continue
		}

		evalValue := maxResult.value

		// Alpha-beta pruning
		if evalValue <= alpha {
			r.set(p, 0)
			result := &pointAndValue{p, evalValue}
			r.putIntoCache(r.hash, step, result)
			return result
		}

		if evalValue < minVal {
			minVal = evalValue
			minPoint = p
		}

		r.set(p, 0)
	}

	if minVal >= beta {
		return nil
	}

	result := &pointAndValue{minPoint, minVal}
	r.putIntoCache(r.hash, step, result)
	return result
}

// getOptimizedCandidates returns better sorted candidates with aggressive pruning
func (r *leanEnhancedRobotPlayer) getOptimizedCandidates(color playerColor) []*pointAndValue {
	var candidates []*pointAndValue

	// Generate candidates only in areas of interest
	for i := 0; i < maxLen; i++ {
		for j := 0; j < maxLen; j++ {
			p := point{j, i}
			if r.get(p) == 0 && r.isNeighbor(p) {
				eval := r.evaluatePoint(p, color)
				candidates = append(candidates, &pointAndValue{p, eval})
			}
		}
	}

	// Sort by evaluation
	sort.Sort(pointAndValueSlice(candidates))

	return candidates
}

// getAdaptiveCandidateLimit returns adaptive candidate limits for better performance
func (r *leanEnhancedRobotPlayer) getAdaptiveCandidateLimit(totalCandidates, depth int) int {
	// Base limit is smaller than original for better performance
	baseLimit := r.maxCountEachLevel

	// Reduce candidates more aggressively at deeper levels
	if depth <= 2 {
		baseLimit += 2 // Slightly more at leaf levels
	} else if depth >= 4 {
		baseLimit -= 4 // Much fewer at deeper levels
	}

	// Game phase adaptation - more aggressive pruning in late game
	if r.count < 8 {
		baseLimit += 1 // Early game - slightly more exploration
	} else if r.count > 20 {
		baseLimit -= 5 // Late game - much more focused
	}

	// Never exceed total available candidates
	return min(baseLimit, totalCandidates)
}

// fastEvaluate - optimized evaluation with caching
func (r *leanEnhancedRobotPlayer) fastEvaluate(color playerColor) int {
	// Check cache first
	if cached, exists := r.evalCache[r.hash]; exists {
		return cached
	}

	// Compute evaluation
	myEval := r.evaluateBoard(color)
	oppEval := r.evaluateBoard(color.conversion())
	result := myEval - oppEval

	// Cache the result
	r.evalCache[r.hash] = result

	// Limit cache size to prevent memory issues
	if len(r.evalCache) > 10000 {
		// Clear cache when it gets too large
		r.evalCache = make(map[uint64]int)
	}

	return result
}
