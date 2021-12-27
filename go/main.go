package main

import (
	"fmt"
	"log"
)

func main() {
	board := make([][]playerColor, maxLen)
	for i := 0; i < maxLen; i++ {
		board[i] = make([]playerColor, maxLen)
	}
	players := []*AIPlayer{newAIPlayer(colorBlack)}
	players = append(players, players[0].copy(colorWhite))
	players[1].rand()
	//players := []player{newRobotPlayer(colorBlack), newHumanPlayer(colorWhite)} // 机器人先
	//players := []player{newHumanPlayer(colorBlack), newRobotPlayer(colorWhite)} // 玩家先
	var watchers []*humanWatcher
	//watchers = append(watchers, newHumanWatcher())
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
	gameCount, winCount := 0, 0
	for {
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
			//fmt.Printf("%s%s\n", board[p.y][p.x], p)
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
		//break
		for i := 0; i < maxLen; i++ {
			board[i] = make([]playerColor, maxLen)
		}
		gameCount++
		if 1-whoseTurn == 1 {
			winCount++
			fmt.Println(winCount, "/", gameCount)
		}
		count = 0
		players[1-whoseTurn].save()
		players = []*AIPlayer{players[1-whoseTurn].copy(colorBlack), players[1-whoseTurn].copy(colorWhite)}
		players[1].rand()
		whoseTurn = 0
	}
	//select {
	//
	//}
}
