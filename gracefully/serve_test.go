package gracefully_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// We initially compile the dummy servers we have, we run them
// as exec commands under this process, send a request to
// them and then send them an interrupt signal.
//
// The dummy servers artificially block the request for 5 seconds
// before responding with a static response.
//
// The intended behaviour is for the request to complete even
// after the interrupt signal sends.
//

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	cleanup()
	os.Exit(c)
}

func TestGracefulShutdown(t *testing.T) {
	t.Run("captures SIGINT", func(t *testing.T) { assertServerGracefullyShutsDown(t, simpleServerBin, syscall.SIGINT) })
	t.Run("captures SIGTERM", func(t *testing.T) { assertServerGracefullyShutsDown(t, simpleServerBin, syscall.SIGTERM) })
}

func TestAfterShutdown(t *testing.T) {
	t.Run("captures SIGINT", func(t *testing.T) { assertAfterShutdownServerAfterStop(t, syscall.SIGINT) })
	t.Run("captures SIGTERM", func(t *testing.T) { assertAfterShutdownServerAfterStop(t, syscall.SIGTERM) })
}

func assertServerGracefullyShutsDown(t *testing.T, bin string, s syscall.Signal) {
	done := make(chan struct{}, 1)
	defer close(done)
	cmd := exec.Command(bin)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	must(cmd.Start())
	time.Sleep(500 * time.Millisecond) // Ensure server is up before starting with anything
	go func() {
		r, err := http.Get("http://localhost:4000/")
		assert.Empty(t, err)
		assert.Equal(t, r.StatusCode, http.StatusOK)
		done <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // Ensure request did start before
	must(cmd.Process.Signal(s))
	<-done
}

func assertAfterShutdownServerAfterStop(t *testing.T, s syscall.Signal) {
	assertServerGracefullyShutsDown(t, afterShutdownServerBin, s)
	time.Sleep(500 * time.Millisecond)      // Give the after shutdown hooks time to run
	assert.FileExists(t, afterShutdownFile) // If this file exists then we know the after shutdown hooks worked on that server
	err := os.Remove(afterShutdownFile)     // This is to ensure the following tests are required to recreate the file
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove the file after asserting for it's existence got: %s", err)
	}
}
