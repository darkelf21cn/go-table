package gotable

type TreeNodeReader interface {
	Children() []TreeNodeReader
	Fields() map[string]any
}

type Row []*Cell
