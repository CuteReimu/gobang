package main

import (
	"fmt"
	"time"
)

func simpleBenchmark() {
	fmt.Println("Running Simple AI Benchmark...")
	fmt.Println("==============================")
	
	// Test original parameters
	fmt.Println("\nTesting Original AI...")
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
	
	fmt.Printf("Original AI (depth 4): %.3f seconds", duration1.Seconds())
	if result1 != nil {
		fmt.Printf(" - Best move: %v (value: %d)\n", result1.p, result1.value)
	} else {
		fmt.Println(" - No move found")
	}
	
	// Test optimized parameters
	fmt.Println("\nTesting Optimized AI...")
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
	
	fmt.Printf("Optimized AI (depth 4): %.3f seconds", duration2.Seconds())
	if result2 != nil {
		fmt.Printf(" - Best move: %v (value: %d)\n", result2.p, result2.value)
	} else {
		fmt.Println(" - No move found")
	}
	
	// Calculate improvement
	if duration1.Seconds() > 0 {
		improvement := ((duration1.Seconds() - duration2.Seconds()) / duration1.Seconds()) * 100
		fmt.Printf("\nPerformance improvement: %.1f%%\n", improvement)
		
		if improvement > 0 {
			fmt.Printf("Optimized AI is %.1fx faster\n", duration1.Seconds()/duration2.Seconds())
		} else {
			fmt.Println("Optimized AI is not faster (may need more tuning)")
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