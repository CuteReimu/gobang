package main

import (
	"fmt"
	"time"
)

// Simple benchmark without GUI dependencies
func main() {
	fmt.Println("Gobang AI 独立性能基准测试")
	fmt.Println("===========================")
	
	runStandaloneBenchmark()
}

func runStandaloneBenchmark() {
	// Test original AI
	fmt.Println("\n测试原始AI...")
	originalRobot := newRobotPlayer(colorBlack).(*robotPlayer)
	originalTime := measureAIStandalone(originalRobot)
	fmt.Printf("原始AI平均思考时间: %.3f 秒\n", originalTime)
	
	// Test optimized AI
	fmt.Println("\n测试优化AI...")
	optimizedRobot := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	optimizedTime := measureAIStandalone(optimizedRobot)
	fmt.Printf("优化AI平均思考时间: %.3f 秒\n", optimizedTime)
	
	// Calculate improvement
	improvement := ((originalTime - optimizedTime) / originalTime) * 100
	fmt.Printf("\n性能对比:\n")
	fmt.Printf("- 原始AI: %.3f 秒/步\n", originalTime)
	fmt.Printf("- 优化AI: %.3f 秒/步\n", optimizedTime)
	fmt.Printf("- 性能提升: %.1f%%\n", improvement)
	
	if improvement > 0 {
		fmt.Printf("- 速度倍数: %.1fx\n", originalTime/optimizedTime)
	}
	
	// Test strategic strength
	fmt.Println("\n测试策略强度...")
	testStrategicStrength()
}

func measureAIStandalone(robot *robotPlayer) float64 {
	// Setup mid-game position
	setupTestPosition(robot)
	
	totalTime := 0.0
	numTests := 5
	
	fmt.Printf("执行 %d 次测试...\n", numTests)
	
	for i := 0; i < numTests; i++ {
		// Reset to same position
		setupTestPosition(robot)
		
		start := time.Now()
		result := robot.max(robot.maxLevelCount, 100000000)
		elapsed := time.Since(start)
		
		totalTime += elapsed.Seconds()
		
		if result != nil {
			fmt.Printf("  测试 %d: %.3f 秒 - 推荐走法: (%d,%d) 评估值: %d\n", 
				i+1, elapsed.Seconds(), result.p.x, result.p.y, result.value)
		} else {
			fmt.Printf("  测试 %d: %.3f 秒 - 未找到走法\n", i+1, elapsed.Seconds())
		}
	}
	
	return totalTime / float64(numTests)
}

func setupTestPosition(robot *robotPlayer) {
	robot.initBoardStatus()
	center := maxLen / 2
	
	// Create a challenging mid-game position
	// Black pieces (AI)
	robot.set(point{center, center}, colorBlack)
	robot.set(point{center + 1, center}, colorBlack)
	robot.set(point{center - 1, center + 1}, colorBlack)
	robot.set(point{center + 2, center - 1}, colorBlack)
	robot.set(point{center, center - 2}, colorBlack)
	
	// White pieces (opponent) - creating threats
	robot.set(point{center, center + 1}, colorWhite)
	robot.set(point{center + 1, center + 1}, colorWhite)
	robot.set(point{center - 1, center}, colorWhite)
	robot.set(point{center + 1, center - 1}, colorWhite)
	robot.set(point{center - 2, center}, colorWhite)
	robot.set(point{center + 3, center}, colorWhite)
}

func testStrategicStrength() {
	// Compare strategic decisions between original and optimized AI
	fmt.Println("比较原始AI和优化AI的策略决策...")
	
	// Test position with clear best move
	original := newRobotPlayer(colorBlack).(*robotPlayer)
	optimized := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	
	// Create a position where AI should defend against immediate threat
	setupThreatPosition(original)
	setupThreatPosition(optimized)
	
	// Get moves from both AIs
	fmt.Println("\n威胁防御测试:")
	originalResult := original.max(original.maxLevelCount, 100000000)
	optimizedResult := optimized.max(optimized.maxLevelCount, 100000000)
	
	if originalResult != nil {
		fmt.Printf("原始AI选择: (%d,%d) 评估值: %d\n", 
			originalResult.p.x, originalResult.p.y, originalResult.value)
	}
	
	if optimizedResult != nil {
		fmt.Printf("优化AI选择: (%d,%d) 评估值: %d\n", 
			optimizedResult.p.x, optimizedResult.p.y, optimizedResult.value)
	}
	
	// Check if both chose to defend
	expectedDefensePoint := point{7, 6} // Should block opponent's threat
	
	fmt.Println("\n策略分析:")
	if originalResult != nil && optimizedResult != nil {
		if originalResult.p == optimizedResult.p {
			fmt.Println("✓ 两个AI选择了相同走法")
		} else {
			fmt.Println("✗ AI选择了不同走法")
		}
		
		if originalResult.p == expectedDefensePoint || optimizedResult.p == expectedDefensePoint {
			fmt.Println("✓ 至少一个AI选择了正确的防御走法")
		} else {
			fmt.Println("✗ 两个AI都未选择预期的防御走法")
		}
	}
}

func setupThreatPosition(robot *robotPlayer) {
	robot.initBoardStatus()
	
	// Create a position where opponent has 3 in a row and needs to be blocked
	// White has 3 in a row horizontally, AI must block at (7,6)
	robot.set(point{7, 7}, colorBlack)  // AI stone
	robot.set(point{6, 8}, colorBlack)  // AI stone
	
	robot.set(point{5, 6}, colorWhite)  // Opponent threat
	robot.set(point{6, 6}, colorWhite)  // Opponent threat  
	robot.set(point{8, 6}, colorWhite)  // Opponent threat
	// Position (7,6) should be blocked by AI
}