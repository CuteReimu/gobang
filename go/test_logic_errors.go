package main

import (
	"fmt"
	"time"
)

func testLogicErrors() {
	fmt.Println("AI逻辑错误检测")
	fmt.Println("================")

	// Test 1: Check if AI can find obvious winning moves
	testObviousWinningMoves()

	// Test 2: Check if AI can block obvious threats
	testObviousThreats()

	// Test 3: Check if AI can recognize forced wins
	testForcedWins()

	// Test 4: Check if AI can avoid obvious blunders
	testObviousBlunders()

	// Test 5: Check evaluation function consistency
	testEvaluationConsistency()
}

func testObviousWinningMoves() {
	fmt.Println("\n--- 测试明显获胜走法 ---")

	// Test position where AI has 4 in a row and can win
	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Setup: Black has 4 in a row, can win in 1 move
	// Position: Black at (7,7), (8,7), (9,7), (10,7) - horizontal line
	ai.set(point{7, 7}, colorBlack)
	ai.set(point{8, 7}, colorBlack)
	ai.set(point{9, 7}, colorBlack)
	ai.set(point{10, 7}, colorBlack)

	// White has some pieces to make it realistic
	ai.set(point{7, 8}, colorWhite)
	ai.set(point{8, 8}, colorWhite)

	start := time.Now()
	move, err := ai.play()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// Expected move should be (11,7) or (6,7) to complete 5 in a row
	expected1 := point{11, 7}
	expected2 := point{6, 7}

	if move == expected1 || move == expected2 {
		fmt.Printf("✓ AI正确找到获胜走法: %v (用时: %v)\n", move, duration)
	} else {
		fmt.Printf("✗ AI未找到明显获胜走法: %v (期望: %v 或 %v)\n", move, expected1, expected2)
	}
}

func testObviousThreats() {
	fmt.Println("\n--- 测试明显威胁阻挡 ---")

	// Test position where opponent has 4 in a row and AI must block
	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Setup: White has 4 in a row, AI must block
	// Position: White at (7,7), (8,7), (9,7), (10,7) - horizontal line
	ai.set(point{7, 7}, colorWhite)
	ai.set(point{8, 7}, colorWhite)
	ai.set(point{9, 7}, colorWhite)
	ai.set(point{10, 7}, colorWhite)

	// Black has some pieces
	ai.set(point{7, 8}, colorBlack)
	ai.set(point{8, 8}, colorBlack)

	start := time.Now()
	move, err := ai.play()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// Expected move should be (11,7) or (6,7) to block
	expected1 := point{11, 7}
	expected2 := point{6, 7}

	if move == expected1 || move == expected2 {
		fmt.Printf("✓ AI正确阻挡威胁: %v (用时: %v)\n", move, duration)
	} else {
		fmt.Printf("✗ AI未阻挡明显威胁: %v (期望: %v 或 %v)\n", move, expected1, expected2)
	}
}

func testForcedWins() {
	fmt.Println("\n--- 测试强制获胜 ---")

	// Test position where AI has multiple threats and can force a win
	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Setup: Black has two live threes that can't both be blocked
	// This is a classic forced win position
	ai.set(point{7, 7}, colorBlack)
	ai.set(point{8, 7}, colorBlack)
	ai.set(point{9, 7}, colorBlack)

	ai.set(point{7, 8}, colorBlack)
	ai.set(point{8, 8}, colorBlack)
	ai.set(point{9, 8}, colorBlack)

	// White has some pieces
	ai.set(point{6, 7}, colorWhite)
	ai.set(point{10, 7}, colorWhite)

	start := time.Now()
	move, err := ai.play()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// AI should find a move that creates multiple threats
	fmt.Printf("AI选择: %v (用时: %v)\n", move, duration)

	// Check if the move creates a winning threat
	ai.set(move, colorBlack)
	winMove, hasWin := ai.findForm5(colorBlack)
	if hasWin {
		fmt.Printf("✓ AI找到获胜走法: %v\n", winMove)
	} else {
		fmt.Printf("? AI走法未直接获胜，需要进一步分析\n")
	}
}

func testObviousBlunders() {
	fmt.Println("\n--- 测试明显失误 ---")

	// Test position where AI should not make obvious blunders
	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Setup: White has 3 in a row, AI should block
	ai.set(point{7, 7}, colorWhite)
	ai.set(point{8, 7}, colorWhite)
	ai.set(point{9, 7}, colorWhite)

	// Black has some pieces
	ai.set(point{7, 8}, colorBlack)
	ai.set(point{8, 8}, colorBlack)

	start := time.Now()
	move, err := ai.play()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// AI should block the threat
	expected1 := point{10, 7}
	expected2 := point{6, 7}

	if move == expected1 || move == expected2 {
		fmt.Printf("✓ AI正确阻挡威胁: %v (用时: %v)\n", move, duration)
	} else {
		fmt.Printf("✗ AI可能犯明显失误: %v (期望: %v 或 %v)\n", move, expected1, expected2)
	}
}

func testEvaluationConsistency() {
	fmt.Println("\n--- 测试评估函数一致性 ---")

	// Test if evaluation function gives consistent results
	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Setup a simple position
	ai.set(point{7, 7}, colorBlack)
	ai.set(point{8, 7}, colorBlack)
	ai.set(point{9, 7}, colorBlack)

	ai.set(point{7, 8}, colorWhite)
	ai.set(point{8, 8}, colorWhite)

	// Evaluate the same position multiple times
	eval1 := ai.evaluateBoard(colorBlack)
	eval2 := ai.evaluateBoard(colorBlack)
	eval3 := ai.evaluateBoard(colorBlack)

	if eval1 == eval2 && eval2 == eval3 {
		fmt.Printf("✓ 评估函数一致性良好: %d\n", eval1)
	} else {
		fmt.Printf("✗ 评估函数不一致: %d, %d, %d\n", eval1, eval2, eval3)
	}

	// Test if evaluation changes appropriately when adding pieces
	ai.set(point{10, 7}, colorBlack)
	evalAfter := ai.evaluateBoard(colorBlack)

	if evalAfter > eval1 {
		fmt.Printf("✓ 评估函数正确响应棋盘变化: %d -> %d\n", eval1, evalAfter)
	} else {
		fmt.Printf("✗ 评估函数未正确响应棋盘变化: %d -> %d\n", eval1, evalAfter)
	}
}

// Test specific pattern recognition
func testPatternRecognition() {
	fmt.Println("\n--- 测试模式识别 ---")

	ai := newRobotPlayer(colorBlack).(*robotPlayer)

	// Test live four pattern
	ai.set(point{7, 7}, colorBlack)
	ai.set(point{8, 7}, colorBlack)
	ai.set(point{9, 7}, colorBlack)
	ai.set(point{10, 7}, colorBlack)

	// Check if AI recognizes the live four
	winMove, hasWin := ai.findForm5(colorBlack)
	if hasWin {
		fmt.Printf("✓ AI正确识别活四模式: %v\n", winMove)
	} else {
		fmt.Printf("✗ AI未识别活四模式\n")
	}

	// Test if AI can block opponent's live four
	ai.initBoardStatus()
	ai.set(point{7, 7}, colorWhite)
	ai.set(point{8, 7}, colorWhite)
	ai.set(point{9, 7}, colorWhite)
	ai.set(point{10, 7}, colorWhite)

	blockMove, hasBlock := ai.stop4(colorBlack)
	if hasBlock {
		fmt.Printf("✓ AI正确识别需要阻挡的活四: %v\n", blockMove)
	} else {
		fmt.Printf("✗ AI未识别需要阻挡的活四\n")
	}
}
