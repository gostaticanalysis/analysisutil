package withoutnames

var a int //lint:ignore test reason

//lint:ignore test reason
var b struct {
	N int //lint:ignore test reason
}

var c int // want "NG" "NG"
