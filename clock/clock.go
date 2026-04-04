package clock

import "time"

var Default = &DefaultClock{}

type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type DefaultClock struct{}

func (c *DefaultClock) Now() time.Time {
	return time.Now()
}

func (c *DefaultClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
