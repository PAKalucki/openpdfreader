# TODO

## In Progress
- [ ] Install Fyne dependencies (gcc, libgl1-mesa-dev, xorg-dev)
- [ ] Run `go mod tidy` to resolve dependencies
- [ ] Implement actual PDF rendering using go-pdfium

## Backlog
- [ ] Add file open via command-line argument
- [ ] Implement text selection and copy
- [ ] Add thumbnail rendering in sidebar
- [ ] Implement zoom fit to page/width calculations
- [ ] Add print dialog
- [ ] Implement merge PDFs dialog
- [ ] Implement split PDF dialog
- [ ] Add annotation tools
- [ ] Implement form field detection and filling
- [ ] Add signature pad
- [ ] Implement password protection dialogs
- [ ] Add redaction tools
- [ ] Implement PDF to image export
- [ ] Add keyboard shortcuts
- [ ] Implement undo/redo system
- [ ] Add dark/light theme support
- [ ] Create application icon
- [ ] Set up CI/CD pipeline
- [ ] Create installers for Windows/Linux

## Completed
- [x] Create initial project plan in project.md (commit: 5609210)
- [x] Update project plan for Go + Fyne stack (commit: 5609210)
- [x] Set up Go project structure (commit: 5609210)
- [x] Create main application entry point (commit: 5609210)
- [x] Create main window with menus and toolbar (commit: 5609210)
- [x] Create PDF viewer widget (placeholder rendering) (commit: 5609210)
- [x] Create sidebar for page list (commit: 5609210)
- [x] Implement basic PDF document loading with pdfcpu (commit: 5609210)
- [x] Add merge/split/extract page operations (commit: 5609210)
- [x] Create Makefile for build automation (commit: 5609210)
- [x] Create README.md (commit: 5609210)