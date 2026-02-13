package pdf

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Security handles PDF encryption and password operations.
type Security struct{}

// NewSecurity creates a new Security instance.
func NewSecurity() *Security {
	return &Security{}
}

// AddPassword encrypts a PDF with the given passwords.
// userPw is required to open the document.
// ownerPw (if different) allows editing without restrictions.
func (s *Security) AddPassword(inputPath, outputPath, userPw, ownerPw string) error {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = userPw
	conf.OwnerPW = ownerPw

	return api.EncryptFile(inputPath, outputPath, conf)
}

// RemovePassword decrypts a PDF and removes password protection.
// The correct password must be provided.
func (s *Security) RemovePassword(inputPath, outputPath, password string) error {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password

	return api.DecryptFile(inputPath, outputPath, conf)
}

// ChangePassword changes the passwords on an encrypted PDF.
// Both user and owner passwords can be changed.
func (s *Security) ChangePassword(inputPath, outputPath, oldPw, newUserPw, newOwnerPw string) error {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = oldPw
	conf.OwnerPW = oldPw

	// Change user password first
	err := api.ChangeUserPasswordFile(inputPath, outputPath, oldPw, newUserPw, conf)
	if err != nil {
		return err
	}

	// If owner password is different, change it too
	if newOwnerPw != "" && newOwnerPw != newUserPw {
		conf.UserPW = newUserPw
		conf.OwnerPW = oldPw
		return api.ChangeOwnerPasswordFile(outputPath, outputPath, oldPw, newOwnerPw, conf)
	}

	return nil
}

// IsEncrypted checks if a PDF is password protected.
func (s *Security) IsEncrypted(path string) (bool, error) {
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return false, err
	}
	return ctx.Encrypt != nil, nil
}
