# OpenPDF Reader - Project Plan

## Project Overview

**OpenPDF Reader** is an open-source, cross-platform PDF application that replicates the core functionality of Adobe Acrobat Reader. The application will support Windows and Linux, using open-source dependencies wherever possible.

---

## Core Features

### 1. PDF Viewing
- Open and display PDF documents
- Multi-page navigation (scroll, page jump, thumbnails)
- Zoom controls (fit to page, fit to width, percentage-based)
- Fullscreen and presentation mode
- Search within documents
- Bookmarks and outline navigation
- Text selection and copy
- Rotate pages (view only)

### 2. PDF Editing
- Add text annotations
- Highlight, underline, strikethrough text
- Add sticky notes/comments
- Draw shapes and freehand
- Add images/stamps
- Edit existing text (where supported by PDF structure)

### 3. Fill in and Sign
- Detect and fill form fields (text, checkbox, radio, dropdown)
- Digital signature support
  - Draw signature with mouse/touchpad
  - Upload signature image
  - Certificate-based digital signatures (X.509)
- Save filled forms
- Flatten form fields

### 4. Printing
- Standard print dialog integration
- Print range selection (all, current, custom range)
- Print quality options
- Print to PDF (virtual printer)
- Double-sided printing support
- Multiple pages per sheet

### 5. Conversion
- PDF to Image (PNG, JPEG, TIFF)
- PDF to Text (plain text extraction)
- PDF to Word (DOCX) - best effort
- PDF to HTML
- Image to PDF
- Merge multiple images into PDF
- Office documents to PDF (via LibreOffice integration)

### 6. Page Management
- Delete pages
- Reorder pages (drag and drop)
- Extract pages to new PDF
- Insert pages from another PDF
- Rotate pages (permanently)
- Split PDF into multiple files

### 7. Merge PDFs
- Combine multiple PDFs into one
- Reorder documents before merge
- Insert at specific position
- Batch merge capability

### 8. Redaction
- Select areas/text for redaction
- Apply redaction (permanent removal)
- Redaction preview
- Search and redact (find all instances)
- Redaction audit trail

### 9. Password Protection Support
- Open password-protected PDFs
- Enter owner/user passwords
- Remove password protection (with valid credentials)
- Add password protection
- Set document permissions

---

## Technical Architecture

### Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Language** | Python 3.11+ | Cross-platform, rich ecosystem, rapid development |
| **GUI Framework** | Qt6 (PySide6) | Native look, cross-platform, mature PDF support |
| **PDF Core Library** | PyMuPDF (fitz) | Fast, feature-rich, AGPL/commercial license |
| **PDF Manipulation** | pypdf / pikepdf | Page manipulation, merging, MIT/MPL licensed |
| **Forms & Signatures** | pdfrw + endesive | Form filling and digital signatures |
| **OCR (optional)** | Tesseract + pytesseract | Text extraction from scanned PDFs |
| **Conversion** | pdf2image, python-docx, LibreOffice headless | Format conversion |
| **Build/Package** | PyInstaller / Nuitka | Cross-platform executables |

### Alternative Stack (Performance-focused)

| Component | Technology |
|-----------|------------|
| **Language** | Rust / C++ |
| **GUI Framework** | Qt6 (C++) or Tauri (Rust + Web) |
| **PDF Library** | MuPDF (C), PDFium, or Poppler |

### Project Structure

