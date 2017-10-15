package codegen

// Label is a label in ir
type Label struct {
	Name string
}

func (label *Label) Line() string {
	return label.Name + ":"
}
