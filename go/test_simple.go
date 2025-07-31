package main

import "fmt"

func testOptimizedAI() {
	fmt.Println("Creating optimized AI...")
	ai := newOptimizedRobotPlayer(colorBlack)
	
	fmt.Println("Testing first move...")
	move1, err := ai.play()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("First move: (%d,%d)\n", move1.x, move1.y)
	
	// Check if it's a valid position
	if move1.x < 0 || move1.x >= maxLen || move1.y < 0 || move1.y >= maxLen {
		fmt.Printf("ERROR: Invalid move coordinates!\n")
		return
	}
	
	// Check if position is actually empty before the move
	ai2 := newOptimizedRobotPlayer(colorBlack)
	if ai2.get(move1) != colorEmpty {
		fmt.Printf("ERROR: AI tried to place on non-empty position!\n")
		return
	}
	
	fmt.Println("Test passed!")
}

func main() {
	testOptimizedAI()
}
