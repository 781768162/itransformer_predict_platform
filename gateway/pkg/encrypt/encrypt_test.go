package encrypt

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	pwd := "P@ssw0rd!"

	hashed, err := HashPassword(pwd)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashed == "" {
		t.Fatalf("expected non-empty hash")
	}

	// correct password should succeed
	if err := CheckPassword(hashed, pwd); err != nil {
		t.Fatalf("CheckPassword failed for correct password: %v", err)
	}

	// wrong password should fail
	if err := CheckPassword(hashed, "wrong"); err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}
}
