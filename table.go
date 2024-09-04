package gotable

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Table struct {
	Layout TableLayout

	columns []*Column
	rows    []Row
	stats   TableStats
	colMap  map[string]int
}

type TableStats struct {
	ColumnWidths []int
	RowHeights   []int
	HeaderHeight int
}

func NewTable(l *TableLayout) *Table {
	var tmp *TableLayout
	if l != nil {
		tmp = l
	} else {
		tmp = DefaultTableLayout()
	}
	t := &Table{
		Layout:  *tmp,
		columns: []*Column{},
		rows:    []Row{},
		colMap:  map[string]int{},
	}
	return t
}

func (a *Table) AppendColumn(cs ...*Column) (*Table, error) {
	for _, c := range cs {
		_, err := a.appendColumn(c)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// AppendRowM add construct a table row from a Row and append to the table
func (a *Table) AppendRowM(m map[string]any) error {
	row, err := a.convFieldsToRow(m)
	if err != nil {
		return err
	}
	a.rows = append(a.rows, row)
	return nil
}

// AppendRow add construct a table row from []any and append to the table
func (a *Table) AppendRow(c ...any) error {
	colCount := len(a.columns)
	if len(c) != colCount {
		return fmt.Errorf("table has %d columns but %d is given", colCount, len(c))
	}
	cells := make([](*Cell), colCount)
	for i, col := range a.columns {
		cells[i] = col.newCell(c[i])
	}
	a.rows = append(a.rows, cells)
	return nil
}

func (a *Table) AppendTrees(sty TreePathStyle, ns ...TreeNodeReader) error {
	a.setTreePathColumn(sty)
	ts := getTreeStatistics(ns)
	rows := []Row{}
	for _, n := range ns {
		tmp, err := a.convTree2Rows(sty, n, ts.MaxDepth, "", false, true)
		if err != nil {
			return err
		}
		rows = append(rows, tmp...)
	}
	a.rows = append(a.rows, rows...)
	return nil
}

func (a *Table) Cell(col, row int) *Cell {
	return a.rows[row][col]
}

func (a *Table) Render(o Output) (string, error) {
	out := ""
	err := func() error {
		err := a.enforceWidth(o)
		if err != nil {
			return err
		}
		header, err := a.renderHeader(o)
		if err != nil {
			return err
		}
		body, err := a.renderBody(o)
		if err != nil {
			return err
		}
		out = header + body
		return nil
	}()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRenderTableFailed, err)
	}
	return out, nil
}

func (a *Table) ResetData() {
	a.rows = []Row{}
	a.stats = TableStats{}
}

func (a *Table) GetColumn(name string) (*Column, error) {
	cIdx, ok := a.colMap[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrColumnNotExist, name)
	}
	return a.columns[cIdx], nil
}

func (a *Table) appendColumn(c *Column) (*Table, error) {
	if _, ok := a.colMap[c.Name()]; ok {
		return nil, fmt.Errorf("%w: %s", ErrColumnAlreadyExist, c.Name())
	}
	cIdx := len(a.columns)
	a.columns = append(a.columns, c)
	a.colMap[c.Name()] = cIdx
	return a, nil
}

func (a *Table) setTreePathColumn(sty TreePathStyle) *Column {
	col := NewTreePathColumn(sty)
	if len(a.columns) == 0 {
		a.columns = append(a.columns, col)
	} else {
		switch a.columns[0].columnCellMaker.(type) {
		case *TreePathColumn:
			a.columns[0] = col
		default:
			a.columns = append([]*Column{col}, a.columns...)
		}
	}
	return col
}

func (a *Table) updateStatistics(o Output) error {
	stats := TableStats{
		ColumnWidths: make([]int, len(a.columns)),
		RowHeights:   make([]int, len(a.rows)),
		HeaderHeight: 0,
	}
	for i, col := range a.columns {
		w, h, err := col.newHeader().stats(a.columns[i].widthLimit, o)
		if err != nil {
			return err
		}
		if h > stats.HeaderHeight {
			stats.HeaderHeight = h
		}
		if w > stats.ColumnWidths[i] {
			stats.ColumnWidths[i] = w
		}
	}
	for ri, row := range a.rows {
		for ci, cell := range row {
			w, h, err := cell.stats(a.columns[ci].widthLimit, o)
			if err != nil {
				return err
			}
			if h > stats.RowHeights[ri] {
				stats.RowHeights[ri] = h
			}
			if w > stats.ColumnWidths[ci] {
				stats.ColumnWidths[ci] = w
			}
		}
	}
	a.stats = stats
	return nil
}

