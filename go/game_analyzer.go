package main

import (
	"fmt"
	"strings"
)

// Game analysis tool to understand why the AI lost
func main() {
	fmt.Println("五子棋游戏分析工具")
	fmt.Println("==================")
	
	// Parse the game sequence from the user's comment
	moves := []string{
		"黑(7,7)", "白(6,8)", "黑(7,6)", "白(7,8)", "黑(8,8)", "白(6,6)", 
		"黑(6,7)", "白(8,7)", "黑(5,8)", "白(8,5)", "黑(6,9)", "白(7,10)",
		"黑(3,10)", "白(4,9)", "黑(7,9)", "白(8,9)", "黑(6,10)", "白(5,11)",
		"黑(10,6)", "白(9,7)", "黑(6,11)", "白(6,12)", "黑(4,8)", "白(5,9)",
		"黑(3,9)", "白(3,8)", "黑(4,10)", "白(2,10)", "黑(2,9)", "白(10,7)",
		"黑(5,12)", "白(4,11)", "黑(9,8)", "白(9,6)", "黑(11,8)", "白(10,5)",
		"黑(11,4)", "白(9,5)", "黑(10,8)", "白(12,8)", "黑(7,5)", "白(9,4)",
		"黑(9,3)", "白(12,7)", "黑(11,7)", "白(11,6)", "黑(8,3)", "白(13,8)",
	}
	
	analyzeGame(moves)
}

func analyzeGame(moves []string) {
	// Create board
	board := make([][]playerColor, 15)
	for i := 0; i < 15; i++ {
		board[i] = make([]playerColor, 15)
	}
	
	robot := newOptimizedRobotPlayer(colorBlack).(*robotPlayer)
	
	fmt.Printf("分析 %d 步棋局...\n\n", len(moves))
	
	criticalMoves := []int{}
	
	for i, moveStr := range moves {
		// Parse move
		var x, y int
		var color string
		if strings.Contains(moveStr, "黑") {
			fmt.Sscanf(moveStr, "黑(%d,%d)", &x, &y)
			color = "黑"
		} else {
			fmt.Sscanf(moveStr, "白(%d,%d)", &x, &y)
			color = "白"
		}
		
		// Convert to internal coordinates (1-based to 0-based)
		x--
		y--
		
		playerColorVal := colorBlack
		if color == "白" {
			playerColorVal = colorWhite
		}
		
		// Update board
		board[y][x] = playerColorVal
		robot.set(point{x, y}, playerColorVal)
		
		// Check for winning move or missed opportunities
		if i%2 == 0 { // AI's move (black)
			// Analyze AI's choice
			aiMove := point{x, y}
			
			// What would the AI choose now?
			if i > 10 { // After some opening moves
				bestMove := robot.max(4, 100000000)
				if bestMove != nil && bestMove.p != aiMove {
					criticalMoves = append(criticalMoves, i+1)
					fmt.Printf("第%d步: AI选择了 (%d,%d)，但最佳可能是 (%d,%d)\n", 
						i+1, x+1, y+1, bestMove.p.x+1, bestMove.p.y+1)
				}
			}
		}
		
		// Check if this move creates immediate threats
		if checkForThreats(board, point{x, y}, playerColorVal) {
			fmt.Printf("第%d步 %s(%d,%d): 创建了威胁\n", i+1, color, x+1, y+1)
		}
		
		// Check if game ends
		if checkGameEnd(board, point{x, y}) {
			fmt.Printf("第%d步 %s(%d,%d): 游戏结束！%s获胜\n", i+1, color, x+1, y+1, color)
			break
		}
	}
	
	fmt.Printf("\n关键时刻: 第%v步可能是关键决策点\n", criticalMoves)
	
	// Provide strategic analysis
	fmt.Println("\n战略分析:")
	fmt.Println("1. AI作为先手方应该有理论优势")
	fmt.Println("2. 可能的问题：")
	fmt.Println("   - 搜索深度不够（当前4层可能不足）")
	fmt.Println("   - 评估函数参数需要调整") 
	fmt.Println("   - 缺乏对特定威胁模式的识别")
	
	// Suggest improvements
	fmt.Println("\n改进建议:")
	fmt.Println("1. 增加关键位置的搜索深度")
	fmt.Println("2. 改进威胁检测和防御")
	fmt.Println("3. 优化开局和中局的评估参数")
}

func checkForThreats(board [][]playerColor, lastMove point, color playerColor) bool {
	// Check if the move creates a threat (3 in a row with open ends)
	for _, dir := range fourDirections {
		count := 1 // Count the piece just placed
		
		// Count in one direction
		for i := 1; i < 5; i++ {
			p := lastMove.move(dir, i)
			if p.checkRange() && board[p.y][p.x] == color {
				count++
			} else {
				break
			}
		}
		
		// Count in opposite direction
		for i := 1; i < 5; i++ {
			p := lastMove.move(dir, -i)
			if p.checkRange() && board[p.y][p.x] == color {
				count++
			} else {
				break
			}
		}
		
		if count >= 3 {
			return true
		}
	}
	return false
}

func checkGameEnd(board [][]playerColor, lastMove point) bool {
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