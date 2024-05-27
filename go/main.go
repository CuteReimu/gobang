package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	hp := newHumanPlayer(colorWhite)
	//hp := newHumanWatcher()
	go func() {
		board := make([][]playerColor, maxLen)
		for i := 0; i < maxLen; i++ {
			board[i] = make([]playerColor, maxLen)
		}
		players := []player{newRobotPlayer(colorBlack), hp} // 机器人先
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
