package main

type st struct {
	o bool
}

func (s st) doSomethingAndReturnSt() st {
	return s
}

func (s st) close() {}

func test1() {
	var s st
	s.doSomethingAndReturnSt() // want `close should be called after calling doSomething`
}

func test2() {
	var s st
	s.doSomethingAndReturnSt().close()
}