```
openpdfreader/
├── cmd/
│   └── openpdfreader/
│       └── main.go             # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go              # Main application struct
│   │   └── config.go           # Configuration management
│   ├── ui/
│   │   ├── window.go           # Main window
│   │   ├── viewer.go           # PDF viewer widget
│   │   ├── toolbar.go          # Toolbar component
│   │   ├── sidebar.go          # Thumbnails, bookmarks
│   │   └── dialogs/
│   │       ├── print.go
│   │       ├── merge.go
│   │       ├── convert.go
│   │       ├── signature.go
│   │       └── password.go
│   ├── pdf/
│   │   ├── document.go         # PDF document abstraction
│   │   ├── renderer.go         # Page rendering
│   │   ├── editor.go           # Editing operations
│   │   ├── forms.go            # Form handling
│   │   ├── signatures.go       # Digital signatures
│   │   ├── merger.go           # Merge/split operations
│   │   ├── converter.go        # Format conversion
│   │   ├── security.go         # Password/encryption
│   │   └── redaction.go        # Redaction operations
│   └── utils/
│       ├── file.go
│       ├── image.go
│       └── platform.go
├── pkg/
│   └── pdfcore/                # Public API (if needed)
│       └── types.go
├── assets/
│   ├── icons/
│   ├── fonts/
│   └── themes/
├── testdata/
│   └── pdfs/
├── docs/
│   ├── user_guide.md
│   └── developer_guide.md
├── scripts/
│   ├── build.sh
│   └── release.sh
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
├── CHANGELOG.md
├── TODO.md
├── AGENTS.md
├── project.md
└── .github/
    └── workflows/
        ├── ci.yml
        ├── release.yml
        └── codeql.yml
```

---

## Development Phases

### Phase 1: Foundation (Weeks 1-4)
**Goal:** Basic PDF viewing with core navigation

- [ ] Project setup (Go modules, Fyne, pdfcpu)
- [ ] Main window with menu bar and toolbar
- [ ] PDF document loading and rendering
- [ ] Page navigation (scroll, page controls)
- [ ] Zoom functionality
- [ ] Text selection and copy
- [ ] Thumbnail sidebar
- [ ] Recent files menu
- [ ] Basic error handling

**Deliverable:** Functional PDF viewer

### Phase 2: Page Management (Weeks 5-6)
**Goal:** Page manipulation capabilities

- [ ] Page thumbnails with drag-and-drop reordering
- [ ] Delete pages
- [ ] Rotate pages
- [ ] Extract pages to new PDF
- [ ] Insert pages from another PDF
- [ ] Split PDF functionality
- [ ] Undo/redo for page operations

**Deliverable:** Page management features complete

### Phase 3: Merge & Conversion (Weeks 7-8)
**Goal:** PDF merging and format conversion

- [ ] Merge multiple PDFs dialog
- [ ] Batch merge capability
- [ ] PDF to PNG/JPEG export
- [ ] PDF to text extraction
- [ ] PDF to Word conversion (basic)
- [ ] Images to PDF conversion
- [ ] Progress indicators for long operations

**Deliverable:** Merge and conversion features complete

### Phase 4: Annotations & Editing (Weeks 9-11)
**Goal:** PDF annotation and basic editing

- [ ] Highlight tool
- [ ] Text annotation tool
- [ ] Sticky notes
- [ ] Freehand drawing
- [ ] Shape tools (rectangle, oval, line, arrow)
- [ ] Stamp/image insertion
- [ ] Annotation properties panel
- [ ] Save annotations to PDF

**Deliverable:** Full annotation support

### Phase 5: Forms & Signatures (Weeks 12-14)
**Goal:** Form filling and digital signatures

- [ ] Form field detection
- [ ] Text field filling
- [ ] Checkbox/radio button support
- [ ] Dropdown selection
- [ ] Draw signature pad
- [ ] Image signature upload
- [ ] Signature placement and sizing
- [ ] Basic certificate signature support
- [ ] Flatten forms option

**Deliverable:** Form and signature features complete

### Phase 6: Security & Redaction (Weeks 15-16)
**Goal:** Password protection and redaction

- [ ] Password-protected PDF support
- [ ] Password entry dialog
- [ ] Add password protection
- [ ] Remove password protection
- [ ] Set document permissions
- [ ] Redaction tool (area selection)
- [ ] Search and redact
- [ ] Apply/burn redactions

**Deliverable:** Security and redaction complete

### Phase 7: Printing & Polish (Weeks 17-18)
**Goal:** Printing and UX improvements

