package routinewrapper

import (
	"sync"
)

var handle func()
var _once sync.Once

func Init(fn func()) {
	_once.Do(func() {
		// this sets the global handle function
		handle = fn
	})
}

func RoutineGenerator(fn func()) {
	defer handle()
	fn()
}
