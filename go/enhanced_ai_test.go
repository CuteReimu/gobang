package main

import (
	"fmt"
	"testing"
	"time"
)

// Test the enhanced AI functionality
func TestEnhancedAI(t *testing.T) {
	// Create enhanced AI
	enhanced := newEnhancedRobotPlayer(colorBlack)
	
	// Test basic initialization
	if enhanced.maxLevelCount != 6 {
		t.Errorf("Expected maxLevelCount to be 6, got %d", enhanced.maxLevelCount)
	}
	
	if enhanced.aspirationWindow != 50000 {
		t.Errorf("Expected aspirationWindow to be 50000, got %d", enhanced.aspirationWindow)
	}
	
	// Test first move (should be center)
	start := time.Now()
	p, err := enhanced.play()
	duration := time.Since(start)
	
	if err != nil {
		t.Errorf("Enhanced AI failed to play first move: %v", err)
	}
	
	expectedCenter := point{maxLen / 2, maxLen / 2}
	if p != expectedCenter {
		t.Errorf("Expected first move to be center %v, got %v", expectedCenter, p)
	}
	
	fmt.Printf("Enhanced AI first move time: %v\n", duration)
	
	// Test second move performance
	start = time.Now()
	p2, err := enhanced.play()
	duration = time.Since(start)
	
	if err != nil {
		t.Errorf("Enhanced AI failed to play second move: %v", err)
	}
	
	fmt.Printf("Enhanced AI second move time: %v\n", duration)
	fmt.Printf("Enhanced AI node count: %d\n", enhanced.nodeCount)
	
	// Verify the move is valid (should be near center)
	if !p2.checkRange() {
		t.Errorf("Enhanced AI made invalid move: %v", p2)
	}
}

// Benchmark enhanced AI vs original AI
func BenchmarkEnhancedAIvsOriginal(b *testing.B) {
	fmt.Println("\n=== Enhanced AI vs Original AI Performance Comparison ===")
	
	// Test scenarios
	scenarios := []struct {
		name string
		moveCount int
	}{
		{"Early Game", 4},
		{"Mid Game", 12},
		{"Late Game", 20},
	}
	
	for _, scenario := range scenarios {
		fmt.Printf("\n--- %s (%d moves) ---\n", scenario.name, scenario.moveCount)
		
		// Test Original AI
		original := newRobotPlayer(colorBlack)
		setupTestBoard(original, scenario.moveCount)
		
		start := time.Now()
		_, err := original.play()
		originalTime := time.Since(start)
		
		if err != nil {
			fmt.Printf("Original AI error: %v\n", err)
			continue
		}
		
		// Test Enhanced AI
		enhanced := newEnhancedRobotPlayer(colorBlack)
		setupTestBoard(enhanced, scenario.moveCount)
		
		start = time.Now()
		_, err = enhanced.play()
		enhancedTime := time.Since(start)
		
		if err != nil {
			fmt.Printf("Enhanced AI error: %v\n", err)
			continue
		}
		
		speedup := float64(originalTime) / float64(enhancedTime)
		fmt.Printf("Original AI (6层):  %v\n", originalTime)
		fmt.Printf("Enhanced AI (6层):  %v\n", enhancedTime)
		fmt.Printf("Enhanced AI nodes: %d\n", enhanced.nodeCount)
		if speedup >= 1.0 {
			fmt.Printf("性能提升: %.1fx 更快\n", speedup)
		} else {
			fmt.Printf("性能: %.1fx 原始AI速度\n", 1.0/speedup)
		}
	}
}

// setupTestBoard creates a test board scenario
func setupTestBoard(ai player, moveCount int) {
	// Get the underlying robot player interface
	var rp *robotPlayer
	switch v := ai.(type) {
	case *robotPlayer:
		rp = v
	case *enhancedRobotPlayer:
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

// Run the benchmark test
func runEnhancedAIBenchmark() {
	fmt.Println("运行增强AI性能基准测试...")
	testing.Benchmark(BenchmarkEnhancedAIvsOriginal)
}