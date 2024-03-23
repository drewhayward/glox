package lox

import "time"

type ClockFn struct{}

func (c ClockFn) Call(runtime *RuntimeState, arguments []Value) any {
	return time.Now().Second()
}

func (c ClockFn) Arity() int { return 0 }
