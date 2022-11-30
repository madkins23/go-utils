package server

import (
	"fmt"
	"os"
	"syscall"
)

func Interrupt() error {
	if proc, err := os.FindProcess(os.Getpid()); err != nil {
		return fmt.Errorf("get current process: %w", err)
	} else if err = proc.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("send SIGINT to current process: %w", err)
	} else {
		return nil
	}
}
