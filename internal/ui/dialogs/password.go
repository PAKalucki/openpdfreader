package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// PasswordCallback is called with the entered password.
type PasswordCallback func(password string)

// ShowPasswordDialog shows a dialog to enter a PDF password.
func ShowPasswordDialog(window fyne.Window, title string, callback PasswordCallback) {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = "Enter password"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			callback(passwordEntry.Text)
		},
	}

	dlg := dialog.NewCustomConfirm(title, "Open", "Cancel", form, func(ok bool) {
		if ok {
			callback(passwordEntry.Text)
		}
	}, window)

	dlg.Resize(fyne.NewSize(350, 150))
	dlg.Show()

	// Focus the password entry
	window.Canvas().Focus(passwordEntry)
}

// ShowSetPasswordDialog shows a dialog to set a password on a PDF.
func ShowSetPasswordDialog(window fyne.Window, callback func(userPw, ownerPw string)) {
	userPwEntry := widget.NewPasswordEntry()
	userPwEntry.PlaceHolder = "Password to open document"

	ownerPwEntry := widget.NewPasswordEntry()
	ownerPwEntry.PlaceHolder = "Password for permissions (optional)"

	confirmEntry := widget.NewPasswordEntry()
	confirmEntry.PlaceHolder = "Re-enter user password"

	form := container.NewVBox(
		widget.NewLabel("User Password (required to open):"),
		userPwEntry,
		widget.NewLabel("Confirm Password:"),
		confirmEntry,
		widget.NewSeparator(),
		widget.NewLabel("Owner Password (for editing permissions):"),
		ownerPwEntry,
	)

	dlg := dialog.NewCustomConfirm("Set PDF Password", "Set Password", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		if userPwEntry.Text == "" {
			dialog.ShowError(errorf("User password cannot be empty"), window)
			return
		}

		if userPwEntry.Text != confirmEntry.Text {
			dialog.ShowError(errorf("Passwords do not match"), window)
			return
		}

		ownerPw := ownerPwEntry.Text
		if ownerPw == "" {
			ownerPw = userPwEntry.Text // Use user password as owner password if not specified
		}

		callback(userPwEntry.Text, ownerPw)
	}, window)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()

	window.Canvas().Focus(userPwEntry)
}

type simpleError struct {
	msg string
}

func (e *simpleError) Error() string {
	return e.msg
}

func errorf(msg string) error {
	return &simpleError{msg: msg}
}
