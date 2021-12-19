package com.github.reimu.gobang;

import java.awt.Font;
import java.awt.GridLayout;
import java.awt.Insets;

import javax.swing.JButton;
import javax.swing.JFrame;

@SuppressWarnings("serial")
public class HumanWatcher extends JFrame {
	private static final String BLACK = "●";
	private static final String BLACK1 = "◆";
	private static final String WHITE = "○";
	private static final String WHITE1 = "◎";
	private byte[][] board = new byte[Constant.MAX_LEN][Constant.MAX_LEN];
	private JButton[][] btn = new JButton[Constant.MAX_LEN][Constant.MAX_LEN];
	private Point lastPoint;

	public HumanWatcher() {
		setLayout(new GridLayout(Constant.MAX_LEN, Constant.MAX_LEN));
		setBounds(10, 10, 600, 600);
		for (int i = 0; i < Constant.MAX_LEN; i++) {
			for (int j = 0; j < Constant.MAX_LEN; j++) {
				btn[i][j] = new JButton(" ");
				btn[i][j].setFont(new Font("宋体", Font.PLAIN, 32));
				btn[i][j].setMargin(new Insets(0, 0, 0, 0));
				add(btn[i][j]);
			}
		}
		setVisible(true);
	}
	
	public void display(Point p, int color) {
		if (board[p.y][p.x] != 0)
			throw new IllegalArgumentException(p.toString() + board[p.y][p.x]);
		board[p.y][p.x] = (byte) color;
		btn[p.y][p.x].setText(color == 1 ? BLACK1 : WHITE1);
		if (lastPoint != null)
			btn[lastPoint.y][lastPoint.x].setText(color == 2 ? BLACK : WHITE);
		lastPoint = p;
	}
	
}
