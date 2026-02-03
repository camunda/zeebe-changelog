package progress

import (
	"github.com/gosuri/uiprogress"
	"sync"
)

type Bar struct {
	bar   *uiprogress.Bar
	mutex sync.Mutex
}

func NewProgressBar(size int) *Bar {
	uiprogress.Start()

	bar := uiprogress.AddBar(size)
	bar.AppendCompleted()
	bar.PrependElapsed()

	return &Bar{bar: bar}
}

func (pb *Bar) Increase() {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.bar.Incr()
}
