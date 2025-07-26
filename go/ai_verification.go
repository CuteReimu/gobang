package main

import (
	"fmt"
)

// Quick AI strength verification
func main() {
	fmt.Println("AI改进验证测试")
	fmt.Println("================")
	
	// Test the same critical position from the user's game
	testCriticalGamePosition()
}

func testCriticalGamePosition() {
	fmt.Println("测试关键游戏位置 (模拟用户对局第13步)...")
	
	// Create AIs
	optimizedAI := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	balancedAI := newBalancedRobotPlayer(colorBlack).(*robotPlayer)
	
	// Setup the position from move 13 in user's game
	setupUserGamePosition13(optimizedAI)
	setupUserGamePosition13(balancedAI)
	
	// Get AI decisions
	fmt.Println("\nAI决策对比:")
	
	optimizedResult := optimizedAI.max(optimizedAI.maxLevelCount, 100000000)
	balancedResult := balancedAI.max(balancedAI.maxLevelCount, 100000000)
	
	if optimizedResult != nil {
		fmt.Printf("优化AI选择: (%d,%d) 评估值: %d\n", 
			optimizedResult.p.x+1, optimizedResult.p.y+1, optimizedResult.value)
	}
	
	if balancedResult != nil {
		fmt.Printf("平衡AI选择: (%d,%d) 评估值: %d\n", 
			balancedResult.p.x+1, balancedResult.p.y+1, balancedResult.value)
	}
	
	// The user said AI chose (3,10) but analysis suggested (4,9) was better
	expectedBetter := point{3, 8} // (4,9) in 1-based coordinates
	
	fmt.Println("\n分析:")
	if balancedResult != nil {
		if balancedResult.p == expectedBetter {
			fmt.Println("✓ 平衡AI选择了分析建议的更好走法")
		} else {
			fmt.Printf("? 平衡AI选择了不同走法: (%d,%d)\n", 
				balancedResult.p.x+1, balancedResult.p.y+1)
		}
		
		if balancedResult.value > optimizedResult.value {
			fmt.Println("✓ 平衡AI显示更好的位置评估")
		}
	}
	
	// Test another critical decision point
	fmt.Println("\n测试第15步决策...")
	testMove15Decision()
}

func setupUserGamePosition13(robot *robotPlayer) {
	robot.initBoardStatus()
	
	// Replay moves 1-12 from user's game
	moves := []struct{p point; color playerColor}{
		{point{6, 6}, colorBlack},   // 黑(7,7)
		{point{5, 7}, colorWhite},   // 白(6,8)
		{point{6, 5}, colorBlack},   // 黑(7,6)
		{point{6, 7}, colorWhite},   // 白(7,8)
		{point{7, 7}, colorBlack},   // 黑(8,8)
		{point{5, 5}, colorWhite},   // 白(6,6)
		{point{5, 6}, colorBlack},   // 黑(6,7)
		{point{7, 6}, colorWhite},   // 白(8,7)
		{point{4, 7}, colorBlack},   // 黑(5,8)
		{point{7, 4}, colorWhite},   // 白(8,5)
		{point{5, 8}, colorBlack},   // 黑(6,9)
		{point{6, 9}, colorWhite},   // 白(7,10)
	}
	
	for _, move := range moves {
		robot.set(move.p, move.color)
	}
	
	// Now AI needs to decide on move 13 (was 黑(3,10) = (2,9))
}

func testMove15Decision() {
	balancedAI := newBalancedRobotPlayer(colorBlack).(*robotPlayer)
	
	// Setup position for move 15
	balancedAI.initBoardStatus()
	
	// Add more moves to reach position 15
	moves := []struct{p point; color playerColor}{
		{point{6, 6}, colorBlack},   // 黑(7,7)
		{point{5, 7}, colorWhite},   // 白(6,8) 
		{point{6, 5}, colorBlack},   // 黑(7,6)
		{point{6, 7}, colorWhite},   // 白(7,8)
		{point{7, 7}, colorBlack},   // 黑(8,8)
		{point{5, 5}, colorWhite},   // 白(6,6)
		{point{5, 6}, colorBlack},   // 黑(6,7)
		{point{7, 6}, colorWhite},   // 白(8,7)
		{point{4, 7}, colorBlack},   // 黑(5,8)
		{point{7, 4}, colorWhite},   // 白(8,5)
		{point{5, 8}, colorBlack},   // 黑(6,9)
		{point{6, 9}, colorWhite},   // 白(7,10)
		{point{2, 9}, colorBlack},   // 黑(3,10)
		{point{3, 8}, colorWhite},   // 白(4,9)
	}
	
	for _, move := range moves {
		balancedAI.set(move.p, move.color)
	}
	
	result := balancedAI.max(balancedAI.maxLevelCount, 100000000)
	if result != nil {
		fmt.Printf("第15步平衡AI选择: (%d,%d) 评估值: %d\n", 
			result.p.x+1, result.p.y+1, result.value)
		
		// User's AI chose (7,9), analysis suggested (6,10) was better
		userChoice := point{6, 8}     // (7,9)
		betterChoice := point{5, 9}   // (6,10)
		
		if result.p == betterChoice {
			fmt.Println("✓ 平衡AI选择了分析建议的更好走法")
		} else if result.p == userChoice {
			fmt.Println("? 平衡AI选择了与原AI相同的走法")
		} else {
			fmt.Println("? 平衡AI选择了完全不同的走法")
		}
	}
}