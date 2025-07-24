package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	
	fmt.Println("Gobang AI Benchmark Tool")
	fmt.Println("========================")
	
	benchmarkAI()
	runSelfPlayTest()
}