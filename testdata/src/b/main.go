package main

type st struct {
	o bool
}

func (s *st) doSomethingAndReturnSt() *st {
	return s
}

func (s *st) close() {}

func test1() {
	var s st
	s.doSomethingAndReturnSt() // want `close should be called after calling doSomething`
}

func test2() {
	var s st
	s.doSomethingAndReturnSt().close()
}

func test3() {
	var s = &st{}
	s.doSomethingAndReturnSt()
	for i := 0; i < 3; i++ {
		// simple loop to check if the analyzer properly stops and lints in cycle.
	}
	s.close()
}

func test4() {
	var s = &st{}
	s.doSomethingAndReturnSt() // want `close should be called after calling doSomethingAndReturnSt`
	for i := 0; i < 3; i++ {
		// simple loop to check if the analyzer properly stops and lints in cycle.
	}
}
