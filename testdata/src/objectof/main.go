package objectof

import (
	"fmt"
	"io"
	"vendored"
)

type A int
var EOF = io.EOF
var _ = vendored.EOF

func main() {
	fmt.Println(EOF)
}
