package dialogs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowSignaturePadDialog shows a dialog for drawing a signature.
func ShowSignaturePadDialog(window fyne.Window, onApply func(signaturePNG []byte) error) {
	pad := NewSignaturePad(520, 180)
	help := widget.NewLabel("Draw your signature below, then click Apply Signature.")

	var dlg dialog.Dialog

	clearBtn := widget.NewButton("Clear", func() {
		pad.Clear()
	})
	cancelBtn := widget.NewButton("Cancel", func() {
		if dlg != nil {
			dlg.Hide()
		}
	})
	applyBtn := widget.NewButton("Apply Signature", func() {
		if !pad.HasInk() {
			dialog.ShowInformation("Empty Signature", "Please draw a signature first.", window)
			return
		}

		signaturePNG, err := pad.PNG()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if err := onApply(signaturePNG); err != nil {
			dialog.ShowError(err, window)
			return
		}

		if dlg != nil {
			dlg.Hide()
		}
	})
	applyBtn.Importance = widget.HighImportance

	content := container.NewBorder(
		help,
		container.NewHBox(clearBtn, widget.NewSeparator(), cancelBtn, applyBtn),
		nil,
		nil,
		pad,
	)

	dlg = dialog.NewCustom("Signature Pad", "Close", content, window)
	dlg.Resize(fyne.NewSize(620, 320))
	dlg.Show()
}

// SignaturePad is a drawable canvas for signatures.
type SignaturePad struct {
	widget.BaseWidget
	img      *image.RGBA
	imageObj *canvas.Image
	drawing  bool
	last     image.Point
	hasInk   bool
}

// NewSignaturePad creates a new signature pad with the given pixel dimensions.
func NewSignaturePad(width, height int) *SignaturePad {
	if width < 100 {
		width = 100
	}
	if height < 60 {
		height = 60
	}

	p := &SignaturePad{
		img: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	p.ExtendBaseWidget(p)
	return p
}

// MinSize returns the minimum size for the pad.
func (p *SignaturePad) MinSize() fyne.Size {
	return fyne.NewSize(float32(p.img.Rect.Dx()), float32(p.img.Rect.Dy()))
}

// CreateRenderer creates the widget renderer.
func (p *SignaturePad) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.White)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{170, 170, 170, 255}
	border.StrokeWidth = 1

	p.imageObj = canvas.NewImageFromImage(p.img)
	p.imageObj.FillMode = canvas.ImageFillStretch
	p.imageObj.ScaleMode = canvas.ImageScaleFastest

	content := container.NewMax(bg, p.imageObj, border)
	return widget.NewSimpleRenderer(content)
}

// Tapped draws a point at the tap location.
func (p *SignaturePad) Tapped(ev *fyne.PointEvent) {
	pt := p.clampPoint(ev.Position)
	p.drawPoint(pt)
	p.hasInk = true
	p.Refresh()
}

// Dragged draws a stroke while dragging.
func (p *SignaturePad) Dragged(ev *fyne.DragEvent) {
	pt := p.clampPoint(ev.Position)
	if !p.drawing {
		p.drawing = true
		p.last = pt
		p.drawPoint(pt)
		p.hasInk = true
		p.Refresh()
		return
	}

	p.drawLine(p.last, pt)
	p.last = pt
	p.hasInk = true
	p.Refresh()
}

// DragEnd ends the current stroke.
func (p *SignaturePad) DragEnd() {
	p.drawing = false
}

// HasInk returns true if any stroke has been drawn.
func (p *SignaturePad) HasInk() bool {
	return p.hasInk
}

// Clear erases all strokes.
func (p *SignaturePad) Clear() {
	for y := 0; y < p.img.Rect.Dy(); y++ {
		for x := 0; x < p.img.Rect.Dx(); x++ {
			p.img.Set(x, y, color.Transparent)
		}
	}
	p.hasInk = false
	p.drawing = false
	p.Refresh()
}

// PNG returns the signature image encoded as PNG.
func (p *SignaturePad) PNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, p.img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *SignaturePad) clampPoint(pos fyne.Position) image.Point {
	x := int(pos.X)
	y := int(pos.Y)
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	maxX := p.img.Rect.Dx() - 1
	maxY := p.img.Rect.Dy() - 1
	if x > maxX {
		x = maxX
	}
	if y > maxY {
		y = maxY
	}
	return image.Pt(x, y)
}

func (p *SignaturePad) drawPoint(pt image.Point) {
	ink := color.RGBA{20, 20, 20, 255}
	radius := 2
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			x := pt.X + dx
			y := pt.Y + dy
			if x >= 0 && x < p.img.Rect.Dx() && y >= 0 && y < p.img.Rect.Dy() {
				p.img.Set(x, y, ink)
			}
		}
	}
}

func (p *SignaturePad) drawLine(start, end image.Point) {
	dx := absInt(end.X - start.X)
	dy := absInt(end.Y - start.Y)
	sx := -1
	if start.X < end.X {
		sx = 1
	}
	sy := -1
	if start.Y < end.Y {
		sy = 1
	}

	err := dx - dy
	x := start.X
	y := start.Y

	for {
		p.drawPoint(image.Pt(x, y))
		if x == end.X && y == end.Y {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
