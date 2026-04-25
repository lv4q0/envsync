package envfile

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched .env file.
type WatchEvent struct {
	Path    string
	Changed bool
	Err     error
}

// FileHash returns the MD5 hash of a file's contents.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls the given file at the specified interval and sends a WatchEvent
// on the returned channel whenever the file changes. It stops when done is closed.
func Watch(path string, interval time.Duration, done <-chan struct{}) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)

	go func() {
		defer close(ch)
		lastHash, err := FileHash(path)
		if err != nil {
			ch <- WatchEvent{Path: path, Err: err}
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				current, err := FileHash(path)
				if err != nil {
					ch <- WatchEvent{Path: path, Err: err}
					return
				}
				if current != lastHash {
					lastHash = current
					ch <- WatchEvent{Path: path, Changed: true}
				}
			}
		}
	}()

	return ch
}
