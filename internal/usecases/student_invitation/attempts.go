package studentinvitation

import (
	"sync"
	"time"
)

const (
	maxAttemptsPerWindow = 10
	attemptWindow        = time.Hour
)

// attemptLimiter frena la fuerza bruta sobre los códigos. Cuenta solo los
// intentos fallidos: al que acierta no hay nada que frenarle.
//
// ponytail: en memoria y por proceso, que alcanza con una sola instancia del
// backend. Si algún día corren dos, cada una lleva su propia cuenta y el techo
// real se duplica: ahí hay que mover esto a Redis o a una tabla.
type attemptLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

var limiter = &attemptLimiter{attempts: make(map[string][]time.Time)}

// allow responde si el alumno puede seguir probando, descartando de paso los
// intentos que ya salieron de la ventana.
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
		// Sin esto el mapa acumula para siempre una entrada por cada alumno
		// que alguna vez se equivocó.
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

// clear borra el historial del alumno que acertó, así el mapa no crece con
// cuentas que ya se vincularon.
func (l *attemptLimiter) clear(studentID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.attempts, studentID)
}