// enforceWidth adjusts widthlimit for each dynamic columns to fit the table width
func (a *Table) enforceWidth(o Output) error {
	// update statistics to get original column width and row height
	err := a.updateStatistics(o)
	if err != nil {
		return err
	}

	l := a.Layout
	if l.Width == 0 {
		return nil
	}
	originalWidth := 0
	// left and right separactors
	if l.ShowSideBorder {
		originalWidth += 2
	}
	colCount := 0
	for i, col := range a.columns {
		if col.hidden {
			continue
		}
		colCount++
		originalWidth += a.stats.ColumnWidths[i]
	}
	// column separactors
	if a.Layout.ShowColumnSeparator {
		originalWidth += colCount - 1
	}
	// do nothing when table width equals to the expected size
	if originalWidth == a.Layout.Width {
		return nil
	}

	// skip columns which meet conditions below
	// column is hidden
	// column autoWidthControl is disabled
	// column width is equal or less than AdjustableColumnMinWidth
	colIndexes := []int{}
	for i, col := range a.columns {
		if col.hidden {
			continue
		}
		if !col.autoWidthControl {
			continue
		}
		if a.stats.ColumnWidths[i] <= AdjustableColumnMinWidth {
			continue
		}
		colIndexes = append(colIndexes, i)
	}
	if len(colIndexes) == 0 {
		return fmt.Errorf("%w: %w", ErrEnforcingTableWidth, ErrNoAdjustableColumn)
	}

	if originalWidth < a.Layout.Width {
		// increase column width
		widthToAdd := a.Layout.Width - originalWidth
		widthToAddPerCol := widthToAdd / len(colIndexes)
		widthToAddRemain := widthToAdd % len(colIndexes)
		for i, ci := range colIndexes {
			w := a.stats.ColumnWidths[ci] + widthToAddPerCol
			if i < widthToAddRemain {
				w += 1
			}
			a.columns[ci].Width(w, false)
		}
	} else {
		// reduce column width for wordwrap enabled columns
		tmp := originalWidth
		for i, colIdx := range colIndexes {
			if !a.columns[colIdx].autoWidthControl {
				colIndexes = append(colIndexes[:i], colIndexes[i+1:]...)
				continue
			}
			tmp -= a.stats.ColumnWidths[colIdx]
		}
		if tmp+(AdjustableColumnMinWidth*len(colIndexes)) > a.Layout.Width {
			return fmt.Errorf("enforcing table width to %d is not possible since rows are too long", a.Layout.Width)
		}
		widthPerCol := (a.Layout.Width - tmp) / len(colIndexes)
		widthLeft := (a.Layout.Width - tmp) % len(colIndexes)
		for i, colIdx := range colIndexes {
			w := widthPerCol
			if i < widthLeft {
				w += 1
			}
			a.columns[colIdx].Width(w, false)
		}
	}

	// table statistics has to be updated since width limit has changed
	err = a.updateStatistics(o)
	if err != nil {
		return err
	}
	return nil
}

func (a *Table) renderHeader(o Output) (string, error) {
	out := a.renderHorizontal("HeaderTop")
	row := make(Row, len(a.columns))
	for i, col := range a.columns {
		row[i] = col.newHeader()
	}
	if a.Layout.ShowHeader {
		tmp, err := a.renderRow(row, a.stats.HeaderHeight, o)
		if err != nil {
			return "", err
		}
		out += tmp
	}
	out += a.renderHorizontal("HeaderBottom")
	return out, nil
}

func (a *Table) renderBody(o Output) (string, error) {
	out := a.renderHorizontal("BodyTop")
	for i := 0; i < len(a.rows); i++ {
		tmp, err := a.renderRow(a.rows[i], a.stats.RowHeights[i], o)
		if err != nil {
			return "", err
		}
		out += tmp
		if i != len(a.rows)-1 {
			out += a.renderHorizontal("Row")
		}
	}
	out += a.renderHorizontal("BodyBottom")
	return out, nil
}

