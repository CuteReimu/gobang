package main

import (
	"fmt"
	"time"
)

func testPerformance() {
	fmt.Println("AI性能优化验证测试")
	fmt.Println("====================")

	// Test all AI variants
	ais := []struct {
		name string
		ai   player
	}{
		{"原始AI", newRobotPlayer(colorBlack)},
		{"优化AI", newOptimizedRobotPlayer(colorBlack)},
		{"平衡AI", newBalancedRobotPlayer(colorBlack)},
		{"增强AI", newEnhancedRobotPlayer(colorBlack)},
	}

	// Test positions
	positions := []struct {
		name  string
		setup func(player)
	}{
		{"开局", setupOpeningPosition},
		{"中局", setupMidGamePositionTest},
		{"终局", setupEndGamePosition},
	}

	for _, pos := range positions {
		fmt.Printf("\n--- %s测试 ---\n", pos.name)

		for _, aiTest := range ais {
			pos.setup(aiTest.ai)

			start := time.Now()
			move, err := aiTest.ai.play()
			duration := time.Since(start)

			if err != nil {
				fmt.Printf("  %s: 错误 - %v\n", aiTest.name, err)
			} else {
				fmt.Printf("  %s: %v 用时: %v\n", aiTest.name, move, duration)
			}
		}
	}
}

func setupOpeningPosition(ai player) {
	// Simple opening position
	center := maxLen / 2

	// Handle different AI types
	switch v := ai.(type) {
	case *robotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorWhite)
		v.set(point{center, center + 1}, colorBlack)
	case *leanEnhancedRobotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorWhite)
		v.set(point{center, center + 1}, colorBlack)
	case *enhancedRobotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorWhite)
		v.set(point{center, center + 1}, colorBlack)
	}
}

func setupMidGamePositionTest(ai player) {
	// Complex mid-game position
	center := maxLen / 2

	// Handle different AI types
	switch v := ai.(type) {
	case *robotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorBlack)
		v.set(point{center - 1, center + 1}, colorBlack)
		v.set(point{center + 2, center - 1}, colorBlack)
		v.set(point{center, center + 1}, colorWhite)
		v.set(point{center + 1, center + 1}, colorWhite)
		v.set(point{center - 1, center}, colorWhite)
		v.set(point{center + 1, center - 1}, colorWhite)
		v.set(point{center - 2, center}, colorWhite)
	case *leanEnhancedRobotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorBlack)
		v.set(point{center - 1, center + 1}, colorBlack)
		v.set(point{center + 2, center - 1}, colorBlack)
		v.set(point{center, center + 1}, colorWhite)
		v.set(point{center + 1, center + 1}, colorWhite)
		v.set(point{center - 1, center}, colorWhite)
		v.set(point{center + 1, center - 1}, colorWhite)
		v.set(point{center - 2, center}, colorWhite)
	case *enhancedRobotPlayer:
		v.set(point{center, center}, colorBlack)
		v.set(point{center + 1, center}, colorBlack)
		v.set(point{center - 1, center + 1}, colorBlack)
		v.set(point{center + 2, center - 1}, colorBlack)
		v.set(point{center, center + 1}, colorWhite)
		v.set(point{center + 1, center + 1}, colorWhite)
		v.set(point{center - 1, center}, colorWhite)
		v.set(point{center + 1, center - 1}, colorWhite)
		v.set(point{center - 2, center}, colorWhite)
	}
}

func setupEndGamePosition(ai player) {
	// Crowded end-game position
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if (i+j)%7 == 0 && i < 10 && j < 10 {
				switch v := ai.(type) {
				case *robotPlayer:
					if (i+j)%2 == 0 {
						v.set(point{i, j}, colorBlack)
					} else {
						v.set(point{i, j}, colorWhite)
					}
				case *leanEnhancedRobotPlayer:
					if (i+j)%2 == 0 {
						v.set(point{i, j}, colorBlack)
					} else {
						v.set(point{i, j}, colorWhite)
					}
				case *enhancedRobotPlayer:
					if (i+j)%2 == 0 {
						v.set(point{i, j}, colorBlack)
					} else {
						v.set(point{i, j}, colorWhite)
					}
				}
			}
		}
	}
}
