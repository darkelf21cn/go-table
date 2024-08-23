package gotable

import (
	"strings"
)

type TableLayout struct {
	HeaderTopLeft          rune
	HeaderTopRight         rune
	HeaderTopSeparator     rune
	HeaderTopHorizontal    rune
	HeaderLeft             rune
	HeaderRight            rune
	HeaderSeparator        rune
	HeaderBottomLeft       rune
	HeaderBottomRight      rune
	HeaderBottomSeparator  rune
	HeaderBottomHorizontal rune
	BodyTopLeft            rune
	BodyTopRight           rune
	BodyTopSeparator       rune
	BodyTopHorizontal      rune
	BodyBottomLeft         rune
	BodyBottomRight        rune
	BodyBottomSeparator    rune
	BodyBottomHorizontal   rune
	RowLeft                rune
	RowRight               rune
	RowSeparator           rune
	RowHorizontal          rune
	ColumnSeparator        rune
	ColumnPaddingLeft      string
	ColumnPaddingRight     string
	CellPadding            string

	ShowHeader             bool
	ShowHeaderTopBorder    bool
	ShowHeaderBottemBorder bool
	ShowBodyTopBorder      bool
	ShowBodyBottomBorder   bool
	ShowSideBorder         bool
	ShowColumnSeparator    bool
	ShowRowSeparator       bool

	Width int
}

func (a *TableLayout) HideHeader() *TableLayout {
	a.ShowHeader = false
	a.ShowHeaderTopBorder = false
	a.ShowHeaderBottemBorder = false
	// when header is off, top and bottom boder of body should be keep consistent
	a.ShowBodyTopBorder = a.ShowBodyBottomBorder
	return a
}

func (a *TableLayout) HideOutterBorder() *TableLayout {
	if a.ShowHeader {
		a.ShowHeaderTopBorder = false
	} else {
		a.ShowBodyTopBorder = false
	}
	a.ShowSideBorder = false
	a.ShowBodyBottomBorder = false
	return a
}

func (a *TableLayout) SplitHeaderAndBody() *TableLayout {
	a.HeaderBottomSeparator = a.BodyBottomSeparator
	a.HeaderBottomLeft = a.BodyBottomLeft
	a.HeaderBottomRight = a.BodyBottomRight
	a.BodyTopSeparator = a.HeaderTopSeparator
	a.BodyTopLeft = a.HeaderTopLeft
	a.BodyTopRight = a.HeaderTopRight
	a.ShowHeader = true
	a.ShowHeaderTopBorder = true
	a.ShowHeaderBottemBorder = true
	a.ShowBodyTopBorder = true
	a.ShowBodyBottomBorder = true
	a.ShowSideBorder = true
	return a
}

func LightTableLayout() *TableLayout {
	return &TableLayout{
		HeaderTopLeft:          '┌',
		HeaderTopRight:         '┐',
		HeaderTopSeparator:     '┬',
		HeaderTopHorizontal:    '─',
		HeaderLeft:             '│',
		HeaderRight:            '│',
		HeaderSeparator:        '│',
		HeaderBottomLeft:       '├',
		HeaderBottomRight:      '┤',
		HeaderBottomSeparator:  '┼',
		HeaderBottomHorizontal: '─',
		BodyTopLeft:            '┌',
		BodyTopRight:           '┐',
		BodyTopSeparator:       '┬',
		BodyTopHorizontal:      '─',
		BodyBottomLeft:         '└',
		BodyBottomRight:        '┘',
		BodyBottomSeparator:    '┴',
		BodyBottomHorizontal:   '─',
		RowLeft:                '│',
		RowRight:               '│',
		RowSeparator:           '┼',
		RowHorizontal:          '─',
		ColumnSeparator:        '│',
		ColumnPaddingLeft:      " ",
		ColumnPaddingRight:     " ",
		CellPadding:            " ",
		ShowHeader:             true,
		ShowHeaderTopBorder:    true,
		ShowHeaderBottemBorder: true,
		ShowBodyTopBorder:      false,
		ShowBodyBottomBorder:   true,
		ShowSideBorder:         true,
		ShowColumnSeparator:    true,
		ShowRowSeparator:       false,
	}
}

func DefaultTableLayout() *TableLayout {
	return &TableLayout{
		HeaderTopLeft:          '+',
		HeaderTopRight:         '+',
		HeaderTopSeparator:     '+',
		HeaderTopHorizontal:    '-',
		HeaderLeft:             '|',
		HeaderRight:            '|',
		HeaderSeparator:        '|',
		HeaderBottomLeft:       '+',
		HeaderBottomRight:      '+',
		HeaderBottomSeparator:  '+',
		HeaderBottomHorizontal: '-',
		BodyTopLeft:            '+',
		BodyTopRight:           '+',
		BodyTopSeparator:       '+',
		BodyTopHorizontal:      '-',
		BodyBottomLeft:         '+',
		BodyBottomRight:        '+',
		BodyBottomSeparator:    '+',
		BodyBottomHorizontal:   '-',
		RowLeft:                '|',
		RowRight:               '|',
		RowSeparator:           '+',
		RowHorizontal:          '-',
		ColumnSeparator:        '|',
		ColumnPaddingLeft:      " ",
		ColumnPaddingRight:     " ",
		CellPadding:            " ",
		ShowHeader:             true,
		ShowHeaderTopBorder:    true,
		ShowHeaderBottemBorder: true,
		ShowBodyTopBorder:      false,
		ShowBodyBottomBorder:   true,
		ShowSideBorder:         true,
		ShowColumnSeparator:    true,
		ShowRowSeparator:       false,
	}
}

type TreePathStyle struct {
	Name          string
	Root          string
	Middle        string
	Terminal      string
	PadLine       string
	PadBlank      string
	Children      string
	PrefixLeveled string
	PrefixBlank   string

	header []TextStyle
	body   []TextStyle
}

func (a *TreePathStyle) Header(ss ...TextStyle) *TreePathStyle {
	a.header = ss
	return a
}

func (a *TreePathStyle) Body(ss ...TextStyle) *TreePathStyle {
	a.body = ss
	return a
}

func (a TreePathStyle) ReplacePathAsExtention(s string) string {
	out := strings.Replace(s, a.Children, a.PrefixLeveled, -1)
	out = strings.Replace(out, a.Middle, a.PrefixLeveled, -1)
	out = strings.Replace(out, a.Terminal, a.PrefixBlank, -1)
	out = strings.Replace(out, a.Root, a.PrefixBlank, -1)
	out = strings.Replace(out, a.PadLine, a.PadBlank, -1)
	return out
}

func DefaultTreePathStyle() *TreePathStyle {
	return &TreePathStyle{
		Name:          "Path",
		Root:          ">-",
		Middle:        "+-",
		Terminal:      "\\-",
		Children:      "+-",
		PrefixLeveled: "| ",
		PrefixBlank:   "  ",
		PadLine:       "-",
		PadBlank:      " ",
		header:        []TextStyle{},
		body:          []TextStyle{},
	}
}

func LightTreePathStyle() *TreePathStyle {
	return &TreePathStyle{
		Name:          "Path",
		Root:          "□─",
		Middle:        "├─",
		Terminal:      "└─",
		Children:      "┬─",
		PrefixLeveled: "│ ",
		PrefixBlank:   "  ",
		PadLine:       "─",
		PadBlank:      " ",
		header:        []TextStyle{},
		body:          []TextStyle{},
	}
}
