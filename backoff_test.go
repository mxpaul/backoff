package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DelayTest struct {
	FailCount int
	Want      time.Duration
}

type TestCase_ExponentialBackoff_Delay struct {
	Desc    string
	Backoff *ExponentialBackoff
	Tests   []DelayTest
}

func (tc TestCase_ExponentialBackoff_Delay) Run(t *testing.T) {
	t.Run(tc.Desc, func(t *testing.T) {
		for _, test := range tc.Tests {
			got := tc.Backoff.Delay(test.FailCount)
			assert.Equalf(t,
				test.Want,
				got,
				"[%d] %v == %v", test.FailCount, test.Want, got,
			)
		}
	})
}

func TestDefaultBackoff(t *testing.T) {
	cases := []TestCase_ExponentialBackoff_Delay{
		TestCase_ExponentialBackoff_Delay{
			Desc:    "nil backof fast exponential deplay",
			Backoff: nil,
			Tests: []DelayTest{
				{FailCount: -1, Want: 500 * time.Millisecond},
				{FailCount: 0, Want: 500 * time.Millisecond},
				{FailCount: 1, Want: 500 * time.Millisecond},
				{FailCount: 2, Want: 2 * time.Second},
				{FailCount: 3, Want: 5 * time.Second},
				{FailCount: 4, Want: 10 * time.Second},
				{FailCount: 5, Want: 30 * time.Second},
				{FailCount: 6, Want: 30 * time.Second},
				{FailCount: 7, Want: 30 * time.Second},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "Base2 Powfactor1",
			Backoff: &ExponentialBackoff{
				MinDelay:  0 * time.Millisecond,
				MaxDelay:  0 * time.Second,
				Step:      1 * time.Millisecond,
				PowFactor: 1,
				Base:      2,
			},
			Tests: []DelayTest{
				{FailCount: -1, Want: 1 * time.Millisecond},
				{FailCount: 0, Want: 1 * time.Millisecond},
				{FailCount: 1, Want: 1 * time.Millisecond},
				{FailCount: 2, Want: 2 * time.Millisecond},
				{FailCount: 3, Want: 4 * time.Millisecond},
				{FailCount: 4, Want: 8 * time.Millisecond},
				{FailCount: 5, Want: 16 * time.Millisecond},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "Base2 Powfactor1 MinDelay",
			Backoff: &ExponentialBackoff{
				MinDelay:  1*time.Millisecond + time.Nanosecond,
				MaxDelay:  0 * time.Second,
				Step:      1 * time.Millisecond,
				PowFactor: 1,
				Base:      2,
			},
			Tests: []DelayTest{
				{FailCount: -1, Want: 1*time.Millisecond + time.Nanosecond},
				{FailCount: 0, Want: 1*time.Millisecond + time.Nanosecond},
				{FailCount: 1, Want: 1*time.Millisecond + time.Nanosecond},
				{FailCount: 2, Want: 2 * time.Millisecond},
				{FailCount: 3, Want: 4 * time.Millisecond},
				{FailCount: 4, Want: 8 * time.Millisecond},
				{FailCount: 5, Want: 16 * time.Millisecond},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "Base2 Powfactor1 MaxDelay",
			Backoff: &ExponentialBackoff{
				MinDelay:  0*time.Millisecond + time.Nanosecond,
				MaxDelay:  16*time.Millisecond - time.Nanosecond,
				Step:      1 * time.Millisecond,
				PowFactor: 1,
				Base:      2,
			},
			Tests: []DelayTest{
				{FailCount: 0, Want: 1 * time.Millisecond},
				{FailCount: 1, Want: 1 * time.Millisecond},
				{FailCount: 2, Want: 2 * time.Millisecond},
				{FailCount: 3, Want: 4 * time.Millisecond},
				{FailCount: 4, Want: 8 * time.Millisecond},
				{FailCount: 5, Want: 16*time.Millisecond - time.Nanosecond},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "Base2 Powfactor2",
			Backoff: &ExponentialBackoff{
				MinDelay:  0 * time.Millisecond,
				MaxDelay:  0 * time.Second,
				Step:      1 * time.Millisecond,
				PowFactor: 2,
				Base:      2,
			},
			Tests: []DelayTest{
				{FailCount: 0, Want: 1 * time.Millisecond},
				{FailCount: 1, Want: 1 * time.Millisecond},
				{FailCount: 2, Want: 4 * time.Millisecond},
				{FailCount: 3, Want: 16 * time.Millisecond},
				{FailCount: 4, Want: 64 * time.Millisecond},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "BaseN PowfactorM",
			Backoff: &ExponentialBackoff{
				MinDelay:  0 * time.Millisecond,
				MaxDelay:  0 * time.Second,
				Step:      500 * time.Millisecond,
				PowFactor: 1.3,
				Base:      3,
			},
			Tests: []DelayTest{
				{FailCount: 0, Want: 500 * time.Millisecond},
				{FailCount: 1, Want: 500 * time.Millisecond},
				{FailCount: 2, Want: 2085583755 * time.Nanosecond},
				{FailCount: 3, Want: 8699319202 * time.Nanosecond},
				{FailCount: 4, Want: 36286317623 * time.Nanosecond},
			},
		},
		TestCase_ExponentialBackoff_Delay{
			Desc: "BaseN PowfactorM Min 1s Max 10s",
			Backoff: &ExponentialBackoff{
				MinDelay:  1 * time.Second,
				MaxDelay:  10 * time.Second,
				Step:      500 * time.Millisecond,
				PowFactor: 1.3,
				Base:      3,
			},
			Tests: []DelayTest{
				{FailCount: 0, Want: 1 * time.Second},
				{FailCount: 1, Want: 1 * time.Second},
				{FailCount: 2, Want: 2085583755 * time.Nanosecond},
				{FailCount: 3, Want: 8699319202 * time.Nanosecond},
				{FailCount: 4, Want: 10 * time.Second},
			},
		},
	}
	for _, tc := range cases {
		tc.Run(t)
	}
}

func Test_NewExponentialBackoff_Error(t *testing.T) {
	_, err := NewExponentialBackoff(10*time.Second, 1*time.Second)
	assert.Error(t, err)
}

func Test_NewExponentialBackoff_Success(t *testing.T) {
	gotBackoff, err := NewExponentialBackoff(1*time.Second, 10*time.Second)

	wantBackoff := &ExponentialBackoff{
		MinDelay:  1 * time.Second,
		MaxDelay:  10 * time.Second,
		Step:      500 * time.Millisecond,
		PowFactor: 1.3,
		Base:      3,
	}

	assert.NoError(t, err)
	assert.Equal(t, gotBackoff, wantBackoff)
}

func Test_NewExponentialBackoff_ZeroMax_NonZeroMin(t *testing.T) {
	gotBackoff, err := NewExponentialBackoff(1*time.Second, 0*time.Second)

	wantBackoff := &ExponentialBackoff{
		MinDelay:  1 * time.Second,
		MaxDelay:  0 * time.Second,
		Step:      500 * time.Millisecond,
		PowFactor: 1.3,
		Base:      3,
	}

	assert.NoError(t, err)
	assert.Equal(t, gotBackoff, wantBackoff)
}
