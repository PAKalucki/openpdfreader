package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

var (
	cliMerge = func(inputs []string, output string) error {
		return pdf.NewMerger().Merge(inputs, output)
	}
	cliSplit = func(input, outputDir string) error {
		return pdf.NewMerger().Split(input, outputDir)
	}
	cliExportImages = func(input, outputDir, format string, scale float64) error {
		_, err := pdf.NewImageExporter().ExportToImages(input, outputDir, format, scale)
		return err
	}
	cliExportText = func(input, output string) error {
		return pdf.NewTextExporter().ExportToText(input, output)
	}
)

// RunCLI executes non-GUI PDF operations.
func RunCLI(args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printCLIUsage(out)
		return nil
	}

	switch args[0] {
	case "merge":
		return runMergeCommand(args[1:], out)
	case "split":
		return runSplitCommand(args[1:], out)
	case "export-images":
		return runExportImagesCommand(args[1:], out)
	case "export-text":
		return runExportTextCommand(args[1:], out)
	default:
		return fmt.Errorf("unknown CLI command: %s", args[0])
	}
}

func runMergeCommand(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	inputsFlag := fs.String("inputs", "", "Comma-separated input PDF files")
	outputFlag := fs.String("output", "", "Output PDF file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	inputs := parseCSV(*inputsFlag)
	if len(inputs) < 2 {
		return errors.New("merge requires at least 2 input files")
	}
	if strings.TrimSpace(*outputFlag) == "" {
		return errors.New("merge requires --output")
	}

	if err := cliMerge(inputs, strings.TrimSpace(*outputFlag)); err != nil {
		return err
	}
	fmt.Fprintf(out, "Merged %d file(s) into %s\n", len(inputs), strings.TrimSpace(*outputFlag))
	return nil
}

func runSplitCommand(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	inputFlag := fs.String("input", "", "Input PDF file")
	outputDirFlag := fs.String("output-dir", "", "Output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*inputFlag) == "" {
		return errors.New("split requires --input")
	}
	if strings.TrimSpace(*outputDirFlag) == "" {
		return errors.New("split requires --output-dir")
	}

	if err := cliSplit(strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputDirFlag)); err != nil {
		return err
	}
	fmt.Fprintf(out, "Split %s into %s\n", strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputDirFlag))
	return nil
}

func runExportImagesCommand(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("export-images", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	inputFlag := fs.String("input", "", "Input PDF file")
	outputDirFlag := fs.String("output-dir", "", "Output directory")
	formatFlag := fs.String("format", "png", "Image format: png or jpg")
	scaleFlag := fs.Float64("scale", 2.0, "Render scale")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*inputFlag) == "" {
		return errors.New("export-images requires --input")
	}
	if strings.TrimSpace(*outputDirFlag) == "" {
		return errors.New("export-images requires --output-dir")
	}
	if *scaleFlag <= 0 {
		return errors.New("export-images requires --scale > 0")
	}

	if err := cliExportImages(strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputDirFlag), strings.TrimSpace(*formatFlag), *scaleFlag); err != nil {
		return err
	}
	fmt.Fprintf(out, "Exported images from %s into %s\n", strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputDirFlag))
	return nil
}

func runExportTextCommand(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("export-text", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	inputFlag := fs.String("input", "", "Input PDF file")
	outputFlag := fs.String("output", "", "Output text file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*inputFlag) == "" {
		return errors.New("export-text requires --input")
	}
	if strings.TrimSpace(*outputFlag) == "" {
		return errors.New("export-text requires --output")
	}

	if err := cliExportText(strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputFlag)); err != nil {
		return err
	}
	fmt.Fprintf(out, "Exported text from %s into %s\n", strings.TrimSpace(*inputFlag), strings.TrimSpace(*outputFlag))
	return nil
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			values = append(values, v)
		}
	}
	return values
}

func printCLIUsage(out io.Writer) {
	fmt.Fprintln(out, "OpenPDF Reader CLI mode")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  openpdfreader --cli <command> [options]")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Commands:")
	fmt.Fprintln(out, "  merge          --inputs a.pdf,b.pdf --output out.pdf")
	fmt.Fprintln(out, "  split          --input in.pdf --output-dir ./out")
	fmt.Fprintln(out, "  export-images  --input in.pdf --output-dir ./out --format png --scale 2.0")
	fmt.Fprintln(out, "  export-text    --input in.pdf --output out.txt")
}
