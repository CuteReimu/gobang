package main

type player interface {
	color() playerColor
	play() (point, error)
	display(p point) error
}

const (
	colorEmpty playerColor = iota
	colorBlack
	colorWhite
)

type playerColor int8

func (c playerColor) String() string {
	switch c {
	case colorEmpty:
		return "无"
	case colorBlack:
		return "黑"
	case colorWhite:
		return "白"
	}
	panic("unreachable")
}

func (c playerColor) getString0() string {
	switch c {
	case colorBlack:
		return "●"
	case colorWhite:
		return "○"
	}
	panic("unreachable")
}

func (c playerColor) getString1() string {
	switch c {
	case colorBlack:
		return "◆"
	case colorWhite:
		return "◎"
	}
	panic("unreachable")
}

func (c playerColor) conversion() playerColor {
	return 3 - c
}
