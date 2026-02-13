package dialogs

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

// MergeDialog handles the PDF merge functionality.
type MergeDialog struct {
	window        fyne.Window
	files         []string
	fileList      *widget.List
	selectedIndex int
	onMerge       func(files []string, output string)
}

// NewMergeDialog creates a new merge dialog.
func NewMergeDialog(window fyne.Window, onMerge func(files []string, output string)) *MergeDialog {
	return &MergeDialog{
		window:        window,
		files:         []string{},
		selectedIndex: -1,
		onMerge:       onMerge,
	}
}

// Show displays the merge dialog.
func (d *MergeDialog) Show() {
	// File list
	d.fileList = widget.NewList(
		func() int { return len(d.files) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("filename.pdf"),
				widget.NewButton("Remove", func() {}),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)

			label.SetText(filepath.Base(d.files[id]))
			btn.OnTapped = func() {
				d.removeFile(id)
			}
		},
	)

	// Track selection
	d.fileList.OnSelected = func(id widget.ListItemID) {
		d.selectedIndex = id
	}
	d.fileList.OnUnselected = func(id widget.ListItemID) {
		d.selectedIndex = -1
	}

	// Buttons
	addBtn := widget.NewButton("Add PDF...", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			reader.Close()
			d.addFile(reader.URI().Path())
		}, d.window)
	})

	moveUpBtn := widget.NewButton("Move Up", func() {
		if d.selectedIndex > 0 {
			d.moveUp(d.selectedIndex)
			d.selectedIndex--
			d.fileList.Select(d.selectedIndex)
		}
	})

	moveDownBtn := widget.NewButton("Move Down", func() {
		if d.selectedIndex >= 0 && d.selectedIndex < len(d.files)-1 {
			d.moveDown(d.selectedIndex)
			d.selectedIndex++
			d.fileList.Select(d.selectedIndex)
		}
	})

	mergeBtn := widget.NewButton("Merge", func() {
		if len(d.files) < 2 {
			dialog.ShowInformation("Error", "Add at least 2 PDF files to merge", d.window)
			return
		}
		d.showSaveDialog()
	})
	mergeBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancel", func() {})

	// Layout
	buttonBar := container.NewHBox(addBtn, moveUpBtn, moveDownBtn)
	actionBar := container.NewHBox(cancelBtn, mergeBtn)

	content := container.NewBorder(
		widget.NewLabel("Select PDF files to merge (drag to reorder):"),
		container.NewBorder(nil, nil, nil, actionBar, buttonBar),
		nil,
		nil,
		d.fileList,
	)

	dlg := dialog.NewCustom("Merge PDFs", "Close", content, d.window)
	dlg.Resize(fyne.NewSize(500, 400))
	dlg.Show()
}

func (d *MergeDialog) addFile(path string) {
	d.files = append(d.files, path)
	d.fileList.Refresh()
}

func (d *MergeDialog) removeFile(index int) {
	if index >= 0 && index < len(d.files) {
		d.files = append(d.files[:index], d.files[index+1:]...)
		d.fileList.Refresh()
	}
}

func (d *MergeDialog) moveUp(index int) {
	if index > 0 {
		d.files[index], d.files[index-1] = d.files[index-1], d.files[index]
		d.fileList.Refresh()
	}
}

func (d *MergeDialog) moveDown(index int) {
	if index < len(d.files)-1 {
		d.files[index], d.files[index+1] = d.files[index+1], d.files[index]
		d.fileList.Refresh()
	}
}

func (d *MergeDialog) showSaveDialog() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		writer.Close()

		outputPath := writer.URI().Path()

		// Perform merge
		merger := pdf.NewMerger()
		err = merger.Merge(d.files, outputPath)
		if err != nil {
			dialog.ShowError(err, d.window)
			return
		}

		dialog.ShowInformation("Success",
			"PDFs merged successfully to:\n"+outputPath,
			d.window)

		if d.onMerge != nil {
			d.onMerge(d.files, outputPath)
		}
	}, d.window)
}
