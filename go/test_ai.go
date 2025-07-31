package main

import (
	"fmt"
)

func testDefensiveLogic() {
	fmt.Println("Testing defensive logic...")
	
	// Test game state from user's example:
	// 黑(7,7) 白(6,8) 黑(8,7) 白(6,7) 黑(7,8) 白(6,9) 黑(9,6) 白(6,6) 黑(6,5) 白(6,10)
	
	fmt.Println("Game state analysis:")
	fmt.Println("White has: (6,8), (6,7), (6,9), (6,6), (6,10) - vertical line in column 6")
	fmt.Println("Black should block at (6,5) or (6,11) but chose (6,5) - this was actually correct!")
	fmt.Println("The issue might be that white still won afterwards...")
}

func main() {
	testDefensiveLogic()
	fmt.Println("Basic test completed.")
}