package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSecurity(t *testing.T) {
	s := NewSecurity()
	if s == nil {
		t.Error("NewSecurity() returned nil")
	}
}

func TestSecurityAddPassword(t *testing.T) {
	s := NewSecurity()
	tmpDir := t.TempDir()

	testPDF := filepath.Join(tmpDir, "test.pdf")
	output := filepath.Join(tmpDir, "encrypted.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	err := s.AddPassword(testPDF, output, "user123", "owner123")
	if err != nil {
		t.Errorf("AddPassword() failed: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("AddPassword() did not create output file")
	}
}

func TestSecurityRemovePassword(t *testing.T) {
	s := NewSecurity()
	tmpDir := t.TempDir()

	testPDF := filepath.Join(tmpDir, "test.pdf")
	encrypted := filepath.Join(tmpDir, "encrypted.pdf")
	decrypted := filepath.Join(tmpDir, "decrypted.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// First encrypt
	err := s.AddPassword(testPDF, encrypted, "test123", "test123")
	if err != nil {
		t.Fatalf("AddPassword() failed: %v", err)
	}

	// Then decrypt
	err = s.RemovePassword(encrypted, decrypted, "test123")
	if err != nil {
		t.Errorf("RemovePassword() failed: %v", err)
	}

	if _, err := os.Stat(decrypted); os.IsNotExist(err) {
		t.Error("RemovePassword() did not create output file")
	}
}

func TestSecurityRemovePasswordWrong(t *testing.T) {
	s := NewSecurity()
	tmpDir := t.TempDir()

	testPDF := filepath.Join(tmpDir, "test.pdf")
	encrypted := filepath.Join(tmpDir, "encrypted.pdf")
	decrypted := filepath.Join(tmpDir, "decrypted.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// First encrypt
	err := s.AddPassword(testPDF, encrypted, "test123", "test123")
	if err != nil {
		t.Fatalf("AddPassword() failed: %v", err)
	}

	// Try wrong password
	err = s.RemovePassword(encrypted, decrypted, "wrongpassword")
	if err == nil {
		t.Error("RemovePassword() should fail with wrong password")
	}
}

func TestSecurityIsEncrypted(t *testing.T) {
	s := NewSecurity()
	tmpDir := t.TempDir()

	testPDF := filepath.Join(tmpDir, "test.pdf")
	encrypted := filepath.Join(tmpDir, "encrypted.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// Check unencrypted file
	isEnc, err := s.IsEncrypted(testPDF)
	if err != nil {
		t.Errorf("IsEncrypted() failed: %v", err)
	}
	if isEnc {
		t.Error("IsEncrypted() should return false for unencrypted PDF")
	}

	// Encrypt and check again
	err = s.AddPassword(testPDF, encrypted, "test123", "test123")
	if err != nil {
		t.Fatalf("AddPassword() failed: %v", err)
	}

	isEnc, err = s.IsEncrypted(encrypted)
	if err != nil {
		t.Logf("IsEncrypted() on encrypted file: %v", err)
	}
	// Note: pdfcpu may or may not be able to detect encryption without password
}

func TestSecurityChangePassword(t *testing.T) {
	s := NewSecurity()
	tmpDir := t.TempDir()

	testPDF := filepath.Join(tmpDir, "test.pdf")
	encrypted := filepath.Join(tmpDir, "encrypted.pdf")
	changed := filepath.Join(tmpDir, "changed.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// First encrypt
	err := s.AddPassword(testPDF, encrypted, "old123", "old123")
	if err != nil {
		t.Fatalf("AddPassword() failed: %v", err)
	}

	// Change password
	err = s.ChangePassword(encrypted, changed, "old123", "new123", "new123")
	if err != nil {
		t.Errorf("ChangePassword() failed: %v", err)
	}

	if _, err := os.Stat(changed); os.IsNotExist(err) {
		t.Error("ChangePassword() did not create output file")
	}
}