func (a *Table) renderRow(cells Row, h int, o Output) (string, error) {
	out := ""
	colAndRows := make([][]string, 0)
	for i, col := range a.columns {
		if col.hidden {
			continue
		}
		tmp, err := cells[i].render(a.stats.ColumnWidths[i], h, o)
		if err != nil {
			return "", err
		}
		colAndRows = append(colAndRows, tmp)
	}
	colSep := string(a.Layout.ColumnSeparator)
	if !a.Layout.ShowColumnSeparator {
		colSep = ""
	}

	strLeft, strRight := "", ""
	if a.Layout.ShowSideBorder {
		strLeft, strRight = string(a.Layout.RowLeft), string(a.Layout.RowRight)
	}
	for r := 0; r < h; r++ {
		lineItems := []string{}
		for c := range colAndRows {
			lineItems = append(lineItems, colAndRows[c][r])
		}
		line := strLeft + strings.Join(lineItems, colSep) + strRight + "\n"
		out += line
	}
	return out, nil
}

func (a *Table) renderHorizontal(t string) string {
	l := a.Layout
	switch t {
	case "HeaderTop":
		return a._renderHorizontal(l.ShowHeaderTopBorder, l.HeaderTopLeft, l.HeaderTopRight, l.HeaderTopSeparator, l.HeaderTopHorizontal)
	case "HeaderBottom":
		return a._renderHorizontal(l.ShowHeaderBottemBorder, l.HeaderBottomLeft, l.HeaderBottomRight, l.HeaderBottomSeparator, l.HeaderBottomHorizontal)
	case "BodyTop":
		return a._renderHorizontal(l.ShowBodyTopBorder, l.BodyTopLeft, l.BodyTopRight, l.BodyTopSeparator, l.BodyTopHorizontal)
	case "BodyBottom":
		return a._renderHorizontal(l.ShowBodyBottomBorder, l.BodyBottomLeft, l.BodyBottomRight, l.BodyBottomSeparator, l.BodyBottomHorizontal)
	case "Row":
		return a._renderHorizontal(l.ShowRowSeparator, l.RowLeft, l.RowRight, l.RowSeparator, l.RowHorizontal)
	default:
		return ""
	}
}

func (a *Table) _renderHorizontal(show bool, left rune, right rune, separator rune, horizontal rune) string {
	if !show {
		return ""
	}
	out := ""
	tmp := []string{}
	for i, col := range a.columns {
		if col.hidden {
			continue
		}
		wCol := a.stats.ColumnWidths[i]
		tmp = append(tmp, strings.Repeat(string(horizontal), wCol))
	}
	strColSep, strLeft, strRight := "", "", ""
	if a.Layout.ShowColumnSeparator {
		strColSep = string(separator)
	}
	if a.Layout.ShowSideBorder {
		strLeft, strRight = string(left), string(right)
	}
	out = fmt.Sprintf("%s%s%s\n", strLeft, strings.Join(tmp, strColSep), strRight)
	return out
}

func (a *Table) convTree2Rows(sty TreePathStyle, node TreeNodeReader, maxDeepth int, prefix string, islast bool, isroot bool) ([]Row, error) {
	// generate data fields of the row
	fields := node.Fields()
	fields[sty.Name] = ""
	row, err := a.convFieldsToRow(fields)
	if err != nil {
		return nil, err
	}

	// generate tree path
	path := ""
	if isroot {
		path = sty.Root
		prefix += sty.PrefixBlank
	} else {
		path += prefix
		if islast {
			path += sty.Terminal
			prefix += sty.PrefixBlank
		} else {
			path += sty.Middle
			prefix += sty.PrefixLeveled
		}
	}
	if len(node.Children()) > 0 {
		path += sty.Children
	}
	pathWidth := 1 + (maxDeepth * 2)
	pad := strings.Repeat(sty.PadLine, pathWidth-utf8.RuneCountInString(path))
	path += pad
	row[0].Value(path)

	out := []Row{row}
	// generate rows for children nodes
	for i, cld := range node.Children() {
		tmp, err := a.convTree2Rows(sty, cld, maxDeepth, prefix, i == len(node.Children())-1, false)
		if err != nil {
			return nil, err
		}
		out = append(out, tmp...)
	}
	return out, nil
}

func (a *Table) convFieldsToRow(fields map[string]any) (Row, error) {
	out := make(Row, len(a.columns))
	for i, col := range a.columns {
		v, ok := fields[col.name]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrFieldIsMissing, col.name)
		}
		c := col.newCell(v)
		out[i] = c
	}
	return out, nil
}
