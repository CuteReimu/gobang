package com.github.reimu.gobang;

public final class Point {
	public final int x;
	public final int y;
	public Point(int x, int y) {
		this.x = x;
		this.y = y;
	}
	Point move(Direction dir, int len) {
		if (len == 0) return new Point(x, y);
		int x = this.x + dir.x * len;
		int y = this.y + dir.y * len;
		if (!checkRange(x, y)) return null;
		return new Point(x, y);
	}
	public boolean nearMidThan(Point p) {
		return Math.max(Math.abs(x - Constant.MAX_LEN / 2), Math.abs(y - Constant.MAX_LEN / 2))
			< Math.max(Math.abs(p.x - Constant.MAX_LEN / 2), Math.abs(p.y - Constant.MAX_LEN / 2));
	}
	@Override
	public boolean equals(Object obj) {
		if (!(obj instanceof Point))
			return false;
		return ((Point) obj).x == x && ((Point) obj).y == y;
	}
	@Override
	public int hashCode() {
		return y * Constant.MAX_LEN + x;
	}
	@Override
	public String toString() {
		return toString(x, y);
	}
	private static String toString(int x, int y) {
		return "(" + x + ',' + y + ')';
	}
	public static boolean checkRange(int x, int y) {
		return x < Constant.MAX_LEN && x >= 0 && y < Constant.MAX_LEN && y >= 0;
	}
}
