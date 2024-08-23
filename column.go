package gotable

type columnCellMaker interface {
	newCell(col *Column, v any) *Cell
}

type Column struct {
	name             string
	hidden           bool
	widthLimit       int
	autoWidthControl bool
	leftPadding      string
	rightPadding     string
	padding          rune
	header           ColumnStyle
	body             ColumnStyle
	columnCellMaker
}

func (a *Column) Name() string {
	return a.name
}

func (a *Column) Hidden(b bool) *Column {
	a.hidden = b
	return a
}

func (a *Column) Width(limit int, autoControl bool) *Column {
	a.widthLimit = limit
	a.autoWidthControl = autoControl
	return a
}

func (a *Column) Padding(p rune, l, r string) *Column {
	a.padding = p
	a.leftPadding = l
	a.rightPadding = r
	return a
}

func (a *Column) HeaderStyle(sty *ColumnStyle) *Column {
	a.header = *sty
	return a
}

func (a *Column) BodyStyle(sty *ColumnStyle) *Column {
	a.body = *sty
	return a
}

func (a *Column) newHeader() *Cell {
	if a.header.overFlowAction == Wordwrap {
		a.header.escapeLineFeed = false
	}
	cell := &Cell{
		leftPadding:  a.leftPadding,
		rightPadding: a.rightPadding,
		style:        a.header.text,
		cellRenderer: &DataCell{
			padding:        a.padding,
			overFlowAction: a.header.overFlowAction,
			escapeLineFeed: a.header.escapeLineFeed,
			align:          a.header.align,
		},
	}
	cell.Value(a.name)
	return cell
}

func (a *Column) newCell(v any) *Cell {
	return a.columnCellMaker.newCell(a, v)
}

type TreePathColumn struct {
	sty TreePathStyle
}

func (a *TreePathColumn) newCell(col *Column, v any) *Cell {
	cell := &Cell{
		leftPadding:  col.leftPadding,
		rightPadding: col.rightPadding,
		cellRenderer: &TreePathCell{
			style: &a.sty,
		},
	}
	cell.Value(v)
	return cell
}

type StandardColumn struct {
}

func (a *StandardColumn) newCell(col *Column, v any) *Cell {
	cell := &Cell{
		leftPadding:  col.leftPadding,
		rightPadding: col.rightPadding,
		style:        col.body.text,
		cellRenderer: &DataCell{
			padding:        col.padding,
			overFlowAction: col.body.overFlowAction,
			escapeLineFeed: col.body.escapeLineFeed,
			align:          col.body.align,
		},
	}
	cell.Value(v)
	return cell
}

func NewStandardColumn(name string) *Column {
	col := &Column{
		name:             name,
		hidden:           false,
		widthLimit:       0,
		autoWidthControl: true,
		leftPadding:      " ",
		rightPadding:     " ",
		padding:          ' ',
		header:           *DefauleHeaderStyle(),
		body:             *DefauleBodyStyle(),
		columnCellMaker:  &StandardColumn{},
	}
	return col
}

func NewTreePathColumn(sty TreePathStyle) *Column {
	header := DefauleHeaderStyle()
	header.text = sty.header
	body := DefauleTreeStyle()
	body.text = sty.body
	col := &Column{
		name:             sty.Name,
		hidden:           false,
		widthLimit:       0,
		autoWidthControl: false,
		leftPadding:      " ",
		rightPadding:     " ",
		padding:          ' ',
		header:           *header,
		body:             *body,
		columnCellMaker: &TreePathColumn{
			sty,
		},
	}
	return col
}

func DefauleHeaderStyle() *ColumnStyle {
	return &ColumnStyle{
		overFlowAction: Wordwrap,
		escapeLineFeed: false,
		align:          AlignCenter,
		text: []TextStyle{
			Bold,
		},
	}
}

func DefauleBodyStyle() *ColumnStyle {
	return &ColumnStyle{
		overFlowAction: Wordwrap,
		escapeLineFeed: false,
		align:          AlignLeft,
		text:           []TextStyle{},
	}
}

func DefauleTreeStyle() *ColumnStyle {
	return &ColumnStyle{
		overFlowAction: Exception,
		escapeLineFeed: false,
		align:          AlignLeft,
		text:           []TextStyle{},
	}
}

type ColumnStyle struct {
	overFlowAction ColumnOverFlowAction
	escapeLineFeed bool
	align          Align
	text           []TextStyle
}

func (a *ColumnStyle) OverFlowAction(ofa ColumnOverFlowAction) *ColumnStyle {
	a.overFlowAction = ofa
	return a
}

func (a *ColumnStyle) EscapeLineFeed(b bool) *ColumnStyle {
	a.escapeLineFeed = b
	return a
}

func (a *ColumnStyle) Align(al Align) *ColumnStyle {
	a.align = al
	return a
}

func (a *ColumnStyle) Text(st ...TextStyle) *ColumnStyle {
	a.text = st
	return a
}