- [ ] Print dialog implementation
- [ ] Print preview
- [ ] Print settings (range, quality, layout)
- [ ] Keyboard shortcuts
- [ ] Preferences dialog
- [ ] Dark/light theme support
- [ ] Performance optimization
- [ ] Memory management improvements

**Deliverable:** Print functionality and polished UX

### Phase 8: Packaging & Release (Weeks 19-20)
**Goal:** Distribution-ready application

- [ ] Windows installer (MSI/NSIS)
- [ ] Linux packages (AppImage, deb, rpm, Flatpak)
- [ ] Auto-update mechanism
- [ ] Documentation
- [ ] User guide
- [ ] Release automation
- [ ] Bug fixes and stabilization

**Deliverable:** Version 1.0 release

---

## Open Source Dependencies

### Core Dependencies

| Library | License | Purpose |
|---------|---------|----------|
| fyne.io/fyne/v2 | BSD-3 | GUI framework |
| pdfcpu/pdfcpu | Apache 2.0 | PDF manipulation & validation |
| klippa-app/go-pdfium | MIT | PDF rendering (PDFium bindings) |
| unidoc/unipdf | AGPL / Commercial | Advanced PDF operations |
| disintegration/imaging | MIT | Image processing |
| signintech/gopdf | MIT | PDF generation |

### Optional Dependencies

| Library | License | Purpose |
|---------|---------|----------|
| otiai10/gosseract | MIT | OCR support (Tesseract) |
| skip2/go-qrcode | MIT | QR code generation |
| golang.org/x/crypto | BSD-3 | Encryption support |

### Development Dependencies

| Library | License | Purpose |
|---------|---------|----------|
| testing (stdlib) | BSD-3 | Testing framework |
| stretchr/testify | MIT | Test assertions |
| golangci-lint | GPL-3 | Linting |
| goreleaser | MIT | Release automation |

---

## Platform Considerations

### Windows
- Native file associations (.pdf)
- Windows installer (MSI)
- Portable version (ZIP)
- Context menu integration
- Jump list support

### Linux
- XDG compliance
- Desktop file integration
- AppImage for universal distribution
- Flatpak for sandboxed installation
- DEB/RPM for traditional package managers
- Wayland and X11 support

---

## Performance Requirements

| Metric | Target |
|--------|--------|
| PDF Open (< 100 pages) | < 2 seconds |
| Page Render | < 100ms |
| Memory (100-page PDF) | < 200 MB |
| Startup Time | < 3 seconds |
| Search (1000 pages) | < 5 seconds |

---

## Security Considerations

- Sandboxed PDF parsing
- Input validation for all file operations
- Secure handling of passwords (no plaintext storage)
- Certificate validation for digital signatures
- Regular dependency security audits
- Code signing for releases

---

## Future Enhancements (Post v1.0)

- Cloud storage integration (Google Drive, Dropbox, OneDrive)
- Collaborative annotations
- Compare PDFs (visual diff)
- PDF/A compliance checking
- Batch processing CLI
- Plugin architecture
- Mobile companion app
- Browser extension

---

## Success Metrics

- Functional parity with Adobe Reader core features
- < 5% crash rate
- User satisfaction score > 4.0/5.0
- Active community contributions
- Cross-platform consistency

---

## License

The project will be released under **Apache License 2.0** to maximize compatibility with Go ecosystem libraries and allow commercial use while requiring attribution.

---

## Getting Started (Development)

```bash
# Clone repository
git clone https://github.com/your-org/openpdfreader.git
cd openpdfreader

# Install Go 1.26+ (if not installed)
curl -LO https://go.dev/dl/go1.26.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Fyne dependencies (Linux)
sudo apt-get install gcc libgl1-mesa-dev xorg-dev

# Download Go dependencies
go mod download

# Run application
go run ./cmd/openpdfreader

# Run tests
go test ./...

# Build executable
go build -o openpdfreader ./cmd/openpdfreader

# Cross-compile for Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o openpdfreader.exe ./cmd/openpdfreader
```

---

*Last Updated: February 2026*
