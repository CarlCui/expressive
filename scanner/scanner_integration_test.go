package scanner

import (
	"io/ioutil"
	"testing"

	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/token"
)

func TestWithRealFiles(t *testing.T) {
	dirName := "./testFiles"

	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		panic("Incorrect test file directory!")
	}

	for _, file := range files {
		fileName := file.Name()

		var fileInput input.File
		fileInput.Init(dirName, fileName)

		var s ExpressiveScanner
		s.Init(&fileInput)

		for cur := s.Next(); cur.TokenType != token.EOF; cur = s.Next() {
			if cur.TokenType == token.ILLEGAL {
				t.Errorf("Illegal token at %v", cur.Locator.Locate())
			}
		}
	}
}
