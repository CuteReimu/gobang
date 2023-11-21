package com.github.reimu.gobang;

import java.util.PriorityQueue;

public class RobotPlayer implements Player {
    private static final int MIN_VALUE = -100000000;
    private static final int MAX_VALUE = 100000000;
    private static final int MAX_LEVEL_COUNT = 6;
    private static final int MAX_COUNT_EACH_LEVEL = 16;
    private static final int MAX_CHECKMATE_COUNT = 12;
    /**
     * 1-黑 2-白
     */
    private final int color;
    /**
     * 一共走了几步
     */
    private final BoardCache boardCache = new BoardCache();
    private final BoardStatus board = new BoardStatus();

    public RobotPlayer(int color) {
        this.color = color;
    }

    @Override
    public void display(Point p) {
        if (board.get(p) != 0)
            throw new IllegalArgumentException(p.toString() + board.get(p));
        board.set(p, (byte) (3 - color));
    }

    @Override
    public Point play() {
        if (board.count() == 0) {
            Point p = new Point(Constant.MAX_LEN / 2, Constant.MAX_LEN / 2);
            board.set(p, color);
            return p;
        }
        Point p1 = findForm5(color);
        if (p1 != null) {
            board.set(p1, color);
            return p1;
        }
        p1 = stop4(color);
        if (p1 != null) {
            board.set(p1, color);
            return p1;
        }
        for (int i = 2; i <= MAX_CHECKMATE_COUNT; i += 2) {
            Point p = calculateKill(color, true, i);
            if (p != null) {
                return p;
            }
        }
        PointAndValue result = max(MAX_LEVEL_COUNT, 100000000);
        if (result == null) {
            throw new RuntimeException("algorithm error");
        }
        board.set(result.point, color);
        return result.point;
    }

