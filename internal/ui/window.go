// Package ui provides the user interface components for OpenPDF Reader.
package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/config"
	"github.com/openpdfreader/openpdfreader/internal/pdf"
	"github.com/openpdfreader/openpdfreader/internal/ui/dialogs"
)

// MainWindow represents the main application window.
type MainWindow struct {
	window       fyne.Window
	config       *config.Config
	viewer       *Viewer
	toolbar      *Toolbar
	sidebar      *Sidebar
	statusBar    *widget.Label
	document     *pdf.Document
	selectedText string
	selectedPage int
}

// NewMainWindow creates a new main window.
func NewMainWindow(fyneApp fyne.App, cfg *config.Config) *MainWindow {
	window := fyneApp.NewWindow("OpenPDF Reader")
	window.Resize(fyne.NewSize(float32(cfg.WindowWidth), float32(cfg.WindowHeight)))

	mw := &MainWindow{
		window:       window,
		config:       cfg,
		statusBar:    widget.NewLabel("Ready"),
		selectedPage: -1,
	}

	mw.setupUI()
	mw.setupMenus()
	mw.setupShortcuts()

	return mw
}

// Show displays the main window.
func (mw *MainWindow) Show() {
	mw.window.CenterOnScreen()
	mw.window.Show()
}

// ShowAndRun displays the window and runs the event loop.
func (mw *MainWindow) ShowAndRun() {
	mw.window.CenterOnScreen()
	mw.window.ShowAndRun()
}

func (mw *MainWindow) setupUI() {
	// Create viewer
	mw.viewer = NewViewer()

	// Create toolbar
	mw.toolbar = NewToolbar(mw)

	// Create sidebar for thumbnails
	mw.sidebar = NewSidebar(mw)

	// Main content area with sidebar and viewer
	split := container.NewHSplit(
		mw.sidebar.Container(),
		mw.viewer.Container(),
	)
	split.SetOffset(0.2) // 20% for sidebar

	// Main layout
	content := container.NewBorder(
		mw.toolbar.Container(), // top
		mw.statusBar,           // bottom
		nil,                    // left
		nil,                    // right
		split,                  // center
	)

	mw.window.SetContent(content)
}

func (mw *MainWindow) setupMenus() {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open...", mw.onOpenFile),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Save", mw.onSave),
		fyne.NewMenuItem("Save As...", mw.onSaveAs),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Print...", mw.onPrint),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exit", func() { mw.window.Close() }),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Undo", mw.onUndo),
		fyne.NewMenuItem("Redo", mw.onRedo),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Copy", mw.onCopy),
		fyne.NewMenuItem("Select All", mw.onSelectAll),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Zoom In", mw.onZoomIn),
		fyne.NewMenuItem("Zoom Out", mw.onZoomOut),
		fyne.NewMenuItem("Fit to Page", mw.onFitToPage),
		fyne.NewMenuItem("Fit to Width", mw.onFitToWidth),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Toggle Thumbnails", mw.onToggleThumbnails),
		fyne.NewMenuItem("Fullscreen", mw.onFullscreen),
	)

	toolsMenu := fyne.NewMenu("Tools",
		fyne.NewMenuItem("Merge PDFs...", mw.onMergePDFs),
		fyne.NewMenuItem("Split PDF...", mw.onSplitPDF),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Extract Pages...", mw.onExtractPages),
		fyne.NewMenuItem("Delete Pages...", mw.onDeletePages),
		fyne.NewMenuItem("Rotate Pages...", mw.onRotatePages),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Add Highlight...", mw.onAddHighlightAnnotation),
		fyne.NewMenuItem("Add Text Annotation...", mw.onAddTextAnnotation),
		fyne.NewMenuItem("Add Shape Annotation...", mw.onAddShapeAnnotation),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Add Password...", mw.onAddPassword),
		fyne.NewMenuItem("Remove Password...", mw.onRemovePassword),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", mw.onAbout),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, editMenu, viewMenu, toolsMenu, helpMenu)
	mw.window.SetMainMenu(mainMenu)
}

