package main

import (
	"fmt"
	"log"
	"time"
)

// benchmarkAI runs a simple benchmark to test AI performance
func benchmarkAI() {
	fmt.Println("运行AI性能基准测试...")
	fmt.Println("=====================================")

	// Test original AI
	fmt.Println("\n测试原始AI...")
	originalRobot := newRobotPlayer(colorBlack).(*robotPlayer)
	originalTime := measureAIThinkingTime(originalRobot)
	fmt.Printf("原始AI平均思考时间: %.2f 秒\n", originalTime)

	// Test optimized AI
	fmt.Println("\n测试优化AI...")
	optimizedRobot := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	optimizedTime := measureAIThinkingTime(optimizedRobot)
	fmt.Printf("优化AI平均思考时间: %.2f 秒\n", optimizedTime)

	// Test balanced AI
	fmt.Println("\n测试平衡AI...")
	balancedRobot := newBalancedRobotPlayer(colorBlack).(*robotPlayer)
	balancedTime := measureAIThinkingTime(balancedRobot)
	fmt.Printf("平衡AI平均思考时间: %.2f 秒\n", balancedTime)

	// Calculate improvements
	optimizedImprovement := ((originalTime - optimizedTime) / originalTime) * 100
	balancedImprovement := ((originalTime - balancedTime) / originalTime) * 100

	fmt.Printf("\n性能对比:\n")
	fmt.Printf("- 原始AI: %.2f 秒/步\n", originalTime)
	fmt.Printf("- 优化AI: %.2f 秒/步 (提升 %.1f%%)\n", optimizedTime, optimizedImprovement)
	fmt.Printf("- 平衡AI: %.2f 秒/步 (提升 %.1f%%)\n", balancedTime, balancedImprovement)

	if optimizedImprovement > 0 {
		fmt.Printf("- 优化AI快了 %.1f倍\n", originalTime/optimizedTime)
	}
	if balancedImprovement > 0 {
		fmt.Printf("- 平衡AI快了 %.1f倍\n", originalTime/balancedTime)
	}
}

// measureAIThinkingTime measures the average thinking time for the AI
func measureAIThinkingTime(robot *robotPlayer) float64 {
	// Simulate a mid-game position with some pieces on the board
	setupMidGamePosition(robot)

	totalTime := 0.0
	numTests := 3 // Reduced number for faster testing

	for i := 0; i < numTests; i++ {
		start := time.Now()

		// Make AI think about next move
		result := robot.max(robot.maxLevelCount, 100000000)

		elapsed := time.Since(start)
		totalTime += elapsed.Seconds()

		if result != nil {
			fmt.Printf("  测试 %d: %.2f 秒 (走法: %v)\n", i+1, elapsed.Seconds(), result.p)
		} else {
			fmt.Printf("  测试 %d: %.2f 秒 (未找到走法)\n", i+1, elapsed.Seconds())
		}
	}

	return totalTime / float64(numTests)
}

// setupMidGamePosition sets up a typical mid-game position for testing
func setupMidGamePosition(robot *robotPlayer) {
	// Clear the board first
	robot.initBoardStatus()

	center := maxLen / 2

	// Set up a typical mid-game scenario
	// Black pieces (robot's color)
	robot.set(point{center, center}, colorBlack)
	robot.set(point{center + 1, center}, colorBlack)
	robot.set(point{center - 1, center + 1}, colorBlack)
	robot.set(point{center + 2, center - 1}, colorBlack)

	// White pieces (opponent)
	robot.set(point{center, center + 1}, colorWhite)
	robot.set(point{center + 1, center + 1}, colorWhite)
	robot.set(point{center - 1, center}, colorWhite)
	robot.set(point{center + 1, center - 1}, colorWhite)
	robot.set(point{center - 2, center}, colorWhite)
}

// runSelfPlayTest runs a self-play test to evaluate parameter effectiveness
func runSelfPlayTest() {
	fmt.Println("\n运行自对弈测试...")
	fmt.Println("=========================")

	// Create two AIs with different parameters
	original := newRobotPlayer(colorBlack).(*robotPlayer)
	optimized := newOptimizedRobotPlayer(colorWhite).(*robotPlayer)

	// Simulate a quick game (limited depth for speed)
	original.maxLevelCount = 3
	optimized.maxLevelCount = 3

	winner, moves := simulateGame(original, optimized)

	fmt.Printf("游戏在 %d 步内完成\n", moves)
	if winner == colorBlack {
		fmt.Println("获胜者: 原始AI (黑)")
	} else if winner == colorWhite {
		fmt.Println("获胜者: 优化AI (白)")
	} else {
		fmt.Println("游戏平局")
	}
}

// simulateGame simulates a game between two AIs
func simulateGame(player1, player2 *robotPlayer) (playerColor, int) {
	// Create a shared board
	board := make([][]playerColor, maxLen)
	for i := 0; i < maxLen; i++ {
		board[i] = make([]playerColor, maxLen)
	}

	// Initialize both players with the same board state
	player1.initBoardStatus()
	player2.initBoardStatus()

	moves := 0
	maxMoves := 50 // Limit game length for testing
	currentPlayer := player1

	for moves < maxMoves {
		var p point
		var err error

		// Get move from current player
		if currentPlayer == player1 {
			p, err = player1.play()
		} else {
			p, err = player2.play()
		}

		if err != nil {
			log.Printf("Error getting move: %v", err)
			break
		}

		// Update shared board
		board[p.y][p.x] = currentPlayer.pColor
		moves++

		// Update both players' internal boards
		player1.set(p, currentPlayer.pColor)
		player2.set(p, currentPlayer.pColor)

		// Check for win
		if checkWin(board, p) {
			return currentPlayer.pColor, moves
		}

		// Switch players
		if currentPlayer == player1 {
			currentPlayer = player2
		} else {
			currentPlayer = player1
		}
	}

	return colorEmpty, moves // Draw or max moves reached
}

// checkWin checks if the last move resulted in a win
func checkWin(board [][]playerColor, lastMove point) bool {
	color := board[lastMove.y][lastMove.x]

	for _, dir := range fourDirections {
		count := 0
		for i := -4; i <= 4; i++ {
			p := lastMove.move(dir, i)
			if p.checkRange() && board[p.y][p.x] == color {
				count++
				if count == 5 {
					return true
				}
			} else {
				count = 0
			}
		}
	}

	return false
}

// setupTestBoard creates a test board scenario
func setupTestBoard(ai player, moveCount int) {
	// Get the underlying robot player interface
	var rp *robotPlayer
	switch v := ai.(type) {
	case *robotPlayer:
		rp = v
	case *optimizedRobotPlayer:
		rp = &v.robotPlayer
	default:
		return
	}

	// Set up a realistic game position
	positions := []point{
		{7, 7}, {6, 8}, {7, 6}, {7, 8}, {8, 8}, {6, 6},
		{6, 7}, {8, 7}, {5, 8}, {8, 5}, {6, 9}, {7, 10},
		{3, 10}, {4, 9}, {7, 9}, {8, 9}, {6, 10}, {5, 11},
		{10, 6}, {9, 7}, {6, 11}, {6, 12}, {4, 8}, {5, 9},
	}

	for i := 0; i < moveCount && i < len(positions); i++ {
		color := colorBlack
		if i%2 == 1 {
			color = colorWhite
		}
		rp.set(positions[i], color)
	}
}
