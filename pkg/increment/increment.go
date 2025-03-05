package increment

import "sync"

type Increment struct {
	increment int
	mutex     *sync.Mutex
}

func NewIncrement() *Increment {
	return &Increment{
		increment: 0,
		mutex:     &sync.Mutex{},
	}
}

func (i *Increment) Add() {
	i.mutex.Lock()
	i.increment++
	i.mutex.Unlock()
}

func (i Increment) Get() int {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.increment
}

var GlobalIncrement Increment

func init() {
	GlobalIncrement = *NewIncrement()
}
