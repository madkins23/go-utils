package cycle

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleTicker() {
	p, err := NewPeriodic(func(cycles uint) error {
		fmt.Println("Cycle Function")
		return nil
	}, func() {
		fmt.Println("Final Function")
	}, nil)
	if err != nil {
		panic(err)
	}
	go p.Ticker(10 * time.Millisecond)
	go func() {
		time.Sleep(35 * time.Millisecond)
		p.Stop()
	}()
	if err = p.Finished(); err != nil {
		panic(err)
	}
}

func TestTickerSimple(t *testing.T) {
	var signal os.Signal
	count := 0
	stopped := false
	p, err := NewPeriodic(func(cycles uint) error {
		count++
		return nil
	}, func() {
		stopped = true
	}, func(sig os.Signal) {
		signal = sig
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(11 * time.Millisecond)
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Stop()
	}()
	assert.NoError(t, p.Finished())
	assert.Contains(t, []int{10, 11}, count)
	assert.True(t, stopped)
	assert.Equal(t, nil, signal)
}

func TestTickerNoCycleFn(t *testing.T) {
	p, err := NewPeriodic(nil, nil, nil)
	require.Error(t, err)
	assert.Equal(t, errNoCycleFn, err)
	assert.Nil(t, p)
}

func TestTickerErrorInCycle0(t *testing.T) {
	started := false
	stopped := false
	p, err := NewPeriodic(func(cycles uint) error {
		started = true
		return errors.New("some sort of error")
	}, func() {
		stopped = true
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	assert.Error(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
}

func TestTickerErrorInCycle3(t *testing.T) {
	started := false
	stopped := false
	p, err := NewPeriodic(func(cycles uint) error {
		started = true
		if cycles == 3 {
			return errors.New("some sort of error")
		}
		return nil
	}, func() {
		stopped = true
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	assert.Error(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
}

func TestTickerInterrupt(t *testing.T) {
	var signal os.Signal
	started := false
	stopped := false
	p, err := NewPeriodic(func(cycles uint) error {
		started = true
		return nil
	}, func() {
		stopped = true
	}, func(sig os.Signal) {
		signal = sig
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	go p.Ticker(5 * time.Millisecond)
	go func() {
		time.Sleep(10 * time.Millisecond)
		p.signals <- syscall.SIGINT
	}()
	assert.NoError(t, p.Finished())
	assert.True(t, started)
	assert.True(t, stopped)
	assert.Equal(t, os.Interrupt, signal)
}

func TestTickerLongCycle(t *testing.T) {
	total := time.Duration(0)
	count := int64(0)
	started := false
	stopped := false
	p, err := NewPeriodic(func(cycles uint) error {
		start := time.Now()
		started = true
		time.Sleep(25 * time.Millisecond)
		count++
		total += time.Now().Sub(start)
		return nil
	}, func() {
		stopped = true
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, p)
	start := time.Now()
	go p.Ticker(5 * time.Millisecond)
	go func() {
		time.Sleep(75 * time.Millisecond)
		p.Stop()
	}()
	assert.NoError(t, p.Finished())
	final := time.Now().Sub(start)
	assert.Equal(t, int64(25), total.Milliseconds()/count)
	assert.Equal(t, int64(3), count)
	assert.True(t, final.Milliseconds() > int64(75) && final.Milliseconds() < int64(100))
	assert.True(t, started)
	assert.True(t, stopped)
}
