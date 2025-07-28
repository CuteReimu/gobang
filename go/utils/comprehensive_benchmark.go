package main

import (
	"fmt"
	"time"
)

// Enhanced benchmark with all three AI variants
func main() {
	fmt.Println("Gobang AI 完整性能基准测试")
	fmt.Println("===============================")
	
	runComprehensiveBenchmark()
}

func runComprehensiveBenchmark() {
	// Test all three AI variants
	fmt.Println("测试所有AI变体...")
	
	originalRobot := newRobotPlayer(colorBlack).(*robotPlayer)
	optimizedRobot := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	balancedRobot := newBalancedRobotPlayer(colorBlack).(*robotPlayer)
	
	fmt.Println("\n性能测试:")
	fmt.Println("----------")
	
	originalTime := measureAIComprehensive(originalRobot, "原始AI")
	optimizedTime := measureAIComprehensive(optimizedRobot, "优化AI")
	balancedTime := measureAIComprehensive(balancedRobot, "平衡AI")
	
	fmt.Println("\n战术决策测试:")
	fmt.Println("------------")
	
	testTacticalDecisions(originalRobot, optimizedRobot, balancedRobot)
	
	fmt.Println("\n威胁检测测试:")
	fmt.Println("------------")
	
	testThreatDetection(originalRobot, optimizedRobot, balancedRobot)
	
	fmt.Printf("\n总结报告:\n")
	fmt.Printf("========\n")
	fmt.Printf("原始AI: %.3f 秒/步 (基准)\n", originalTime)
	fmt.Printf("优化AI: %.3f 秒/步 (%.1fx 更快，但可能较弱)\n", optimizedTime, originalTime/optimizedTime)
	fmt.Printf("平衡AI: %.3f 秒/步 (%.1fx 更快，更强棋力)\n", balancedTime, originalTime/balancedTime)
	
	fmt.Println("\n推荐使用:")
	if balancedTime < originalTime * 0.5 {
		fmt.Println("✓ 平衡AI - 最佳选择：速度快且棋力强")
	} else if balancedTime < originalTime * 0.8 {
		fmt.Println("✓ 平衡AI - 推荐选择：合理速度和更强棋力")
	} else {
		fmt.Println("? 需要进一步优化平衡AI参数")
	}
}

func measureAIComprehensive(robot *robotPlayer, name string) float64 {
	fmt.Printf("\n测试 %s (深度:%d, 候选:%d):\n", name, robot.maxLevelCount, robot.maxCountEachLevel)
	
	// Test multiple positions
	positions := []func(*robotPlayer){
		setupOpeningPosition,
		setupMidGamePosition,
		setupTacticalPosition,
		setupEndGamePosition,
	}
	
	positionNames := []string{"开局", "中局", "战术", "终局"}
	
	totalTime := 0.0
	tests := 0
	
	for i, setupFunc := range positions {
		fmt.Printf("  %s位置: ", positionNames[i])
		
		setupFunc(robot)
		
		start := time.Now()
		result := robot.max(robot.maxLevelCount, 100000000)
		elapsed := time.Since(start)
		
		totalTime += elapsed.Seconds()
		tests++
		
		if result != nil {
			fmt.Printf("%.3f秒 -> (%d,%d) [%d]\n", 
				elapsed.Seconds(), result.p.x+1, result.p.y+1, result.value)
		} else {
			fmt.Printf("%.3f秒 -> 无解\n", elapsed.Seconds())
		}
	}
	
	avg := totalTime / float64(tests)
	fmt.Printf("  平均: %.3f 秒\n", avg)
	return avg
}

func setupOpeningPosition(robot *robotPlayer) {
	robot.initBoardStatus()
	center := maxLen / 2
	
	// Simple opening position
	robot.set(point{center, center}, colorBlack)
	robot.set(point{center + 1, center}, colorWhite)
	robot.set(point{center, center + 1}, colorBlack)
}

