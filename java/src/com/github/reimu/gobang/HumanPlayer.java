package com.github.reimu.gobang;

import javax.swing.*;
import java.awt.*;

@SuppressWarnings("serial")
public class HumanPlayer extends JFrame implements Player {
	private static final String BLACK = "●";
	private static final String BLACK1 = "◆";
	private static final String WHITE = "○";
	private static final String WHITE1 = "◎";
	private final byte[][] board = new byte[Constant.MAX_LEN][Constant.MAX_LEN];
	private final JButton[][] btn = new JButton[Constant.MAX_LEN][Constant.MAX_LEN];
	private boolean isTurn = false;
	private com.github.reimu.gobang.Point p;
	private final int color;

	public HumanPlayer(int color) {
		setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
		this.color = color;
		setLayout(new GridLayout(Constant.MAX_LEN, Constant.MAX_LEN));
		setBounds(10, 10, 600, 600);
		for (int i = 0; i < Constant.MAX_LEN; i++) {
			for (int j = 0; j < Constant.MAX_LEN; j++) {
				final int ii = i, jj = j;
				btn[i][j] = new JButton(" ");
				btn[i][j].setMargin(new Insets(0, 0, 0, 0));
				btn[i][j].setFont(new Font("宋体", Font.PLAIN, 32));
				btn[i][j].addActionListener(e -> {
					if (isTurn) {
						synchronized (HumanPlayer.this) {
							if (isTurn && board[ii][jj] == 0) {
								board[ii][jj] = (byte) HumanPlayer.this.color;
								btn[ii][jj].setText(HumanPlayer.this.color == 1 ? BLACK1 : WHITE1);
								if (p != null)
									btn[p.y][p.x].setText(HumanPlayer.this.color == 2 ? BLACK : WHITE);
								p = new com.github.reimu.gobang.Point(jj, ii);
								isTurn = false;
								HumanPlayer.this.notify();
							}
						}
					}
				});
				add(btn[i][j]);
			}
		}
		isTurn = false;
		setVisible(true);
	}
	
	@Override
	public com.github.reimu.gobang.Point play() {
		synchronized (this) {
			isTurn = true;
			try {
				wait();
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
		}
		return p;
	}

	@Override
	public void display(Point p) {
		if (board[p.y][p.x] != 0)
			throw new IllegalArgumentException(p.toString() + board[p.y][p.x]);
		board[p.y][p.x] = (byte) (3 - color);
		btn[p.y][p.x].setText(color == 2 ? BLACK1 : WHITE1);
		if (this.p != null)
			btn[this.p.y][this.p.x].setText(color == 1 ? BLACK : WHITE);
		this.p = p;
	}
}
