// Package invitecode genera y normaliza los códigos que el docente dicta a sus
// alumnos. El código se lee en voz alta y se copia a mano, así que el alfabeto
// pesa más que la longitud.
package invitecode

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// Crockford base32 sin I, L, O ni U: las tres primeras se confunden con 1 y 0
// al dictarlas, y sacar la U evita que salgan palabras por casualidad.
const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// Length son 8 caracteres: 40 bits de entropía. Con el límite de intentos del
// canje, probar códigos al azar no es una vía práctica.
const Length = 8

// New devuelve un código nuevo en mayúsculas y sin separadores.
func New() (string, error) {
	max := big.NewInt(int64(len(alphabet)))

	var sb strings.Builder
	sb.Grow(Length)
	for range Length {
		// crypto/rand y no math/rand: el código es lo único que hay entre un
		// desconocido y aparecer en la lista de alumnos de un docente.
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(alphabet[n.Int64()])
	}

	return sb.String(), nil
}

// Normalize deja el código como se guarda. El alumno lo escribe como lo ve o
// como se lo dictaron: en minúscula, con el guion de la pantalla, con espacios
// pegados al pegar desde el portapapeles.
func Normalize(raw string) string {
	var sb strings.Builder
	sb.Grow(len(raw))
	for _, r := range strings.ToUpper(strings.TrimSpace(raw)) {
		if strings.ContainsRune(alphabet, r) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// Format agrega el guion del medio, solo para mostrar.
func Format(code string) string {
	if len(code) != Length {
		return code
	}

	return code[:4] + "-" + code[4:]
}
