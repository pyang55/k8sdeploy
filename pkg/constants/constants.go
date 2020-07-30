// Package constants provides one place to store constants used by other
// packages
package constants

import "time"

const (
	// WatcherTimeout is how long hte deployment watcher should wait before timeout
	WatcherTimeout = 45 * time.Minute
)
