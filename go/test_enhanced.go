package main

import (
	"fmt"
	"time"
)

// Simple standalone test for enhanced AI
func main() {
	fmt.Println("=== 增强AI（6层深度）性能测试 ===")
	
	// Test enhanced AI
	fmt.Println("\n创建增强AI（6层深度）...")
	enhanced := newEnhancedRobotPlayer(colorBlack).(*enhancedRobotPlayer)
	
	// Test first move
	fmt.Println("测试第一步（应该是中心位置）...")
	start := time.Now()
	p1, err := enhanced.play()
	firstMoveTime := time.Since(start)
	
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	fmt.Printf("第一步: %v, 用时: %v\n", p1, firstMoveTime)
	
	// Test second move
	fmt.Println("测试第二步...")
	start = time.Now()
	p2, err := enhanced.play()
	secondMoveTime := time.Since(start)
	
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	fmt.Printf("第二步: %v, 用时: %v\n", p2, secondMoveTime)
	fmt.Printf("搜索节点数: %d\n", enhanced.nodeCount)
	
	// Compare with original AI
	fmt.Println("\n=== 对比原始AI（6层深度）===")
	
	// Test original AI with same board state
	original := newRobotPlayer(colorBlack).(*robotPlayer)
	original.set(p1, colorBlack)
	
	start = time.Now()
	p3, err := original.play()
	originalTime := time.Since(start)
	
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	fmt.Printf("原始AI第二步: %v, 用时: %v\n", p3, originalTime)
	
	// Calculate speedup
	if originalTime > secondMoveTime {
		speedup := float64(originalTime) / float64(secondMoveTime)
		fmt.Printf("\n性能提升: %.1fx 更快\n", speedup)
	} else {
		ratio := float64(secondMoveTime) / float64(originalTime)
		fmt.Printf("\n性能: %.1fx 原始AI速度\n", ratio)
	}
	
	// Test mid-game scenario
	fmt.Println("\n=== 中盘测试 ===")
	testMidGame()
}

func testMidGame() {
	// Create a mid-game scenario
	enhanced := newEnhancedRobotPlayer(colorBlack).(*enhancedRobotPlayer)
	original := newRobotPlayer(colorBlack).(*robotPlayer)
	
	// Set up identical board positions for both AIs
	positions := []point{
		{7, 7}, {6, 8}, {7, 6}, {7, 8}, {8, 8}, {6, 6},
		{6, 7}, {8, 7}, {5, 8}, {8, 5}, {6, 9}, {7, 10},
	}
	
	for i, pos := range positions {
		color := colorBlack
		if i%2 == 1 {
			color = colorWhite
		}
		enhanced.set(pos, color)
		original.set(pos, color)
	}
	
	// Test enhanced AI
	start := time.Now()
	enhancedMove, err := enhanced.play()
	enhancedTime := time.Since(start)
	
	if err != nil {
		fmt.Printf("增强AI错误: %v\n", err)
		return
	}
	
	// Test original AI
	start = time.Now()
	originalMove, err := original.play()
	originalTime := time.Since(start)
	
	if err != nil {
		fmt.Printf("原始AI错误: %v\n", err)
		return
	}
	
	fmt.Printf("增强AI (6层): 着法 %v, 用时 %v, 节点 %d\n", enhancedMove, enhancedTime, enhanced.nodeCount)
	fmt.Printf("原始AI (6层): 着法 %v, 用时 %v\n", originalMove, originalTime)
	
	if originalTime > enhancedTime {
		speedup := float64(originalTime) / float64(enhancedTime)
		fmt.Printf("中盘性能提升: %.1fx 更快\n", speedup)
	} else {
		ratio := float64(enhancedTime) / float64(originalTime)
		fmt.Printf("中盘性能: %.1fx 原始AI速度\n", ratio)
	}
}