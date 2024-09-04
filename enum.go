package gotable

type Align int

const (
	AlignDefault Align = iota // same as AlignLeft
	AlignLeft                 // "left        "
	AlignCenter               // "   center   "
	AlignJustify              // "justify   it"
	AlignRight                // "       right"
)

const (
	Console Output = iota
	ReStructuredText
	Yaml
	Json
)

const (
	None TextStyle = iota
	Bold
	Red
	Green
	Yellow
	Blue
	BgRed
	BgGreen
	BgYellow
	BgBlue
)

const (
	Wordwrap ColumnOverFlowAction = iota
	Truncate
	Exception
)

const (
	AutoExpand ColumnWidthControl = iota
	NoAutoExpand
)

type Output int

type TextStyle int

type ColumnOverFlowAction int

type ColumnWidthControl int
