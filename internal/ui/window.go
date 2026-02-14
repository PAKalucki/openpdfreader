// Package ui provides the user interface components for OpenPDF Reader.
package ui

import (
	"errors"
	"fmt"
	"path/filepath"
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
	tabs         *container.AppTabs
	viewer       *Viewer
	toolbar      *Toolbar
	sidebar      *Sidebar
	statusBar    *widget.Label
	document     *pdf.Document
	selectedText string
	selectedPage int
	openTabs     []*DocumentTab
}

// DocumentTab represents one open PDF tab.
type DocumentTab struct {
	item         *container.TabItem
	path         string
	document     *pdf.Document
	viewer       *Viewer
	sidebar      *Sidebar
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
	// Create toolbar
	mw.toolbar = NewToolbar(mw)

	mw.tabs = container.NewAppTabs()
	mw.tabs.SetTabLocation(container.TabLocationTop)
	mw.tabs.OnSelected = func(item *container.TabItem) {
		tab := mw.findTabByItem(item)
		if tab == nil {
			return
		}
		mw.activateTab(tab)
		mw.statusBar.SetText("Active: " + tab.path)
	}

	// Main layout
	content := container.NewBorder(
		mw.toolbar.Container(), // top
		mw.statusBar,           // bottom
		nil,                    // left
		nil,                    // right
		mw.tabs,                // center
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
		fyne.NewMenuItem("List Form Fields", mw.onListFormFields),
		fyne.NewMenuItem("Fill Form Fields...", mw.onFillFormFields),
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
	tab := mw.newDocumentTab(doc, path)
	mw.openTabs = append(mw.openTabs, tab)
	mw.tabs.Append(tab.item)
	mw.tabs.Select(tab.item)
	mw.activateTab(tab)
	mw.statusBar.SetText("Loaded: " + path)

	mw.config.AddRecentFile(path)
	mw.config.Save()
}

func (mw *MainWindow) newDocumentTab(doc *pdf.Document, path string) *DocumentTab {
	viewer := NewViewer()
	viewer.SetDocument(doc)

	sidebar := NewSidebar(viewer)
	sidebar.SetDocument(doc)

	split := container.NewHSplit(
		sidebar.Container(),
		viewer.Container(),
	)
	split.SetOffset(0.2)

	title := tabTitleForPath(path)
	item := container.NewTabItem(title, split)

	return &DocumentTab{
		item:         item,
		path:         path,
		document:     doc,
		viewer:       viewer,
		sidebar:      sidebar,
		selectedText: "",
		selectedPage: -1,
	}
}

func (mw *MainWindow) findTabByItem(item *container.TabItem) *DocumentTab {
	for _, tab := range mw.openTabs {
		if tab.item == item {
			return tab
		}
	}
	return nil
}

func (mw *MainWindow) currentTab() *DocumentTab {
	if mw.tabs == nil {
		return nil
	}
	return mw.findTabByItem(mw.tabs.Selected())
}

func (mw *MainWindow) activateTab(tab *DocumentTab) {
	mw.document = tab.document
	mw.viewer = tab.viewer
	mw.sidebar = tab.sidebar
	mw.selectedText = tab.selectedText
	mw.selectedPage = tab.selectedPage
	mw.window.SetTitle("OpenPDF Reader - " + tab.path)
}

func (mw *MainWindow) syncSelectionToCurrentTab() {
	tab := mw.currentTab()
	if tab == nil {
		return
	}
	tab.selectedText = mw.selectedText
	tab.selectedPage = mw.selectedPage
}

func (mw *MainWindow) updateCurrentTabPath(path string) {
	tab := mw.currentTab()
	if tab == nil {
		return
	}
	tab.path = path
	tab.item.Text = tabTitleForPath(path)
	mw.tabs.Refresh()
	mw.window.SetTitle("OpenPDF Reader - " + path)
}

func tabTitleForPath(path string) string {
	if path == "" {
		return "Untitled"
	}
	base := filepath.Base(path)
	if base == "." || base == "/" || base == "" {
		return path
	}
	return base
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
			return
		}
		mw.updateCurrentTabPath(path)
		mw.statusBar.SetText("Saved as: " + path)
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
		mw.syncSelectionToCurrentTab()
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
		mw.syncSelectionToCurrentTab()
		dialog.ShowInformation("No Text", "No selectable text found on this page", mw.window)
		return
	}

	mw.selectedText = selectedText
	mw.selectedPage = currentPage
	mw.syncSelectionToCurrentTab()
	mw.statusBar.SetText(fmt.Sprintf("Selected all text on page %d", currentPage+1))
}

