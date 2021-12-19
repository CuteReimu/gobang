package com.github.reimu.gobang;

public enum Direction {
	LEFT(-1, 0), UPLEFT(-1, -1), UP(0, -1), UPRIGHT(1, -1), RIGHT(1, 0), DOWNRIGHT(1, 1), DOWN(0, 1), DOWNLEFT(-1, 1);
	public final int x;
	public final int y;
	private Direction(int x, int y) {
		this.x = x;
		this.y = y;
	}
	private static final Direction[] _4dir = new Direction[]{LEFT, UPLEFT, UP, UPRIGHT};
	public static Direction[] get4Directions() {
		return _4dir;
	}
}
