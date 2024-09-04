package gotable

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

type cellRenderer interface {
	render(c *Cell, w int, h int, o Output) ([]string, error)
	// stats returns the width and height that requires by the cell
	stats(c *Cell, wlimit int, o Output) (int, int, error)
}

type Cell struct {
	rawData      string
	leftPadding  string
	rightPadding string
	style        []TextStyle
	cellRenderer
}

func (a *Cell) String() string {
	return a.rawData
}

func (a *Cell) Style(tss ...TextStyle) {
	a.style = tss
}

func (a *Cell) Value(data any) {
	tmp := strings.Replace(fmt.Sprintf("%v", data), "\t", "    ", -1)
	tmp = strings.Replace(tmp, "\r\n", "\n", -1)
	a.rawData = tmp
}

func (a *Cell) Padding(l, r string) {
	a.leftPadding = l
	a.rightPadding = r
}

func (a *Cell) render(w int, h int, o Output) ([]string, error) {
	return a.cellRenderer.render(a, w, h, o)
}

func (a *Cell) stats(wlimit int, o Output) (int, int, error) {
	return a.cellRenderer.stats(a, wlimit, o)
}

func (a *Cell) sidePaddingWidth() int {
	return runewidth.StringWidth(a.leftPadding + a.rightPadding)
}

func (a *Cell) formatText(s string, o Output) string {
	switch o {
	case Console:
		return formatConsoleText(s, a.style...)
	case ReStructuredText:
		// TODO: not implemented
		return s
	default:
		return s
	}
}

func (a *Cell) textStylerWidth(o Output) int {
	switch o {
	case Console:
		return 0
	case ReStructuredText:
		// TODO: not implemented
		return 0
	default:
		return 0
	}
}

type DataCell struct {
	padding        rune
	align          Align
	overFlowAction ColumnOverFlowAction
	escapeLineFeed bool
}

func (a *DataCell) render(c *Cell, w int, h int, o Output) ([]string, error) {
	lines, err := a.splitIntoLines(c, w, o)
	if err != nil {
		return nil, err
	}
	wContent := w - c.sidePaddingWidth() - c.textStylerWidth(o)
	if w == 0 || wContent < 0 {
		return nil, ErrInvalidCellWidth
	}
	if h == 0 {
		return nil, ErrInvalidCellHeight
	}

	// format and render cell
	if len(lines) > h {
		return nil, ErrInsufficientColumnHeight
	}
	out := make([]string, h)
	for i := range out {
		tmp := ""
		if i < len(lines) {
			tmp = lines[i]
		}
		tmp = c.leftPadding + formatAlignment(tmp, wContent, a.padding, a.align) + c.rightPadding
		tmp = c.formatText(tmp, o)
		out[i] = tmp
	}
	return out, nil
}

// splitIntoLines splits cell data into multiple lines when overflow action is set to wordwrap
func (a *DataCell) splitIntoLines(c *Cell, wlimit int, o Output) ([]string, error) {
	content := c.String()
	if a.escapeLineFeed {
		content = escapeLineFeed(content)
	}
	wContent := runewidth.StringWidth(content)
	wContentLimit := wlimit - c.sidePaddingWidth() - c.textStylerWidth(o)
	if wlimit == 0 {
		wContentLimit = 0
	} else if wContentLimit < MinColumnWidth {
		return nil, fmt.Errorf("%w: no sufficient space after padding", ErrInsufficientColumnWidth)
	}

	// split string into string array according to cell overflow action
	contentLines := strings.Split(content, "\n")
	out := []string{}
	for _, cl := range contentLines {
		var lines []string
		switch a.overFlowAction {
		case Truncate:
			if wContentLimit != 0 && wContent > wContentLimit {
				wTailer := runewidth.StringWidth(UnfinishedCellTailer)
				wContentNew := wContentLimit - wTailer
				tmp := runewidth.Truncate(cl, wContentNew, "")
				wPadding := wContentNew - runewidth.StringWidth(tmp)
				tmp += strings.Repeat(string(a.padding), wPadding) + UnfinishedCellTailer
				cl = tmp
			}
			lines = []string{cl}
		case Wordwrap:
			lines = splitStringByWidth(cl, wContentLimit)
		case Exception:
			if wContent > wContentLimit {
				return nil, fmt.Errorf("%w: cell overflow action is set to exception", ErrInsufficientColumnWidth)
			}
			lines = []string{cl}
		}
		out = append(out, lines...)
	}
	return out, nil
}

func (a *DataCell) stats(c *Cell, wlimit int, o Output) (int, int, error) {
	lines, err := a.splitIntoLines(c, wlimit, o)
	if err != nil {
		return 0, 0, err
	}
	height := len(lines)
	if wlimit != 0 {
		return wlimit, height, nil
	}

	// use the longest line as the width of the cell when wlimit is set to 0
	width := 0
	for _, l := range lines {
		if tmp := runewidth.StringWidth(l); tmp > width {
			width = tmp
		}
	}
	width += c.sidePaddingWidth() + c.textStylerWidth(o)
	return width, height, nil
}

type TreePathCell struct {
	style *TreePathStyle
}

func (a *TreePathCell) render(c *Cell, w int, h int, o Output) ([]string, error) {
	out := make([]string, h)
	for i := 0; i < h; i++ {
		if i == 0 {
			out[i] = c.String()
		} else {
			out[i] = a.style.ReplacePathAsExtention(c.String())
		}
		out[i] = c.leftPadding + out[i] + c.rightPadding
		out[i] = c.formatText(out[i], o)
	}
	return out, nil
}

func (a *TreePathCell) stats(c *Cell, wlimit int, o Output) (int, int, error) {
	// use utf8.RuneCountInString instead of runewidth.StringWidth since some character
	rw := &runewidth.Condition{
		EastAsianWidth:     false,
		StrictEmojiNeutral: true,
	}
	w := c.sidePaddingWidth() + c.textStylerWidth(o) + rw.StringWidth(c.String())
	return w, 1, nil
}
