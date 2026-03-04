package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	lockFile     = ".scaffold.lock"
	retryTimeout = 5 * time.Second
	retryInterval = 50 * time.Millisecond
	staleAge     = 30 * time.Second
)

// Lock represents an acquired file lock.
type Lock struct {
	path string
}

// Acquire creates a lock file in dir. It retries for up to 5 seconds,
// removing stale locks (older than 30s) before retrying.
func Acquire(dir string) (*Lock, error) {
	path := filepath.Join(dir, lockFile)
	deadline := time.Now().Add(retryTimeout)

	for {
		lk, err := tryAcquire(path)
		if err == nil {
			return lk, nil
		}

		// Check for stale lock.
		if info, statErr := os.Stat(path); statErr == nil {
			if time.Since(info.ModTime()) > staleAge {
				os.Remove(path)
				continue
			}
		}

		if time.Now().After(deadline) {
			pid := readLockPID(path)
			return nil, fmt.Errorf("scaffold is already running (lock held by PID %s)", pid)
		}

		time.Sleep(retryInterval)
	}
}

func tryAcquire(path string) (*Lock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d:%d", os.Getpid(), time.Now().Unix())
	return &Lock{path: path}, nil
}

// Release removes the lock file.
func (l *Lock) Release() error {
	return os.Remove(l.path)
}

func readLockPID(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) < 1 {
		return "unknown"
	}
	pid := strings.TrimSpace(parts[0])
	if _, err := strconv.Atoi(pid); err != nil {
		return "unknown"
	}
	return pid
}
