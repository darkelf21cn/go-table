package gotable

import "errors"

var (
	ErrColumnAlreadyExist       = errors.New("column already exist")
	ErrColumnNotExist           = errors.New("column does not exist")
	ErrEnforcingTableWidth      = errors.New("enforcing table width failed")
	ErrFieldIsMissing           = errors.New("field is missing")
	ErrInsufficientColumnHeight = errors.New("insufficient column height")
	ErrInsufficientColumnWidth  = errors.New("insufficient column width")
	ErrInvalidCellHeight        = errors.New("invalid cell height")
	ErrInvalidCellWidth         = errors.New("invalid cell width")
	ErrNoAdjustableColumn       = errors.New("no adjustable column")
	ErrRenderTableFailed        = errors.New("render table failed")
	ErrTableNotEmpty            = errors.New("table is not empty")
)
