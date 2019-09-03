package main

type st struct {
	o bool
}

func (*st) doSomething() {}

func (s *st) close() {}

func test1() {
	var s st
	s.doSomething() // want `close should be called after calling doSomething`
}

func test2() {
	var s st
	s.doSomething()
	s.close()
}
