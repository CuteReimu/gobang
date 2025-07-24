package main

import (
	"fmt"
	"time"
)

func simpleBenchmark() {
	fmt.Println("运行简单AI基准测试...")
	fmt.Println("==============================")
	
	// Test original parameters
	fmt.Println("\n测试原始AI...")
	original := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            colorBlack,
		maxLevelCount:     6,
		maxCountEachLevel: 16,
		maxCheckmateCount: 12,
		evalParams:        getDefaultEvaluationParams(),
	}
	original.initBoardStatus()
	
	// Setup test position
	setupTestPosition(original)
	
	start := time.Now()
	result1 := original.max(4, 100000000) // Use smaller depth for testing
	duration1 := time.Since(start)
	
	fmt.Printf("原始AI (深度 4): %.3f 秒", duration1.Seconds())
	if result1 != nil {
		fmt.Printf(" - 最佳走法: %v (价值: %d)\n", result1.p, result1.value)
	} else {
		fmt.Println(" - 未找到走法")
	}
	
	// Test optimized parameters
	fmt.Println("\n测试优化AI...")
	optimized := &robotPlayer{
		boardCache:        make(boardCache),
		pColor:            colorBlack,
		maxLevelCount:     5,
		maxCountEachLevel: 12,
		maxCheckmateCount: 10,
		evalParams:        getOptimizedEvaluationParams(),
	}
	optimized.initBoardStatus()
	
	// Setup same test position
	setupTestPosition(optimized)
	
	start = time.Now()
	result2 := optimized.max(4, 100000000) // Use same depth for fair comparison
	duration2 := time.Since(start)
	
	fmt.Printf("优化AI (深度 4): %.3f 秒", duration2.Seconds())
	if result2 != nil {
		fmt.Printf(" - 最佳走法: %v (价值: %d)\n", result2.p, result2.value)
	} else {
		fmt.Println(" - 未找到走法")
	}
	
	// Calculate improvement
	if duration1.Seconds() > 0 {
		improvement := ((duration1.Seconds() - duration2.Seconds()) / duration1.Seconds()) * 100
		fmt.Printf("\n性能提升: %.1f%%\n", improvement)
		
		if improvement > 0 {
			fmt.Printf("优化AI快了 %.1f倍\n", duration1.Seconds()/duration2.Seconds())
		} else {
			fmt.Println("优化AI没有更快 (可能需要更多调优)")
		}
	}
}

func setupTestPosition(robot *robotPlayer) {
	center := maxLen / 2
	
	// Set up a typical mid-game scenario for testing
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

func main() {
	simpleBenchmark()
}