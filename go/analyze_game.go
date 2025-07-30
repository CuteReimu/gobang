package main

import (
	"fmt"
	"time"
)

// Analyze the specific game where optimized AI lost
func analyzeUserGame() {
	fmt.Println("分析用户对局中AI的表现")
	fmt.Println("==========================")
	
	// Recreate the game state at various critical points
	testGamePosition()
}

func testGamePosition() {
	fmt.Println("\n--- 分析用户提供的对局 ---")
	
	// Create different AI types for comparison
	originalAI := newRobotPlayer(colorBlack).(*robotPlayer)
	optimizedAI := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	balancedAI := newBalancedRobotPlayer(colorBlack).(*robotPlayer)
	
	// Replay the game moves up to a critical point
	moves := []struct{
		color playerColor
		p point
	}{
		{colorBlack, point{7,7}},   // 黑(7,7)
		{colorWhite, point{6,8}},   // 白(6,8)
		{colorBlack, point{7,8}},   // 黑(7,8)
		{colorWhite, point{7,9}},   // 白(7,9)
		{colorBlack, point{5,7}},   // 黑(5,7)
		{colorWhite, point{8,7}},   // 白(8,7)
		{colorBlack, point{6,6}},   // 黑(6,6)
		{colorWhite, point{7,5}},   // 白(7,5)
		{colorBlack, point{5,5}},   // 黑(5,5)
		{colorWhite, point{8,8}},   // 白(8,8)
		{colorBlack, point{5,4}},   // 黑(5,4)
		{colorWhite, point{5,6}},   // 白(5,6)
		{colorBlack, point{6,4}},   // 黑(6,4)
		{colorWhite, point{4,4}},   // 白(4,4) - up to this point
	}
	
	// Apply moves to all AI instances
	for _, move := range moves {
		originalAI.set(move.p, move.color)
		optimizedAI.set(move.p, move.color)
		balancedAI.set(move.p, move.color)
	}
	
	fmt.Println("当前局面（前14步后）:")
	printBoardState(originalAI)
	
	// Now test what each AI would do next
	fmt.Println("\n各AI的下一步选择:")
	
	// Test original AI
	start := time.Now()
	originalMove, err := originalAI.play()
	originalTime := time.Since(start)
	if err == nil {
		originalEval := originalAI.evaluateBoard(colorBlack)
		fmt.Printf("原始AI: %v (评估: %d, 用时: %v)\n", originalMove, originalEval, originalTime)
	}
	
	// Test optimized AI  
	start = time.Now()
	optimizedMove, err := optimizedAI.play()
	optimizedTime := time.Since(start)
	if err == nil {
		optimizedEval := optimizedAI.evaluateBoard(colorBlack)
		fmt.Printf("优化AI: %v (评估: %d, 用时: %v)\n", optimizedMove, optimizedEval, optimizedTime)
	}
	
	// Test balanced AI
	start = time.Now()
	balancedMove, err := balancedAI.play()
	balancedTime := time.Since(start)
	if err == nil {
		balancedEval := balancedAI.evaluateBoard(colorBlack)
		fmt.Printf("平衡AI: %v (评估: %d, 用时: %v)\n", balancedMove, balancedEval, balancedTime)
	}
	
	// Analyze the position evaluation differences
	fmt.Println("\n位置评估比较:")
	originalEval := originalAI.evaluateBoard(colorBlack)
	optimizedEval := optimizedAI.evaluateBoard(colorBlack)
	balancedEval := balancedAI.evaluateBoard(colorBlack)
	
	fmt.Printf("原始AI评估: %d\n", originalEval)
	fmt.Printf("优化AI评估: %d\n", optimizedEval)
	fmt.Printf("平衡AI评估: %d\n", balancedEval)
	
	// Test threat recognition
	fmt.Println("\n威胁识别测试:")
	testThreatRecognition(originalAI, "原始AI")
	testThreatRecognition(optimizedAI, "优化AI")
	testThreatRecognition(balancedAI, "平衡AI")
}

func testThreatRecognition(ai *robotPlayer, name string) {
	// Test if AI can find immediate wins
	winMove, hasWin := ai.findForm5(colorBlack)
	if hasWin {
		fmt.Printf("%s找到获胜走法: %v\n", name, winMove)
	}
	
	// Test if AI recognizes threats to block
	blockMove, hasBlock := ai.stop4(colorBlack)
	if hasBlock {
		fmt.Printf("%s识别需要阻挡的威胁: %v\n", name, blockMove)
	}
	
	// Test checkmate calculation
	for i := 2; i <= 4; i += 2 {
		if killMove, hasKill := ai.calculateKill(colorBlack, true, i); hasKill {
			fmt.Printf("%s找到%d步杀: %v\n", name, i, killMove)
			break
		}
	}
}

func printBoardState(ai *robotPlayer) {
	fmt.Println("  0 1 2 3 4 5 6 7 8 9 A B C D E")
	for y := 0; y < 15; y++ {
		fmt.Printf("%X ", y)
		for x := 0; x < 15; x++ {
			p := point{x, y}
			if ai.at(p) == colorBlack {
				fmt.Print("● ")
			} else if ai.at(p) == colorWhite {
				fmt.Print("○ ")
			} else {
				fmt.Print("· ")
			}
		}
		fmt.Println()
	}
}