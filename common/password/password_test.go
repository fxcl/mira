package password

import (
	"testing"

	"mira/common/xerrors"
)

func TestGenerateAndVerify(t *testing.T) {
	password := "plain-text-password"
	hash, err := Generate(password)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if hash == "" {
		t.Fatal("Generate() returned empty hash")
	}

	err = Verify(hash, password)
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}

	err = Verify(hash, "wrong-password")
	if err != xerrors.ErrMismatchedPassword {
		t.Errorf("Verify() expected ErrMismatchedPassword, got %v", err)
	}
}
