# OpenPDF Reader

An open-source, cross-platform PDF viewer and editor.

## Features

- **View PDFs** - Open and navigate PDF documents with zoom, scroll, and page controls
- **Tabbed Documents** - Open multiple PDF files in separate tabs
- **Print** - Send the currently opened PDF to the system default printer
- **Text Copy** - Select all text on the current page and copy it to clipboard
- **Edit** - Add annotations, highlights, and notes
- **Fill & Sign** - Complete form fields and add signatures
- **Page Management** - Delete, reorder, rotate, extract, and merge pages
- **Conversion** - Export to images and other formats
- **Security** - Support for password-protected PDFs

## Requirements

- Go 1.26 or later
- GCC (for CGO dependencies)
- OpenGL development libraries

### Linux Dependencies

```bash
# Ubuntu/Debian
sudo apt-get install gcc libgl1-mesa-dev xorg-dev

# Fedora
sudo dnf install gcc mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libXxf86vm-devel
```

### Windows Dependencies

- MinGW-w64 or TDM-GCC

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/your-org/openpdfreader.git
cd openpdfreader

# Download dependencies
go mod download

# Build
make build

# Or run directly
make run
```

### Build for Distribution

```bash
# Linux
make build

# Windows (cross-compile from Linux)
make build-windows

# All platforms (requires fyne-cross)
make cross-compile
```

## Usage

```bash
# Run the application
./build/openpdfreader

# Open a specific file
./build/openpdfreader /path/to/document.pdf
```

## Development

```bash
# Run with hot reload (via make)
make dev

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

## Project Structure

```
openpdfreader/
├── cmd/openpdfreader/     # Application entry point
├── internal/
│   ├── app/               # Application logic and config
│   ├── ui/                # User interface components
│   ├── pdf/               # PDF handling
│   └── utils/             # Utility functions
├── pkg/pdfcore/           # Public API (if needed)
├── assets/                # Icons, fonts, themes
├── testdata/              # Test PDF files
└── docs/                  # Documentation
```

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please read the project guidelines before submitting PRs.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## Acknowledgments

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library
