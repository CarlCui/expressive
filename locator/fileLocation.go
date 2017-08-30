package locator

import (
	"path"
	"strconv"
)

type FileLocation struct {
	Row      int
	Col      int
	FileName string
	DirName  string
}

func (loc *FileLocation) Locate() string {
	return "file " + path.Join(loc.DirName, loc.FileName) +
		": row " + strconv.Itoa(loc.Row) + ", column " + strconv.Itoa(loc.Col)
}
