package main

import (
	"fmt"
	"log"
	"time"
)

// benchmarkAI runs a simple benchmark to test AI performance
func benchmarkAI() {
	fmt.Println("Running AI Performance Benchmark...")
	fmt.Println("=====================================")
	
	// Test original AI
	fmt.Println("\nTesting Original AI...")
	originalRobot := newRobotPlayer(colorBlack).(*robotPlayer)
	originalTime := measureAIThinkingTime(originalRobot)
	fmt.Printf("Original AI average thinking time: %.2f seconds\n", originalTime)
	
	// Test optimized AI
	fmt.Println("\nTesting Optimized AI...")
	optimizedRobot := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	optimizedTime := measureAIThinkingTime(optimizedRobot)
	fmt.Printf("Optimized AI average thinking time: %.2f seconds\n", optimizedTime)
	
	// Calculate improvement
	improvement := ((originalTime - optimizedTime) / originalTime) * 100
	fmt.Printf("\nPerformance improvement: %.1f%%\n", improvement)
	
	if improvement > 0 {
		fmt.Printf("Optimized AI is %.1fx faster\n", originalTime/optimizedTime)
	} else {
		fmt.Println("Optimized AI is not faster (may need more tuning)")
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
			fmt.Printf("  Test %d: %.2f seconds (move: %v)\n", i+1, elapsed.Seconds(), result.p)
		} else {
			fmt.Printf("  Test %d: %.2f seconds (no move found)\n", i+1, elapsed.Seconds())
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
	fmt.Println("\nRunning Self-Play Test...")
	fmt.Println("=========================")
	
	// Create two AIs with different parameters
	original := newRobotPlayer(colorBlack).(*robotPlayer)
	optimized := newOptimizedRobotPlayer(colorWhite).(*robotPlayer)
	
	// Simulate a quick game (limited depth for speed)
	original.maxLevelCount = 3
	optimized.maxLevelCount = 3
	
	winner, moves := simulateGame(original, optimized)
	
	fmt.Printf("Game completed in %d moves\n", moves)
	if winner == colorBlack {
		fmt.Println("Winner: Original AI (Black)")
	} else if winner == colorWhite {
		fmt.Println("Winner: Optimized AI (White)")
	} else {
		fmt.Println("Game ended in a draw")
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