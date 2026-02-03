package progress

import "github.com/gosuri/uiprogress"

type Bar struct {
	bar *uiprogress.Bar
}

func NewProgressBar(size int) *Bar {
	uiprogress.Start()

	bar := uiprogress.AddBar(size)
	bar.AppendCompleted()
	bar.PrependElapsed()

	return &Bar{bar: bar}
}

func (pb *Bar) Increase() {
	pb.bar.Incr()
}
