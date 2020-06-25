package withnames

var a int //lint:ignore check1 reason

//lint:ignore check1 reason
var b struct {
	N int //lint:ignore check2 reason
}

var c int // want "NG" "NG"
