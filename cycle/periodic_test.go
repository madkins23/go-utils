package cycle

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/madkins23/go-utils/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTickerSimple(t *testing.T) {
	count := 0
	stopped := false
	p, err := NewPeriodic(log.Logger(), func(cycles uint) error {
		count++
		return nil
	}, func() {
		stopped = true
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(10 * time.Millisecond)
	time.Sleep(105 * time.Millisecond)
	p.Stop()
	assert.NoError(t, p.Finished())
	assert.Equal(t, 10, count)
	assert.True(t, stopped)
}

func TestTickerNoCycleFn(t *testing.T) {
	p, err := NewPeriodic(log.Logger(), nil, nil)
	require.Error(t, err)
	require.Nil(t, p)
}

func TestTickerErrorInCycle(t *testing.T) {
	started := false
	stopped := false
	p, err := NewPeriodic(log.Logger(), func(cycles uint) error {
		started = true
		return errors.New("some sort of error")
	}, func() {
		stopped = true
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	assert.Error(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
}

func TestTickerInterrupt(t *testing.T) {
	started := false
	stopped := false
	p, err := NewPeriodic(log.Logger(), func(cycles uint) error {
		started = true
		return nil
	}, func() {
		stopped = true
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	p.signals <- syscall.SIGINT
	assert.NoError(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
}

func TestTickerStop(t *testing.T) {
	started := false
	stopped := false
	p, err := NewPeriodic(log.Logger(), func(cycles uint) error {
		started = true
		return nil
	}, func() {
		stopped = true
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(15 * time.Millisecond)
	go func() {
		time.Sleep(5 * time.Millisecond)
		p.Stop()
	}()
	assert.NoError(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
}

func TestTickerLongCycle(t *testing.T) {
	total := time.Duration(0)
	count := int64(0)
	started := false
	stopped := false
	p, err := NewPeriodic(log.Logger(), func(cycles uint) error {
		start := time.Now()
		started = true
		time.Sleep(25 * time.Millisecond)
		count++
		total += time.Now().Sub(start)
		return nil
	}, func() {
		stopped = true
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	start := time.Now()
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(75 * time.Millisecond)
	p.Stop()
	assert.NoError(t, p.Finished())
	final := time.Now().Sub(start)
	assert.Equal(t, int64(25), total.Milliseconds()/count)
	assert.Equal(t, int64(3), count)
	fmt.Println("Final time:", final.Milliseconds(), "ms")
	assert.True(t, final.Milliseconds() > int64(75) && final.Milliseconds() < int64(100))
	assert.True(t, started)
	assert.True(t, stopped)
}
