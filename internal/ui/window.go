// Package ui provides the user interface components for OpenPDF Reader.
package ui

import (
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
	window    fyne.Window
	config    *config.Config
	viewer    *Viewer
	toolbar   *Toolbar
	sidebar   *Sidebar
	statusBar *widget.Label
	document  *pdf.Document
}

// NewMainWindow creates a new main window.
func NewMainWindow(fyneApp fyne.App, cfg *config.Config) *MainWindow {
	window := fyneApp.NewWindow("OpenPDF Reader")
	window.Resize(fyne.NewSize(float32(cfg.WindowWidth), float32(cfg.WindowHeight)))

	mw := &MainWindow{
		window:    window,
		config:    cfg,
		statusBar: widget.NewLabel("Ready"),
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
		mw.statusBar,          // bottom
		nil,                   // left
		nil,                   // right
		split,                 // center
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

func (mw *MainWindow) onPrint()             { /* TODO: Implement print dialog */ }
func (mw *MainWindow) onUndo()              { /* TODO: Implement undo */ }
func (mw *MainWindow) onRedo()              { /* TODO: Implement redo */ }
func (mw *MainWindow) onCopy()              { /* TODO: Implement copy */ }
func (mw *MainWindow) onSelectAll()         { /* TODO: Implement select all */ }
func (mw *MainWindow) onZoomIn()            { mw.viewer.ZoomIn() }
func (mw *MainWindow) onZoomOut()           { mw.viewer.ZoomOut() }
func (mw *MainWindow) onFitToPage()         { mw.viewer.FitToPage() }
func (mw *MainWindow) onFitToWidth()        { mw.viewer.FitToWidth() }
func (mw *MainWindow) onToggleThumbnails()  { mw.sidebar.Toggle() }
func (mw *MainWindow) onFullscreen()        { mw.window.SetFullScreen(!mw.window.FullScreen()) }

func (mw *MainWindow) onMergePDFs() {
	dlg := dialogs.NewMergeDialog(mw.window, func(files []string, output string) {
		// Optionally open the merged file
		mw.OpenFile(output)
	})
	dlg.Show()
}

func (mw *MainWindow) onSplitPDF()          { /* TODO: Show split dialog */ }
func (mw *MainWindow) onExtractPages()      { /* TODO: Show extract dialog */ }
func (mw *MainWindow) onDeletePages()       { /* TODO: Show delete dialog */ }
func (mw *MainWindow) onRotatePages()       { /* TODO: Show rotate dialog */ }

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
