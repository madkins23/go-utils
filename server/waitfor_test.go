package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	badURL  = "http://localhost:666/dead"
	goodURL = "http://localhost:8080/ping"
)

func TestIsReady_NoService(t *testing.T) {
	assert.ErrorIs(t, IsReady(badURL), syscall.ECONNREFUSED)
}

func TestIsReady_Ping(t *testing.T) {
	server := runPingService()
	require.NotNil(t, server)
	time.Sleep(100 * time.Millisecond) // wait for server to start without WaitFor()
	assert.NoError(t, IsReady(server.URL))
	server.Close()
}

func TestWaitFor_NoService(t *testing.T) {
	assert.ErrorIs(t, WaitFor(badURL, 10*time.Millisecond), ErrServerNotReady)
}

func TestWaitFor_Bad(t *testing.T) {
	server := runBadService()
	require.NotNil(t, server)
	assert.ErrorIs(t, WaitFor(badURL, 10*time.Millisecond), ErrServerNotReady)
	server.Close()
}

func TestWaitFor_Ping(t *testing.T) {
	server := runPingService()
	require.NotNil(t, server)
	assert.NoError(t, WaitFor(server.URL, 100*time.Millisecond))
	server.Close()
}

//////////////////////////////////////////////////////////////////////////

func ExampleWaitFor() {
	if err := IsReady(goodURL); err != nil {
		fmt.Printf("Not Ready 1: %s\n", err)
	}

	server := runPingService()
	defer server.Close()

	if err := WaitFor(server.URL, 100*time.Millisecond); err != nil {
		fmt.Printf("Not Ready 2: %s\n", err)
	} else {
		fmt.Println("Success!")
	}

	// Output:
	// Not Ready 1: Get "http://localhost:8080/ping": dial tcp 127.0.0.1:8080: connect: connection refused
	// Success!
}

//////////////////////////////////////////////////////////////////////////

func runBadService() *httptest.Server {
	return runService(true)
}

func runPingService() *httptest.Server {
	return runService(false)
}

func runService(bad bool) *httptest.Server {
	var status int
	if bad {
		status = http.StatusNotImplemented
	} else {
		status = http.StatusOK
	}
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte("Data"))
		}))
}

//func stopPingService(srvr *http.Server) error {
//	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	return srvr.Shutdown(ctxt)
//}