func (mw *MainWindow) setupShortcuts() {
	canvas := mw.window.Canvas()

	// File shortcuts
	canvas.AddShortcut(&fyne.ShortcutCopy{}, func(_ fyne.Shortcut) {
		mw.onCopy()
	})

	// Custom shortcuts using desktop package
	mw.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		// Handle key events
		switch ev.Name {
		case fyne.KeyPageUp:
			if mw.viewer != nil {
				mw.viewer.GoToPage(mw.viewer.currentPage - 1)
			}
		case fyne.KeyPageDown:
			if mw.viewer != nil {
				mw.viewer.GoToPage(mw.viewer.currentPage + 1)
			}
		case fyne.KeyHome:
			if mw.viewer != nil {
				mw.viewer.GoToPage(0)
			}
		case fyne.KeyEnd:
			if mw.viewer != nil && mw.document != nil {
				mw.viewer.GoToPage(mw.document.PageCount() - 1)
			}
		case fyne.KeyF11:
			mw.onFullscreen()
		case fyne.KeyEscape:
			if mw.window.FullScreen() {
				mw.window.SetFullScreen(false)
			}
		case fyne.KeyPlus, fyne.KeyEqual:
			mw.onZoomIn()
		case fyne.KeyMinus:
			mw.onZoomOut()
		case fyne.Key0:
			mw.onFitToPage()
		}
	})
}

// OpenFile opens a PDF file.
func (mw *MainWindow) OpenFile(path string) error {
	doc, err := pdf.Open(path)
	if err != nil {
		// Check if it's a password-protected PDF
		if pdf.IsPasswordError(err) {
			mw.openPasswordProtectedFile(path)
			return nil
		}
		return err
	}

	mw.setDocument(doc, path)
	return nil
}

func (mw *MainWindow) openPasswordProtectedFile(path string) {
	dialogs.ShowPasswordDialog(mw.window, "Password Protected PDF", func(password string) {
		doc, err := pdf.OpenWithPassword(path, password)
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		mw.setDocument(doc, path)
	})
}

func (mw *MainWindow) setDocument(doc *pdf.Document, path string) {
	mw.document = doc
	mw.selectedText = ""
	mw.selectedPage = -1
	mw.viewer.SetDocument(doc)
	mw.sidebar.SetDocument(doc)
	mw.window.SetTitle("OpenPDF Reader - " + path)
	mw.statusBar.SetText("Loaded: " + path)

	mw.config.AddRecentFile(path)
	mw.config.Save()
}

// Menu action handlers

func (mw *MainWindow) onOpenFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		reader.Close()

		path := reader.URI().Path()
		if err := mw.OpenFile(path); err != nil {
			dialog.ShowError(err, mw.window)
		}
	}, mw.window)
}

func (mw *MainWindow) onSave() {
	if mw.document == nil {
		return
	}
	if err := mw.document.Save(); err != nil {
		dialog.ShowError(err, mw.window)
	}
}

func (mw *MainWindow) onSaveAs() {
	if mw.document == nil {
		return
	}
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		writer.Close()

		path := writer.URI().Path()
		if err := mw.document.SaveAs(path); err != nil {
			dialog.ShowError(err, mw.window)
		}
	}, mw.window)
}

func (mw *MainWindow) onPrint() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	dialog.ShowConfirm("Print", "Send this document to the default printer?", func(confirmed bool) {
		if !confirmed {
			return
		}

		if err := pdf.PrintFile(mw.document.Path()); err != nil {
			dialog.ShowError(err, mw.window)
			return
		}

		mw.statusBar.SetText("Sent to printer: " + mw.document.Path())
	}, mw.window)
}

func (mw *MainWindow) onUndo() { /* TODO: Implement undo */ }
func (mw *MainWindow) onRedo() { /* TODO: Implement redo */ }
func (mw *MainWindow) onCopy() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	currentPage := mw.viewer.CurrentPage()
	selectedText := strings.TrimSpace(mw.selectedText)
	if mw.selectedPage != currentPage || selectedText == "" {
		text, err := mw.document.ExtractText(currentPage)
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		selectedText = strings.TrimSpace(text)
		mw.selectedText = selectedText
		mw.selectedPage = currentPage
	}

	if selectedText == "" {
		dialog.ShowInformation("No Text", "No selectable text found on this page", mw.window)
		return
	}

	mw.window.Clipboard().SetContent(selectedText)
	mw.statusBar.SetText(fmt.Sprintf("Copied text from page %d", currentPage+1))
}

func (mw *MainWindow) onSelectAll() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	currentPage := mw.viewer.CurrentPage()
	text, err := mw.document.ExtractText(currentPage)
	if err != nil {
		dialog.ShowError(err, mw.window)
		return
	}

	selectedText := strings.TrimSpace(text)
	if selectedText == "" {
		mw.selectedText = ""
		mw.selectedPage = -1
		dialog.ShowInformation("No Text", "No selectable text found on this page", mw.window)
		return
	}

	mw.selectedText = selectedText
	mw.selectedPage = currentPage
	mw.statusBar.SetText(fmt.Sprintf("Selected all text on page %d", currentPage+1))
}