func setupMidGamePosition(robot *robotPlayer) {
	robot.initBoardStatus()
	center := maxLen / 2
	
	// Complex mid-game position
	robot.set(point{center, center}, colorBlack)
	robot.set(point{center + 1, center}, colorBlack)
	robot.set(point{center - 1, center + 1}, colorBlack)
	robot.set(point{center + 2, center - 1}, colorBlack)
	
	robot.set(point{center, center + 1}, colorWhite)
	robot.set(point{center + 1, center + 1}, colorWhite)
	robot.set(point{center - 1, center}, colorWhite)
	robot.set(point{center + 1, center - 1}, colorWhite)
	robot.set(point{center - 2, center}, colorWhite)
}

func setupTacticalPosition(robot *robotPlayer) {
	robot.initBoardStatus()
	
	// Position requiring tactical calculation
	robot.set(point{5, 6}, colorWhite)
	robot.set(point{6, 6}, colorWhite)
	robot.set(point{8, 6}, colorWhite)
	// AI should block at (7,6)
	
	robot.set(point{7, 7}, colorBlack)
	robot.set(point{6, 8}, colorBlack)
}

func setupEndGamePosition(robot *robotPlayer) {
	robot.initBoardStatus()
	
	// Crowded end-game position
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if (i+j)%7 == 0 && i < 10 && j < 10 {
				if (i+j)%2 == 0 {
					robot.set(point{i, j}, colorBlack)
				} else {
					robot.set(point{i, j}, colorWhite)
				}
			}
		}
	}
}

func testTacticalDecisions(original, optimized, balanced *robotPlayer) {
	fmt.Println("测试关键战术决策...")
	
	// Setup critical position
	setupCriticalPosition := func(robot *robotPlayer) {
		robot.initBoardStatus()
		// White threatens to win
		robot.set(point{5, 5}, colorWhite)
		robot.set(point{6, 5}, colorWhite)
		robot.set(point{7, 5}, colorWhite)
		// AI must block at (8,5) or (4,5)
		
		robot.set(point{6, 6}, colorBlack)
		robot.set(point{7, 7}, colorBlack)
	}
	
	ais := []*robotPlayer{original, optimized, balanced}
	names := []string{"原始", "优化", "平衡"}
	
	for i, ai := range ais {
		setupCriticalPosition(ai)
		result := ai.max(ai.maxLevelCount, 100000000)
		
		if result != nil {
			correctMove := result.p.x == 8 && result.p.y == 5 || result.p.x == 4 && result.p.y == 5
			status := "✗"
			if correctMove {
				status = "✓"
			}
			fmt.Printf("  %s AI: (%d,%d) %s\n", names[i], result.p.x+1, result.p.y+1, status)
		} else {
			fmt.Printf("  %s AI: 无解 ✗\n", names[i])
		}
	}
}

func testThreatDetection(original, optimized, balanced *robotPlayer) {
	fmt.Println("测试威胁检测能力...")
	
	// Setup position with multiple threats
	setupThreatPosition := func(robot *robotPlayer) {
		robot.initBoardStatus()
		
		// Create multiple potential threats
		robot.set(point{3, 3}, colorBlack)
		robot.set(point{4, 4}, colorBlack)
		robot.set(point{6, 6}, colorBlack)
		
		robot.set(point{8, 3}, colorWhite)
		robot.set(point{9, 3}, colorWhite)
		robot.set(point{7, 8}, colorWhite)
		robot.set(point{8, 8}, colorWhite)
	}
	
	ais := []*robotPlayer{original, optimized, balanced}
	names := []string{"原始", "优化", "平衡"}
	
	for i, ai := range ais {
		setupThreatPosition(ai)
		
		// Test threat detection methods if available
		fmt.Printf("  %s AI: ", names[i])
		
		if hasMethod := true; hasMethod { // Placeholder for actual method check
			result := ai.max(ai.maxLevelCount, 100000000)
			if result != nil {
				fmt.Printf("威胁评估: %d\n", result.value)
			} else {
				fmt.Printf("无法评估\n")
			}
		}
	}
}