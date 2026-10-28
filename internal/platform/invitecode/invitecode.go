package invitecode

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

const Length = 8

func New() (string, error) {
	max := big.NewInt(int64(len(alphabet)))

	var sb strings.Builder
	sb.Grow(Length)
	for range Length {

		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(alphabet[n.Int64()])
	}

	return sb.String(), nil
}

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

func Format(code string) string {
	if len(code) != Length {
		return code
	}

	return code[:4] + "-" + code[4:]
}
