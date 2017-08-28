package ast

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintAst(node Node) {
	b, err := json.MarshalIndent(node, "", "    ")

	if err != nil {
		fmt.Println("error: ", err)
	}

	os.Stdout.Write(b)
	fmt.Println()
}
