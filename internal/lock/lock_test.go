package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAcquireRelease(t *testing.T) {
	dir := t.TempDir()

	lk, err := Acquire(dir)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	lockPath := filepath.Join(dir, lockFile)
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatal("lock file not created")
	}

	if err := lk.Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatal("lock file not removed after Release")
	}
}

func TestLockFileContents(t *testing.T) {
	dir := t.TempDir()

	lk, err := Acquire(dir)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	defer lk.Release()

	data, err := os.ReadFile(filepath.Join(dir, lockFile))
	if err != nil {
		t.Fatalf("reading lock file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, ":") {
		t.Errorf("expected PID:TIMESTAMP format, got %q", content)
	}
	if pid := readLockPID(filepath.Join(dir, lockFile)); pid == "unknown" {
		t.Error("PID should be readable from lock file")
	}
}

func TestStaleLockRemoved(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, lockFile)

	// Write a fake lock file with a PID that is not running.
	if err := os.WriteFile(lockPath, []byte("99999:0"), 0600); err != nil {
		t.Fatal(err)
	}

	// Backdate mtime so the lock looks stale (60s old).
	oldTime := time.Now().Add(-60 * time.Second)
	if err := os.Chtimes(lockPath, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	lk, err := Acquire(dir)
	if err != nil {
		t.Fatalf("expected stale lock to be cleared, got: %v", err)
	}
	defer lk.Release()
}

func TestLockContention(t *testing.T) {
	dir := t.TempDir()

	lk1, err := Acquire(dir)
	if err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}

	// Release the first lock after a short delay so the second can succeed.
	go func() {
		time.Sleep(200 * time.Millisecond)
		lk1.Release()
	}()

	lk2, err := Acquire(dir)
	if err != nil {
		t.Fatalf("second Acquire failed after contention: %v", err)
	}
	defer lk2.Release()
}

func TestLockTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 5s timeout test in short mode")
	}

	dir := t.TempDir()
	lockPath := filepath.Join(dir, lockFile)

	// Create a fresh (non-stale) lock with a recognisable fake PID.
	content := fmt.Sprintf("12345:%d", time.Now().Unix())
	if err := os.WriteFile(lockPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Acquire(dir)
	if err == nil {
		t.Fatal("expected error when lock cannot be acquired, got nil")
	}
	if !strings.Contains(err.Error(), "scaffold is already running") {
		t.Errorf("unexpected error message: %v", err)
	}
	if !strings.Contains(err.Error(), "12345") {
		t.Errorf("expected PID 12345 in error message, got: %v", err)
	}
}