func (mw *MainWindow) onZoomIn()           { mw.viewer.ZoomIn() }
func (mw *MainWindow) onZoomOut()          { mw.viewer.ZoomOut() }
func (mw *MainWindow) onFitToPage()        { mw.viewer.FitToPage() }
func (mw *MainWindow) onFitToWidth()       { mw.viewer.FitToWidth() }
func (mw *MainWindow) onToggleThumbnails() { mw.sidebar.Toggle() }
func (mw *MainWindow) onFullscreen()       { mw.window.SetFullScreen(!mw.window.FullScreen()) }

func (mw *MainWindow) onMergePDFs() {
	dlg := dialogs.NewMergeDialog(mw.window, func(files []string, output string) {
		// Optionally open the merged file
		mw.OpenFile(output)
	})
	dlg.Show()
}

func (mw *MainWindow) onSplitPDF() {
	defaultInputPath := ""
	if mw.document != nil {
		defaultInputPath = mw.document.Path()
	}

	dlg := dialogs.NewSplitDialog(mw.window, defaultInputPath, func(outputDir string) {
		mw.statusBar.SetText("Split complete: " + outputDir)
	})
	dlg.Show()
}
func (mw *MainWindow) onExtractPages() { /* TODO: Show extract dialog */ }
func (mw *MainWindow) onDeletePages()  { /* TODO: Show delete dialog */ }
func (mw *MainWindow) onRotatePages()  { /* TODO: Show rotate dialog */ }

func (mw *MainWindow) onAddHighlightAnnotation() {
	mw.promptAnnotationContents("Add Highlight", "Highlight content", "Highlight", func(contents string) error {
		return pdf.NewAnnotator().AddHighlight(mw.document.Path(), "", mw.viewer.CurrentPage(), contents)
	})
}

func (mw *MainWindow) onAddTextAnnotation() {
	mw.promptAnnotationContents("Add Text Annotation", "Text note", "Note", func(contents string) error {
		return pdf.NewAnnotator().AddText(mw.document.Path(), "", mw.viewer.CurrentPage(), contents)
	})
}

func (mw *MainWindow) onAddShapeAnnotation() {
	mw.promptAnnotationContents("Add Shape Annotation", "Shape label", "Shape", func(contents string) error {
		return pdf.NewAnnotator().AddShape(mw.document.Path(), "", mw.viewer.CurrentPage(), contents)
	})
}

func (mw *MainWindow) promptAnnotationContents(
	title string,
	fieldLabel string,
	defaultValue string,
	apply func(contents string) error,
) {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	entry := widget.NewMultiLineEntry()
	entry.SetText(defaultValue)

	form := dialog.NewForm(
		title,
		"Apply",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem(fieldLabel, entry),
		},
		func(ok bool) {
			if !ok {
				return
			}

			contents := strings.TrimSpace(entry.Text)
			if contents == "" {
				contents = defaultValue
			}

			page := mw.viewer.CurrentPage()
			if err := apply(contents); err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			if err := mw.document.Reload(); err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			mw.viewer.SetDocument(mw.document)
			mw.sidebar.SetDocument(mw.document)
			mw.viewer.GoToPage(page)
			mw.statusBar.SetText(fmt.Sprintf("%s added on page %d", title, page+1))
		},
		mw.window,
	)
	form.Resize(fyne.NewSize(460, 220))
	form.Show()
}

func (mw *MainWindow) onAddPassword() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	dialogs.ShowSetPasswordDialog(mw.window, func(userPw, ownerPw string) {
		// Save to a new file
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			writer.Close()

			outputPath := writer.URI().Path()
			security := pdf.NewSecurity()
			err = security.AddPassword(mw.document.Path(), outputPath, userPw, ownerPw)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			dialog.ShowInformation("Success", "Password added to:\n"+outputPath, mw.window)
		}, mw.window)
	})
}

func (mw *MainWindow) onRemovePassword() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	dialogs.ShowPasswordDialog(mw.window, "Enter Current Password", func(password string) {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			writer.Close()

			outputPath := writer.URI().Path()
			security := pdf.NewSecurity()
			err = security.RemovePassword(mw.document.Path(), outputPath, password)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			dialog.ShowInformation("Success", "Password removed. Saved to:\n"+outputPath, mw.window)
		}, mw.window)
	})
}

func (mw *MainWindow) onAbout() {
	dialog.ShowInformation("About OpenPDF Reader",
		"OpenPDF Reader v0.1.0\n\nAn open-source PDF viewer and editor.\n\nLicensed under Apache 2.0",
		mw.window)
}
