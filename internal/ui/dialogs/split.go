package dialogs

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

// SplitDialog handles splitting a PDF into individual page files.
type SplitDialog struct {
	window    fyne.Window
	inputPath string
	outputDir string
	onSplit   func(outputDir string)
}

// NewSplitDialog creates a new split dialog.
func NewSplitDialog(window fyne.Window, defaultInputPath string, onSplit func(outputDir string)) *SplitDialog {
	return &SplitDialog{
		window:    window,
		inputPath: defaultInputPath,
		onSplit:   onSplit,
	}
}

// Show displays the split dialog.
func (d *SplitDialog) Show() {
	inputEntry := widget.NewEntry()
	inputEntry.SetText(d.inputPath)

	outputEntry := widget.NewEntry()
	outputEntry.SetText(d.outputDir)

	chooseInputBtn := widget.NewButton("Choose PDF...", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			d.inputPath = reader.URI().Path()
			inputEntry.SetText(d.inputPath)
		}, d.window)
	})

	chooseOutputBtn := widget.NewButton("Choose Output Folder...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}

			d.outputDir = uri.Path()
			outputEntry.SetText(d.outputDir)
		}, d.window)
	})

	content := container.NewVBox(
		widget.NewLabel("Split a PDF into one file per page."),
		widget.NewSeparator(),
		widget.NewLabel("Input PDF:"),
		inputEntry,
		chooseInputBtn,
		widget.NewSeparator(),
		widget.NewLabel("Output folder:"),
		outputEntry,
		chooseOutputBtn,
	)

	dlg := dialog.NewCustomConfirm("Split PDF", "Split", "Cancel", content, func(ok bool) {
		if !ok {
			return
		}

		if err := d.validate(); err != nil {
			dialog.ShowError(err, d.window)
			return
		}

		merger := pdf.NewMerger()
		if err := merger.Split(d.inputPath, d.outputDir); err != nil {
			dialog.ShowError(err, d.window)
			return
		}

		dialog.ShowInformation("Success", "Split completed in:\n"+d.outputDir, d.window)
		if d.onSplit != nil {
			d.onSplit(d.outputDir)
		}
	}, d.window)

	dlg.Resize(fyne.NewSize(520, 330))
	dlg.Show()
}

func (d *SplitDialog) validate() error {
	if d.inputPath == "" {
		return errors.New("select an input PDF file")
	}
	if d.outputDir == "" {
		return errors.New("select an output folder")
	}
	return nil
}
