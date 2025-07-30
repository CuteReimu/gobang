package main

import (
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	optimized := flag.Bool("optimized", false, "使用优化AI参数以获得更好性能")
	balanced := flag.Bool("balanced", false, "使用平衡AI参数以获得更强棋力和合理速度")
	enhanced := flag.Bool("enhanced", false, "使用增强AI算法以获得6层深度的强棋力和优化性能")
	benchmark := flag.Bool("benchmark", false, "运行AI性能基准测试")
	flag.Parse()

	if *benchmark {
		benchmarkAI()
		runSelfPlayTest()
		runEnhancedAIBenchmark()
		testPerformance()
		testLogicErrors()
		return
	}

	runGameWithGUI(*optimized, *balanced, *enhanced)
}

func runGameWithGUI(optimized, balanced, enhanced bool) {
	hp := newHumanPlayer(colorWhite)
	//hp := newHumanWatcher()
	go func() {
		board := make([][]playerColor, maxLen)
		for i := 0; i < maxLen; i++ {
			board[i] = make([]playerColor, maxLen)
		}

		var robot player
		if enhanced {
			robot = newEnhancedRobotPlayer(colorBlack)
			fmt.Println("使用增强AI算法以获得6层深度的强棋力和优化性能")
		} else if balanced {
			robot = newBalancedRobotPlayer(colorBlack)
			fmt.Println("使用平衡AI参数以获得更强棋力和合理速度")
		} else if optimized {
			robot = newOptimizedRobotPlayer(colorBlack)
			fmt.Println("使用增强优化AI以获得强棋力和良好性能")
			fmt.Println("- 搜索深度：6层（确保战术强度）")
			fmt.Println("- 候选着法：18个（全面战术分析）")
			fmt.Println("- 自适应深度：复杂局面自动加深至8层")
			fmt.Println("- 强化评估：更准确的威胁识别和防御")
			fmt.Println("- 平衡策略：保持速度的同时不牺牲棋力")
		} else {
			robot = newRobotPlayer(colorBlack)
			fmt.Println("使用默认AI参数")
		}

		players := []player{robot, hp} // 机器人先
		//players := []player{hp, newRobotPlayer(colorWhite)} // 玩家先
		//players := []player{newRobotPlayer(colorBlack), newRobotPlayer(colorWhite)}
		var watchers []*humanWatcher
		//watchers = append(watchers, hp)
		count := 0
		whoseTurn := 0
		checkForWin := func(p point) bool {
			whose := board[p.y][p.x]
			for _, dir := range fourDirections {
				count := 0
				for i := -4; i <= 4; i++ {
					pi := p.move(dir, i)
					if !pi.checkRange() || board[pi.y][pi.x] != whose {
						count = 0
					} else {
						count++
						if count == 5 {
							return true
						}
					}
				}
			}
			return false
		}
		for {
			p, err := players[whoseTurn].play()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if board[p.y][p.x] != 0 {
				log.Printf("illegal argument: %s%s\n", p, board[p.y][p.x])
				continue
			}
			board[p.y][p.x] = players[whoseTurn].color()
			fmt.Printf("%s%s\n", board[p.y][p.x], p)
			whoseTurn = 1 - whoseTurn
			if err := players[whoseTurn].display(p); err != nil {
				log.Println(err.Error())
			}
			for _, watcher := range watchers {
				if err := watcher.display(board[p.y][p.x], p); err != nil {
					log.Println(err.Error())
				}
			}
			count++
			if count == maxLen*maxLen || checkForWin(p) {
				break
			}
		}
		select {}
	}()
	ebiten.SetWindowSize(35*(maxLen+1), 35*(maxLen+1))
	ebiten.SetWindowTitle("gobang")
	if err := ebiten.RunGame(hp); err != nil {
		panic(err)
	}
}
