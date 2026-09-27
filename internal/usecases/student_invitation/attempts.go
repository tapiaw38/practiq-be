package studentinvitation

import (
	"sync"
	"time"
)

const (
	maxAttemptsPerWindow = 10
	attemptWindow        = time.Hour
)

type attemptLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

var limiter = &attemptLimiter{attempts: make(map[string][]time.Time)}

func (l *attemptLimiter) allow(studentID string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := now.Add(-attemptWindow)
	recent := l.attempts[studentID][:0]
	for _, at := range l.attempts[studentID] {
		if at.After(cutoff) {
			recent = append(recent, at)
		}
	}
	if len(recent) == 0 {

		delete(l.attempts, studentID)
	} else {
		l.attempts[studentID] = recent
	}

	return len(recent) < maxAttemptsPerWindow
}

func (l *attemptLimiter) recordFailure(studentID string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.attempts[studentID] = append(l.attempts[studentID], now)
}

func (l *attemptLimiter) clear(studentID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.attempts, studentID)
}
