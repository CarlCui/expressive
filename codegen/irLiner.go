package codegen

// IrLiner represents one line of code in IR
type IrLiner interface {
	Line() string
}
