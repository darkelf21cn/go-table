package gotable

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

type TreeStatistics struct {
	MaxDepth int
}

func escapeLineFeed(s string) string {
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

func getTreeStatistics(tr []TreeNodeReader) TreeStatistics {
	ts := TreeStatistics{}
	for i, nn := range tr {
		tmp := _getTreeStatistics(nn)
		if i == 0 {
			ts = tmp
			continue
		}
		if tmp.MaxDepth > ts.MaxDepth {
			ts.MaxDepth = tmp.MaxDepth
		}
	}
	return ts
}

func _getTreeStatistics(tn TreeNodeReader) TreeStatistics {
	ts := TreeStatistics{
		MaxDepth: 0,
	}
	for _, nn := range tn.Children() {
		tmp := _getTreeStatistics(nn)
		if tmp.MaxDepth > ts.MaxDepth {
			ts.MaxDepth = tmp.MaxDepth
		}
	}
	ts.MaxDepth++
	return ts
}

func splitStringByWidth(s string, wLimit int) []string {
	lines := strings.Split(s, "\n")
	str := lines[0]
	wContent := runewidth.StringWidth(str)
	var out, remain string
	if wLimit == 0 || wContent <= wLimit {
		out = str
		remain = strings.Join(lines[1:], "\n")
	} else {
		out = runewidth.Truncate(str, wLimit, "")
		remain = strings.Replace(str, out, "", 1)
		remain += "\n" + strings.Join(lines[1:], "\n")
	}
	if remain == "" {
		return []string{out}
	} else {
		return append([]string{out}, splitStringByWidth(remain, wLimit)...)
	}
}

func formatConsoleText(str string, tss ...TextStyle) string {
	if len(tss) == 0 {
		return str
	}
	out := str
	fmtColor := "%s"
	fmtBgColor := "%s"
	for _, ts := range tss {
		switch ts {
		case None:
			return str
		case Bold:
			out = fmt.Sprintf("\033[1m%s\033[22m", out)
		case Red:
			fmtColor = "\033[31m%s\033[0m"
		case Green:
			fmtColor = "\033[32m%s\033[0m"
		case Yellow:
			fmtColor = "\033[33m%s\033[0m"
		case Blue:
			fmtColor = "\033[34m%s\033[0m"
		case BgRed:
			fmtBgColor = "\033[41m%s\033[0m"
		case BgGreen:
			fmtBgColor = "\033[42m%s\033[0m"
		case BgYellow:
			fmtBgColor = "\033[43m%s\033[0m"
		case BgBlue:
			fmtBgColor = "\033[44m%s\033[0m"
		}
	}
	out = fmt.Sprintf(fmtColor, out)
	out = fmt.Sprintf(fmtBgColor, out)
	return out
}

func formatAlignment(s string, w int, padding rune, align Align) string {
	ws := runewidth.StringWidth(s)
	// no check on negative number of padCount since it should be handled before invoking this function
	padCount := w - ws
	padStr := string(padding)
	out := ""
	if s == "" {
		return strings.Repeat(padStr, padCount)
	}
	switch align {
	case AlignDefault, AlignLeft:
		out = s + strings.Repeat(padStr, padCount)
	case AlignRight:
		out = strings.Repeat(padStr, padCount) + s
	case AlignCenter:
		padLeft := padCount / 2
		padRight := padCount - padLeft
		out = strings.Repeat(padStr, padLeft) + s + strings.Repeat(padStr, padRight)
	case AlignJustify:
		tmp := strings.Split(s, padStr)
		for padCount > 0 {
			for i := 0; i < len(tmp)-1; i++ {
				tmp[i] += padStr
				padCount--
				if padCount == 0 {
					break
				}
			}
		}
		out = strings.Join(tmp, padStr)
	}
	return out
}
