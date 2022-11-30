package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/madkins23/go-utils/msg"
)

// IsReady checks to see if the specified URL can be returned correctly.
// A nil error return signifies the server is ready.
func IsReady(url string) error {
	if resp, err := http.Get(url); err != nil {
		// Erroneous response.
		return err
	} else if err = resp.Body.Close(); err != nil {
		// Shouldn't ever happen.
		return fmt.Errorf("close response body: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	} else {
		// Server is up and running.
		return nil
	}
}

// ErrServerNotReady is returned when WaitFor times out.
const ErrServerNotReady msg.ConstError = "server not ready"

// WaitFor pings server until it is actively serving requests for the specified url.
// If the server does not properly respond within the timeout an error is returned.
// A nil error return signifies the server is ready.
func WaitFor(url string, timeout time.Duration) error {
	const loopWait = 25 * time.Millisecond
	tooLate := time.Now().Add(timeout)
	for time.Now().Before(tooLate) {
		time.Sleep(loopWait)
		if err := IsReady(url); err == nil {
			return nil
		}
	}

	return ErrServerNotReady
}
