package secretbox

import "testing"

func TestSealOpenRoundTrip(t *testing.T) {
	box, err := New("a passphrase of any length at all")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	const secret = "sk-0123456789abcdef"
	sealed, err := box.Seal(secret)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if sealed == secret {
		t.Fatal("sealed value equals the plaintext")
	}

	opened, err := box.Open(sealed)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if opened != secret {
		t.Fatalf("Open returned %q, want %q", opened, secret)
	}
}

func TestSealIsNonDeterministic(t *testing.T) {
	box, _ := New("passphrase")
	first, _ := box.Seal("sk-same-input")
	second, _ := box.Seal("sk-same-input")
	if first == second {
		t.Fatal("two seals of the same plaintext produced the same ciphertext")
	}
}

func TestOpenRejectsTamperedAndForeignCiphertext(t *testing.T) {
	box, _ := New("passphrase")
	sealed, _ := box.Seal("sk-secret")

	tampered := []byte(sealed)
	tampered[len(tampered)-2] ^= 'A' ^ 'B'
	if _, err := box.Open(string(tampered)); err == nil {
		t.Fatal("Open accepted a tampered ciphertext")
	}

	other, _ := New("a different passphrase")
	if _, err := other.Open(sealed); err == nil {
		t.Fatal("Open accepted a ciphertext sealed with another key")
	}

	if _, err := box.Open("not base64 at all !!"); err == nil {
		t.Fatal("Open accepted a non-base64 value")
	}
	if _, err := box.Open(""); err == nil {
		t.Fatal("Open accepted an empty value")
	}
}

func TestNewRejectsEmptyKey(t *testing.T) {
	if _, err := New("   "); err != ErrNoKey {
		t.Fatalf("New(blank) returned %v, want ErrNoKey", err)
	}
}

func TestLast4(t *testing.T) {
	if got := Last4("sk-0123456789abcdef"); got != "cdef" {
		t.Fatalf("Last4 = %q, want %q", got, "cdef")
	}
	if got := Last4("abcd"); got != "" {
		t.Fatalf("Last4 of a short secret = %q, want empty", got)
	}
}
