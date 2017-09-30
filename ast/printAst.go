package ast

import (
	"encoding/json"
	"fmt"
	"os"
)

func SerializeAst(node Node) []byte {
	b, err := json.MarshalIndent(node, "", "    ")

	if err != nil {
		fmt.Println("error: ", err)
	}

	return b
}

func PrintAst(node Node) {
	b := SerializeAst(node)

	os.Stdout.Write(b)
	fmt.Println()
}