    private Point calculateKill(int color, boolean aggresive, int step) {
        if (step == 0) return null;
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                Point p = new Point(j, i);
                if (board.get(p) == 0) {
                    board.set(p, color);
                    if (!exists4(3 - color) && (!aggresive || exists4(color)) && null == calculateKill(3 - color, !aggresive, step - 1)) {
                        board.set(p, 0);
                        return new Point(j, i);
                    }
                    board.set(p, 0);
                }
            }
        }
        return null;
    }

    private Point stop4(int color) {
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                Point p = new Point(j, i);
                if (board.get(p) == 0) {
                    for (Direction dir : Direction.get4Directions()) {
                        int leftCount = 0, rightCount = 0;
                        for (int k = -1; k >= -4; k--) {
                            Point p1 = p.move(dir, k);
                            if (board.get(p1) == 3 - color) {
                                leftCount++;
                            } else {
                                break;
                            }
                        }
                        for (int k = 1; k <= 4; k++) {
                            Point p1 = p.move(dir, k);
                            if (board.get(p1) == 3 - color) {
                                rightCount++;
                            } else {
                                break;
                            }
                        }
                        if (leftCount + rightCount >= 4) {
                            return p;
                        }
                    }
                }
            }
        }
        return null;
    }

    private boolean exists4(int color) {
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                if (board.get(j, i) == color || board.get(j, i) == 0) {
                    Point p = new Point(j, i);
                    for (Direction dir : Direction.get4Directions()) {
                        int count0 = 0, count1 = 0;
                        for (int k = 0; k <= 4; k++) {
                            int kColor = board.get(p.move(dir, k));
                            if (kColor == 0) count0++;
                            else if (kColor == color) count1++;
                        }
                        if (count0 == 1 && count1 == 4) {
                            return true;
                        }
                    }
                }
            }
        }
        return false;
    }

    private Point findForm5(int color) {
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                if (board.get(j, i) == 0) {
                    Point p = new Point(j, i);
                    for (Direction dir : Direction.get4Directions()) {
                        int leftCount = 0, rightCount = 0;
                        for (int k = -1; k >= -4; k--) {
                            if (board.get(p.move(dir, k)) == color) leftCount++;
                            else break;
                        }
                        for (int k = 1; k <= 4; k++) {
                            if (board.get(p.move(dir, k)) == color) rightCount++;
                            else break;
                        }
                        if (leftCount + rightCount >= 4)
                            return p;
                    }
                }
            }
        }
        return null;
    }

    private PointAndValue max(int step, int foundminVal) {
        PointAndValue cacheResult = boardCache.get(board.getZobrist(), step);
        if (null != cacheResult) return cacheResult;
        PriorityQueue<PointAndValue> queue = new PriorityQueue<>((obj1, obj2) -> Integer.compare(obj2.value, obj1.value));
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                Point p = new Point(j, i);
                if (board.get(p) == 0 && board.isNeighbor(p)) {
                    int evathis = evaluatePoint(p, color);
                    queue.add(new PointAndValue(p, evathis));
                }
            }
        }
        if (step == 1) {
            assert queue.peek() != null;
            Point p = queue.peek().point;
            board.setIfEmpty(p, color);
            int val = evaluateBoard(color) - evaluateBoard(3 - color);
            board.set(p, 0);
            PointAndValue result = new PointAndValue(p, val);
            boardCache.put(board.getZobrist(), result, step);
            return result;
        }
        Point max = null;
        int maxVal = MIN_VALUE;
        int i = 0;
        PointAndValue obj;
        while ((obj = queue.poll()) != null) {
            if (++i > MAX_COUNT_EACH_LEVEL) break;
            Point p = obj.point;
            board.set(p, color);
            int boardVal = evaluateBoard(color) - evaluateBoard(3 - color);
            if (boardVal > 800000) {
                board.set(p, 0);
                PointAndValue result = new PointAndValue(p, boardVal);
                boardCache.put(board.getZobrist(), result, step);
                return result;
            }
            PointAndValue result = min(step - 1, maxVal);//最大值最小值法
            assert result != null;
            int evathis = result.value;
            if (evathis >= foundminVal) {
                board.set(p, 0);
                return new PointAndValue(p, evathis);
            }
            if (max == null || evathis > maxVal || evathis == maxVal && p.nearMidThan(max)) {
                maxVal = evathis;
                max = p;
            }
            board.set(p, 0);
        }
        if (max == null) return null;
        PointAndValue result = new PointAndValue(max, maxVal);
        boardCache.put(board.getZobrist(), result, step);
        return result;
    }

    private PointAndValue min(int step, int foundmaxVal) {
        PointAndValue cacheResult = boardCache.get(board.getZobrist(), step);
        if (null != cacheResult) return cacheResult;
        PriorityQueue<PointAndValue> queue = new PriorityQueue<>((obj1, obj2) -> Integer.compare(obj2.value, obj1.value));
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                Point p = new Point(j, i);
                if (board.get(p) == 0 && board.isNeighbor(p)) {
                    int evathis = evaluatePoint(p, 3 - color);
                    queue.add(new PointAndValue(p, evathis));
                }
            }
        }
        if (step == 1) {
            assert queue.peek() != null;
            Point p = queue.peek().point;
            board.setIfEmpty(p, 3 - color);
            int val = evaluateBoard(color) - evaluateBoard(3 - color);
            board.set(p, 0);
            PointAndValue result = new PointAndValue(p, val);
            boardCache.put(board.getZobrist(), result, step);
            return result;
        }
        Point min = null;
        int minVal = MAX_VALUE;
        int i = 0;
        PointAndValue obj;
        while ((obj = queue.poll()) != null) {
            if (++i > MAX_COUNT_EACH_LEVEL) break;
            Point p = obj.point;
            board.set(p, 3 - color);
            int boardVal = evaluateBoard(color) - evaluateBoard(3 - color);
            if (boardVal < -800000) {
                board.set(p, 0);
                PointAndValue result = new PointAndValue(p, boardVal);
                boardCache.put(board.getZobrist(), result, step);
                return result;
            }
            PointAndValue result = max(step - 1, minVal);//最大值最小值法
            assert result != null;
            int evathis = result.value;
            if (evathis <= foundmaxVal) {
                board.set(p, 0);
                return new PointAndValue(p, evathis);
            }
            if (min == null || evathis < minVal || evathis == minVal && p.nearMidThan(min)) {
                minVal = evathis;
                min = p;
            }
            board.set(p, 0);
        }
        if (min == null) return null;
        PointAndValue result = new PointAndValue(min, minVal);
        boardCache.put(board.getZobrist(), result, step);
        return result;
    }

    private int evaluatePoint(Point p, int color) {
        return evaluatePoint(p, color, 1) + evaluatePoint(p, color, 2);
    }

    private int getLine(Point p, Direction dir, int j) {
        Point p2 = p.move(dir, j);
        if (p2 != null) {
            return board.get(p2);
        }
        return -1;
    }

    private int evaluatePoint(Point p, int me, int plyer) {
        int value = 0;
        int numoftwo = 0;
        for (Direction dir : Direction.values()) { // 8个方向
            // 活四 01111* *代表当前空位置 0代表其他空位置 下同
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer && getLine(p, dir, -5) == 0) {
                value += 300000;
                if (me != plyer)
                    value -= 500;
                continue;
            }
            // 死四A 21111*
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer && (getLine(p, dir, -5) == 3 - plyer || getLine(p, dir, -5) == -1)) {
                value += 2500009;
                if (me != plyer)
                    value -= 500;
                continue;
            }
            // 死四B 111*1
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, 1) == plyer) {
                value += 240000;
                if (me != plyer)
                    value -= 500;
                continue;
            }
            // 死四C 11*11
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer && getLine(p, dir, 2) == plyer) {
                value += 230000;
                if (me != plyer)
                    value -= 500;
                continue;
            }
            // 活三 近3位置 111*0
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer) {
                if (getLine(p, dir, 1) == 0) {
                    value += 1450;
                    if (getLine(p, dir, -4) == 0) {
                        value += 6000;
                        if (me != plyer)
                            value -= 300;
                    }
                }
                if ((getLine(p, dir, 1) == 3 - plyer || getLine(p, dir, 1) == -1) && getLine(p, dir, -4) == 0)
                    value += 500;
                if ((getLine(p, dir, -4) == 3 - plyer || getLine(p, dir, -4) == -1) && getLine(p, dir, 1) == 0)
                    value += 500;
                continue;
            }
            // 活三 远3位置 1110*
            if (getLine(p, dir, -1) == 0 && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == plyer && getLine(p, dir, -4) == plyer) {
                value += 350;
                continue;
            }
            // 死三 11*1
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, 1) == plyer) {
                value += 700;
                if (getLine(p, dir, -3) == 0 && getLine(p, dir, 2) == 0) {
                    value += 6700;
                    continue;
                }
                if ((getLine(p, dir, -3) == 3 - plyer || getLine(p, dir, -3) == -1) && (getLine(p, dir, 2) == 3 - plyer || getLine(p, dir, 2) == -1))
                    value -= 700;
                else
                    value += 800;
                continue;
            }
            // 活二的个数（因为会算2次，就2倍）
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == plyer && getLine(p, dir, -3) == 0 && getLine(p, dir, 1) == 0) {
                if (getLine(p, dir, 2) == 0 || getLine(p, dir, -4) == 0)
                    numoftwo += 2;
                else
                    value += 250;
            }
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, -2) == 0 && getLine(p, dir, 2) == plyer && getLine(p, dir, 1) == 0 && getLine(p, dir, 3) == 0)
                numoftwo += 2;
            if (getLine(p, dir, -1) == 0 && getLine(p, dir, 4) == 0 && getLine(p, dir, 3) == plyer && (getLine(p, dir, 2) == plyer && getLine(p, dir, 1) == 0 || getLine(p, dir, 1) == plyer && getLine(p, dir, 2) == 0))
                numoftwo += 2;
            if (getLine(p, dir, -1) == plyer && getLine(p, dir, 1) == plyer && getLine(p, dir, -2) == 0 && getLine(p, dir, 2) == 0) {
                if (getLine(p, dir, 3) == 0 || getLine(p, dir, -3) == 0)
                    numoftwo++;
                else
                    value += 125;
            }
            // 其余散棋
            int numOfplyer = 0;
            for (int k = -4; k <= 0; k++) { // ++++* +++*+ ++*++ +*+++ *++++
                int temp = 0;
                for (int l = 0; l <= 4; l++) {
                    if (getLine(p, dir, k + l) == plyer) {
                        temp += 5 - Math.abs(k + l);
                    } else if (getLine(p, dir, k + l) == 3 - plyer || getLine(p, dir, k + l) == -1) {
                        temp = 0;
                        break;
                    }
                }
                numOfplyer += temp;
            }
            value += numOfplyer * 5;
        }
        numoftwo /= 2;
        if (numoftwo >= 2) {
            value += 3000;
            if (me != plyer)
                value -= 100;
        } else if (numoftwo == 1) {
            value += 2725;
            if (me != plyer)
                value -= 10;
        }
        return value;
    }

    private int evaluateBoard(int color) {
        int values = 0;
        for (int i = 0; i < Constant.MAX_LEN; i++) {
            for (int j = 0; j < Constant.MAX_LEN; j++) {
                Point p = new Point(j, i);
                if (board.get(p) != color)
                    continue;
                for (Direction dir : Direction.values()) {
                    int[] colors = new int[9];
                    for (int k = 0; k < 9; k++) {
                        Point pk = p.move(dir, k - 4);
                        if (pk != null)
                            colors[k] = board.get(pk);
                        else
                            colors[k] = -1;
                    }
                    if (colors[5] == color && colors[6] == color && colors[7] == color && colors[8] == color) {
                        values += 1000000;
                        continue;
                    }
                    if (colors[5] == color && colors[6] == color && colors[7] == color && colors[3] == 0) {
                        if (colors[8] == 0) //?AAAA?
                            values += 300000 / 2;
                        else //AAAA?
                            values += 25000;
                        continue;
                    }
                    final boolean b = colors[3] != 0 && colors[3] != color && colors[7] == 0 && colors[8] == 0;
                    if (colors[5] == color && colors[6] == color) {
                        if (colors[7] == 0 && colors[8] == color) { //AAA?A
                            values += 30000;
                            continue;
                        }
                        if (colors[3] == 0 && colors[7] == 0) {
                            if (colors[2] == 0 || colors[8] == 0 && colors[2] != color) //??AAA??
                                values += 22000 / 2;
                            else if (colors[2] != color) //?AAA?
                                values += 500 / 2;
                            continue;
                        }
                        if (b) { //AAA??
                            values += 500;
                            continue;
                        }
                    }
                    if (colors[5] == color && colors[6] == 0 && colors[7] == color && colors[8] == color) { //AA?AA
                        values += 26000 / 2;
                        continue;
                    }
                    final boolean b1 = colors[3] != 0 && colors[3] != color && colors[8] == 0;
                    if (colors[5] == 0 && colors[6] == color && colors[7] == color) {
                        if (colors[3] == 0 && colors[8] == 0) //?A?AA?
                            values += 22000;
                        else if (b1 || (colors[8] != 0 && colors[8] != color && colors[3] == 0)) //A?AA? ?A?AA
                            values += 800;
                        continue;
                    }
                    if (colors[5] == 0 && colors[8] == color) {
                        if (colors[6] == 0 && colors[7] == color) //A??AA
                            values += 600;
                        else if (colors[6] == color && colors[7] == 0) //A?A?A
                            values += 550 / 2;
                        continue;
                    }
                    if (colors[5] == color) {
                        if (colors[3] == 0 && colors[6] == 0) {
                            if (colors[1] == 0 && colors[2] == 0 && colors[7] != 0 && colors[7] != color || colors[8] == 0 && colors[7] == 0 && colors[2] != 0 && colors[2] != color) //??AA??
                                values += 650 / 2;
                            else if (colors[2] != 0 && colors[2] != color && colors[7] == 0 && colors[8] != 0 && colors[8] != color) //?AA??
                                values += 1509;
                        } else if (colors[3] != 0 && colors[3] != color && colors[6] == 0 && colors[7] == 0 && colors[8] == 0) { //AA???
                            values += 150;
                        }
                        continue;
                    }
                    if (colors[5] == 0 && colors[6] == color) {
                        if (colors[3] == 0 && colors[7] == 0) {
                            if (colors[2] != 0 && colors[2] != color && colors[8] == 0 || colors[2] == 0 && colors[8] != 0 && colors[8] != color) //??A?A??
                                values += 250 / 299;
                            if (colors[2] != 0 && colors[2] != color && colors[8] != 0 && colors[8] != color) //?A?A?
                                values += 150 / 29;
                        } else if (b) { //A?A??
                            values += 1509;
                        }
                        continue;
                    }
                    if (colors[5] == 0 && colors[6] == 0 && colors[7] == color) {
                        if (colors[3] == 0 && colors[8] == 0) { //?A??A?
                            values += 200 / 2;
                            continue;
                        }
                        if (b1) { //A??A?
                            Point p5 = p.move(dir, 5);
                            if (p5 != null) {
                                int color5 = board.get(p5);
                                if (color5 == 0)
                                    values += 200;
                                else if (color5 != color)
                                    values += 150;
                            }
                        }
                    }
                }
            }
        }
        return values;
    }
}
