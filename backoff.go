package backoff

import (
	"fmt"
	"math"
	"time"
)

type Backoff interface {
	Delay(int) time.Duration
}

type ExponentialBackoff struct {
	MinDelay  time.Duration
	MaxDelay  time.Duration
	Step      time.Duration
	PowFactor float64
	Base      float64
}

func FallbackDelay(attempt int) time.Duration {
	switch {
	case attempt <= 0:
		return 500 * time.Millisecond
	case attempt == 1:
		return 2 * time.Second
	case attempt == 2:
		return 5 * time.Second
	case attempt == 3:
		return 10 * time.Second
	}
	return 30 * time.Second
}

func (b *ExponentialBackoff) Delay(failCount int) time.Duration {
	attempt := failCount - 1
	if failCount <= 0 {
		attempt = 0
	}

	if b == nil {
		return FallbackDelay(attempt)
	}

	delay := time.Duration(float64(b.Step) * math.Pow(b.Base, b.PowFactor*float64(attempt)))

	if delay < b.MinDelay {
		delay = b.MinDelay
	} else if b.MaxDelay > 0 && delay > b.MaxDelay {
		delay = b.MaxDelay
	}

	return delay
}

func NewExponentialBackoff(min, max time.Duration) (*ExponentialBackoff, error) {
	if max > 0 && max < min {
		return nil, fmt.Errorf("fuck")
	}

	b := &ExponentialBackoff{
		MinDelay:  min,
		MaxDelay:  max,
		Step:      500 * time.Millisecond,
		PowFactor: 1.3,
		Base:      3,
	}

	return b, nil
}
