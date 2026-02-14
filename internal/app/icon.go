package app

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/app-icon.svg
var appIconSVG []byte

func appIconResource() fyne.Resource {
	return fyne.NewStaticResource("app-icon.svg", appIconSVG)
}
