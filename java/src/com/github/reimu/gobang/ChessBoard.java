package com.github.reimu.gobang;

import java.util.ArrayList;
import java.util.List;

public class ChessBoard {
	public static void main(String[] args) {
		ChessBoard b = new ChessBoard();
//		b.setPlayer(new HumanPlayer(1), new RobotPlayer(2)); //人对机
		b.setPlayer(new RobotPlayer(1), new HumanPlayer(2)); //机对人
//		b.setPlayer(new RobotPlayer(1), new RobotPlayer(2)); //机对机
//		b.addWatcher(new HumanWatcher());					 //机对机的时候加一个旁观的人
		b.clear();
		while (b.getWinner() == null) {
//			try {
//				Thread.sleep(500);
//			} catch (InterruptedException e) {
//				e.printStackTrace();
//			}
			b.play();
		}
	}
	
	private final byte[][] board = new byte[Constant.MAX_LEN][Constant.MAX_LEN];
	private final Player[] players = new Player[2];
	private int whoseTurn = 0;
	private int count = 0;
	private boolean isEnd = false;
	private final List<HumanWatcher> watchers = new ArrayList<>();
	
	public void setPlayer(Player a, Player b) {
		players[0] = a;
		players[1] = b;
	}
	
	public void addWatcher(HumanWatcher watcher) {
		watchers.add(watcher);
	}
	
	public void clear() {
		for (int i = 0; i < Constant.MAX_LEN; i++)
			for (int j = 0; j < Constant.MAX_LEN; j++)
				board[i][j] = 0;
		isEnd = false;
		whoseTurn = 0;
	}
	
	public void play() {
		if (isEnd)
			return;
		Point p = players[whoseTurn].play();
		if (board[p.y][p.x] != 0)
			throw new IllegalArgumentException(p.toString() + board[p.y][p.x]);
		board[p.y][p.x] = (byte) (whoseTurn + 1);
		System.out.println((whoseTurn == 0 ? "黑" : "白") + p.toString());
		if (++count == Constant.MAX_LEN * Constant.MAX_LEN || checkForWin(p))
			isEnd = true;
		whoseTurn = 1 - whoseTurn;
		players[whoseTurn].display(p);
		for (HumanWatcher watcher : watchers)
			watcher.display(p, 2 - whoseTurn);
	}
	
	public Player getWinner() {
		if (!isEnd) return null;
		return players[1 - whoseTurn];
	}
	
	private boolean checkForWin(Point p) {
		int whose = board[p.y][p.x];
		int count = 0;
		for (int i = -4; i <= 4; i++) {
			if (!Point.checkRange(p.x + i, p.y + i) || board[p.y + i][p.x + i] != whose)
				count = 0;
			else if (++count == 5)
				return true;
		}
		count = 0;
		for (int i = -4; i <= 4; i++) {
			if (!Point.checkRange(p.x + i, p.y) || board[p.y][p.x + i] != whose)
				count = 0;
			else if (++count == 5)
				return true;
		}
		count = 0;
		for (int i = -4; i <= 4; i++) {
			if (!Point.checkRange(p.x - i, p.y + i) || board[p.y + i][p.x - i] != whose)
				count = 0;
			else if (++count == 5)
				return true;
		}
		count = 0;
		for (int i = -4; i <= 4; i++) {
			if (!Point.checkRange(p.x, p.y + i) || board[p.y + i][p.x] != whose)
				count = 0;
			else if (++count == 5)
				return true;
		}
		return false;
	}
}
