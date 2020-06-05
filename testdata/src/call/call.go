package call

var Flag bool

// res is a resource which must be closed after using.
type res struct {}

func newRes() *res {
	return &res{}
}

func (r *res) close() {}

func test1() {
	r := newRes() // want `NG`
	_ = r
}

func test2() {
	r := newRes() // OK
	r.close()
}

func test3() {
	r := newRes() // want `NG`
	for i := 0; i < 3; i++ {}
	_ = r
}

func test4() {
	r := newRes() // OK
	for i := 0; i < 3; i++ {}
	r.close()
}

func test5() {
	r := newRes() // want `NG`
	if Flag {
		// this flow does not close res
		return 
	}
	r.close()
}

func test6() {
	r := newRes() // OK
	if Flag {
		r.close()
		return 
	}
	r.close()
}

func test7() {
	r := newRes() // OK
	defer r.close()
	if Flag {
		return 
	}
}

func test8() {
	r := newRes() // OK
	defer r.close()
	func() {
		func(r *res) {}(r)
	}()
}

func test9() {
	r := newRes() // NG
	func() {
		func(r *res) {}(r)
	}()
}

func test10() {
	r := newRes() // OK

	if Flag { // divide into multiple blocks
		println("hoge")
	}

	defer r.close()
	func() {
		func(r *res) {}(r)
	}()
}

func test11() {
	r := newRes() // NG

	if Flag { // divide into multiple blocks
		println("hoge")
	}

	func() {
		func(r *res) {}(r)
	}()
}

func test12() {
	r := newRes() // OK
	defer r.close()
	func() {
		_ = r // OK (free var)
	}()
}

func test13() {
	r := newRes() // OK
	defer r.close()
	func() {
		r2 := r // OK (free var)
		_ = r2
	}()
}

func test14() {
	r := newRes() // OK
	defer r.close()
	go func() {
		r2 := r // OK (free var)
		_ = r2
	}()
}

var pkgRes = newRes()
func test15() {
	r := pkgRes // OK (package variable)
	_ = r
}