func (mw *MainWindow) onZoomIn() {
	if mw.viewer != nil {
		mw.viewer.ZoomIn()
	}
}
func (mw *MainWindow) onZoomOut() {
	if mw.viewer != nil {
		mw.viewer.ZoomOut()
	}
}
func (mw *MainWindow) onFitToPage() {
	if mw.viewer != nil {
		mw.viewer.FitToPage()
	}
}
func (mw *MainWindow) onFitToWidth() {
	if mw.viewer != nil {
		mw.viewer.FitToWidth()
	}
}
func (mw *MainWindow) onToggleThumbnails() {
	if mw.sidebar != nil {
		mw.sidebar.Toggle()
	}
}
func (mw *MainWindow) onFullscreen() { mw.window.SetFullScreen(!mw.window.FullScreen()) }

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

func (mw *MainWindow) onListFormFields() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	fields, err := pdf.NewFormManager().ListFields(mw.document.Path())
	if err != nil {
		dialog.ShowError(err, mw.window)
		return
	}
	if len(fields) == 0 {
		dialog.ShowInformation("Form Fields", "No form fields found in this document", mw.window)
		return
	}

	lines := make([]string, 0, len(fields)+1)
	lines = append(lines, "Pages | Type | Name | ID | Value")
	for _, field := range fields {
		name := field.Name
		if name == "" {
			name = "(unnamed)"
		}
		lines = append(lines, fmt.Sprintf(
			"%s | %s | %s | %s | %s",
			pageNumbersToString(field.Pages),
			field.Type,
			name,
			field.ID,
			field.Value,
		))
	}

	entry := widget.NewMultiLineEntry()
	entry.SetText(strings.Join(lines, "\n"))
	entry.Disable()

	content := container.NewBorder(
		widget.NewLabel("Use field name or ID in Fill Form Fields..."),
		nil,
		nil,
		nil,
		container.NewScroll(entry),
	)

	info := dialog.NewCustom("Form Fields", "Close", content, mw.window)
	info.Resize(fyne.NewSize(760, 460))
	info.Show()
}

func (mw *MainWindow) onFillFormFields() {
	if mw.document == nil {
		dialog.ShowInformation("No Document", "Open a PDF file first", mw.window)
		return
	}

	fields, err := pdf.NewFormManager().ListFields(mw.document.Path())
	if err != nil {
		dialog.ShowError(err, mw.window)
		return
	}
	if len(fields) == 0 {
		dialog.ShowInformation("Form Fields", "No form fields found in this document", mw.window)
		return
	}

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("fieldName=value\ncheckField=true\nlistField=a,b")

	formDialog := dialog.NewForm(
		"Fill Form Fields",
		"Apply",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Assignments", entry),
		},
		func(ok bool) {
			if !ok {
				return
			}

			assignments, parseErr := parseFieldAssignments(entry.Text)
			if parseErr != nil {
				dialog.ShowError(parseErr, mw.window)
				return
			}

			page := mw.viewer.CurrentPage()
			manager := pdf.NewFormManager()
			if err := manager.FillFields(mw.document.Path(), "", assignments); err != nil {
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
			mw.statusBar.SetText(fmt.Sprintf("Updated form fields (%d assignments)", len(assignments)))
		},
		mw.window,
	)
	formDialog.Resize(fyne.NewSize(540, 280))
	formDialog.Show()
}

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

func parseFieldAssignments(input string) (map[string]string, error) {
	assignments := map[string]string{}
	lines := strings.Split(input, "\n")

	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid assignment on line %d: expected field=value", i+1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("invalid assignment on line %d: missing field name", i+1)
		}

		assignments[key] = value
	}

	if len(assignments) == 0 {
		return nil, errors.New("no assignments provided")
	}

	return assignments, nil
}

func pageNumbersToString(pages []int) string {
	if len(pages) == 0 {
		return "-"
	}

	parts := make([]string, 0, len(pages))
	for _, page := range pages {
		parts = append(parts, fmt.Sprintf("%d", page))
	}
	return strings.Join(parts, ",")
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
