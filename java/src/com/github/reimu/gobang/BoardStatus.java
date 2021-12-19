package com.github.reimu.gobang;

import java.util.Random;

public class BoardStatus {
	private static final long[][] BLACK_HASH = new long[Constant.MAX_LEN][Constant.MAX_LEN];
	private static final long[][] WHITE_HASH = new long[Constant.MAX_LEN][Constant.MAX_LEN];
	
	private final byte[][] board = new byte[Constant.MAX_LEN][Constant.MAX_LEN];
	private long hash;
	private int cnt;
	
	/**
	 * construct an empty BoardStatus
	 */
	public BoardStatus() {
		hash = 0;
		cnt = 0;
	}
	
	/**
	 * if (i, j) is 0, then set
	 * @param color color
	 */
	public void setIfEmpty(Point p, int color) {
		if (board[p.y][p.x] != 0) return;
		if (color == 0)
			return;
		else if (color == 1)
			hash ^= BLACK_HASH[p.y][p.x];
		else if (color == 2)
			hash ^= WHITE_HASH[p.y][p.x];
		else
			throw new IllegalArgumentException(p.toString() + color);
		board[p.y][p.x] = (byte) color;
		cnt++;
	}
	
	public void set(Point p, int color) {
		if (board[p.y][p.x] == color) return;
		if (color == 1) {
			hash ^= BLACK_HASH[p.y][p.x];
			cnt++;
		} else if (color == 2) {
			hash ^= WHITE_HASH[p.y][p.x];
			cnt++;
		} else if (color != 0)
			throw new IllegalArgumentException(p.toString() + color);
		if (board[p.y][p.x] == 1) {
			hash ^= BLACK_HASH[p.y][p.x];
			cnt--;
		} else if (board[p.y][p.x] == 2) {
			hash ^= WHITE_HASH[p.y][p.x];
			cnt--;
		}
		board[p.y][p.x] = (byte) color;
	}

	public byte get(Point p) {
		if (null == p) return -1;
		return board[p.y][p.x];
	}

	public byte get(int x, int y) {
		if (!Point.checkRange(x, y)) return -1;
		return board[y][x];
	}

	public int count() {
		return cnt;
	}
	
	public long getZobrist() {
		return hash;
	}
	
	/**
	 * check if p is not too far from the exist points
	 * @param p point
	 * @return boolean
	 */
	public boolean isNeighbor(Point p) {
		if (null == p) return false;
		for (int i = -2; i <= 2; i++)
			for (int j = -2; j <= 2; j++)
				if (get(p.x + j, p.y + i) > 0)
					return true;
		return false;
	}
	
	static {
		Random rand = new Random(1551980916123L);
		for (int i = 0; i < Constant.MAX_LEN; i++) {
			for (int j = 0; j < Constant.MAX_LEN; j++) {
				BLACK_HASH[i][j] = rand.nextLong();
				WHITE_HASH[i][j] = rand.nextLong();
			}
		}
	}
}
