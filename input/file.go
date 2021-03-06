package input

import (
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"unicode/utf8"

	"github.com/carlcui/expressive/locator"
)

// File encapsulates all operations related to file IO
type File struct {
	curRow    int
	curColumn int

	curRead int

	src      []byte
	filename string
	dirname  string
}

// Init initializes File by reading and storing the contents in the file to src
func (file *File) Init(dirname, filename string) {
	filePath := path.Join(dirname, filename)

	src, err := ioutil.ReadFile(filePath)

	if err != nil {
		// TODO: proper error reporting
		panic(err)
	}

	file.curColumn = 0
	file.curColumn = 0

	file.curRead = 0

	file.src = src
	file.filename = filename
	file.dirname = dirname
}

// NextRune returns the next rune. The user needs to check IsEOF before calling this function.
func (file *File) NextRune() rune {
	if file.IsEOF() {
		panic("EOF in " + file.dirname + "/" + file.filename)
	}

	r, size := utf8.DecodeRune(file.src[file.curRead:])

	if r == utf8.RuneError {
		file.reportError("Unable to parse UTF-8 rune")
	}

	file.curRead += size

	if r == '\n' { // new line
		file.curRow++
		file.curColumn = 0
	} else {
		file.curColumn++
	}

	return r
}

// Peek returns the next rune without modifying any field (does not considered as a read)
func (file *File) Peek() rune {
	if file.IsEOF() {
		panic("EOF in " + file.dirname + "/" + file.filename)
	}

	r, _ := utf8.DecodeRune(file.src[file.curRead:])

	if r == utf8.RuneError {
		file.reportError("Unable to parse UTF-8 rune")
	}

	return r
}

// IsEOF returns true if nothing can be read further
func (file *File) IsEOF() bool {
	return file.curRead >= len(file.src)
}

func (file *File) CurLoc() locator.Locator {
	var loc locator.FileLocation
	loc.Col = file.curColumn
	loc.Row = file.curRow
	loc.FileName = file.filename
	loc.DirName = file.dirname

	return &loc
}

// ReportError reports error message at current reading location
func (file *File) reportError(message string) {
	log.Fatal(message + " at row " + strconv.Itoa(file.curRow) + ", column " + strconv.Itoa(file.curColumn))
}
