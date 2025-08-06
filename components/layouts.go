package components

import (
	"fyne.io/fyne/v2"
)

// FormLayout para dos columnas con tamaños dinámicos
type FormLayout struct {
	separator     float32
	margin        float32
	minEntryWidth float32
	maxLabelWidth float32
}

func (f *FormLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	totalHeight := f.separator
	for i := 0; i < len(objects); i += 2 {
		label := objects[i]
		entry := objects[i+1]
		rowHeight := fyne.Max(label.MinSize().Height, entry.MinSize().Height)
		totalHeight += rowHeight + f.separator
		f.maxLabelWidth = fyne.Max(label.MinSize().Width, f.maxLabelWidth)
	}
	return fyne.NewSize(f.maxLabelWidth+f.minEntryWidth+f.separator+2*f.margin, totalHeight)
}

func (f *FormLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	yPos := f.separator
	EntryWidth := containerSize.Width - f.maxLabelWidth - f.separator - 2*f.margin
	EntryHpos := f.maxLabelWidth + f.separator + f.margin
	for i := 0; i < len(objects); i += 2 {
		label := objects[i]
		entry := objects[i+1]
		rowHeight := fyne.Max(label.MinSize().Height, entry.MinSize().Height)
		labelHPos := f.maxLabelWidth - label.MinSize().Width + f.margin

		label.Move(fyne.NewPos(labelHPos, yPos))
		entry.Move(fyne.NewPos(EntryHpos, yPos))
		entry.Resize(fyne.NewSize(EntryWidth, rowHeight))

		yPos += rowHeight + f.separator
	}
}

type DynTableLayout struct {
	colWidths   []float32
	minDynWidth float32
	separetor   float32
	cols        int
}

func (t *DynTableLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	t.cols = len(t.colWidths)
	minWidth := t.separetor
	rowHeight := float32(0)
	minHeight := t.separetor
	for i, obj := range objects {
		if t.colWidths[i%t.cols] == float32(0) {
			minWidth += t.minDynWidth + t.separetor
		} else {
			if obj.MinSize().Width > t.colWidths[i%t.cols] {
				t.colWidths[i%t.cols] = obj.MinSize().Width
			}
			minWidth += t.colWidths[i%t.cols] + t.separetor
		}
		rowHeight = fyne.Max(rowHeight, obj.MinSize().Height)
		if i%t.cols == t.cols-1 {
			minHeight += rowHeight + t.separetor
			rowHeight = 0
		}
	}
	if rowHeight > 0 {
		minHeight += rowHeight + t.separetor
	}
	return fyne.NewSize(minWidth, minHeight)
}

func (t *DynTableLayout) Layout(objets []fyne.CanvasObject, containerSize fyne.Size) {
	xPos := t.separetor
	dynCount := int(0)
	for width := range t.colWidths {
		if width == 0 {
			dynCount += 1
		}
	}
	dynWidth := (containerSize.Width - float32(t.cols+1)*(t.separetor)) / float32(dynCount)
	yPos := t.separetor
	rowHeight := float32(0)
	for i, obj := range objets {
		obj.Move(fyne.NewPos(xPos, yPos))
		if t.colWidths[i%t.cols] == float32(0) {
			xPos += dynWidth + t.separetor
		} else {
			xPos += t.colWidths[i%t.cols] + t.separetor
		}
		rowHeight = fyne.Max(rowHeight, obj.MinSize().Height)
		if i%t.cols == t.cols-1 {
			xPos = 0
			yPos += rowHeight + t.separetor
			rowHeight = 0
		}
	}
}

type ExpandeVBox struct {
	separator float32
}

func (b *ExpandeVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minWidth := float32(0)
	minHeight := b.separator
	for _, obj := range objects {
		minWidth = fyne.Max(minWidth, obj.MinSize().Width)
		minHeight += obj.MinSize().Height + b.separator
	}
	return fyne.NewSize(minWidth, minHeight)
}

func (b *ExpandeVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	posY := b.separator
	width := containerSize.Width - 2*b.separator
	for _, obj := range objects {
		height := obj.MinSize().Height
		obj.Resize(fyne.NewSize(width, height))
		obj.Move(fyne.NewPos(b.separator, posY))
		posY += obj.MinSize().Height + b.separator
	}
}

type ExpandeHbox struct {
	dynIndex int
}

func (l *ExpandeHbox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	width := float32(0)
	height := float32(0)
	for _, obj := range objects {
		width += obj.MinSize().Width
		height = fyne.Max(height, obj.MinSize().Height)
	}
	return fyne.NewSize(width, height)
}

func (l *ExpandeHbox) Layout(objects []fyne.CanvasObject, contSize fyne.Size) {
	dynWidth := contSize.Width
	height := float32(0)
	for i, obj := range objects {
		if i != l.dynIndex {
			dynWidth -= obj.MinSize().Width
		}
		height = fyne.Max(height, obj.MinSize().Height)
	}
	xPos := float32(0)
	for i, obj := range objects {
		yPos := (height - obj.MinSize().Height) / 2
		obj.Move(fyne.NewPos(xPos, yPos))
		if i != l.dynIndex {
			xPos += obj.MinSize().Width
		} else {
			xPos += dynWidth
			obj.Resize(fyne.NewSize(dynWidth, obj.MinSize().Height))
		}
	}
}

type ItemsRowLayout struct {
	separator   float32
	widths      []float32
	minDynWidth float32
}

func (l *ItemsRowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minWidth := l.separator
	minHeight := float32(0)
	for i, obj := range objects {
		if i < len(l.widths) && l.widths[i] > 0 {
			l.widths[i] = fyne.Max(l.widths[i], obj.MinSize().Width)
			minWidth += l.widths[i] + l.separator
		} else {
			minWidth += l.minDynWidth + l.separator
		}
		minHeight = fyne.Max(minHeight, obj.MinSize().Height)
	}
	return fyne.NewSize(minWidth, minHeight)
}

func (l *ItemsRowLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	xPos := l.separator
	dynCount := 0
	fixedWidth := float32(0)

	// Calcular anchos fijos y contar dinámicos
	for i, obj := range objects {
		if i < len(l.widths) && l.widths[i] > 0 {
			l.widths[i] = fyne.Max(l.widths[i], obj.MinSize().Width)
			fixedWidth += l.widths[i]
		} else {
			dynCount++
		}
	}

	// Calcular ancho dinámico con mínimo respetado
	dynWidth := float32(0)
	if dynCount > 0 {
		dynWidth = (containerSize.Width - fixedWidth - l.separator*float32(len(objects)+1)) / float32(dynCount)
		//		if dynWidth < l.minDynWidth {
		//			dynWidth = l.minDynWidth
		//		}
		dynWidth = fyne.Max(dynWidth, l.minDynWidth)
	}

	// Posicionar y redimensionar elementos
	for i, obj := range objects {
		height := obj.MinSize().Height // Respetamos la altura mínima del widget
		if i < len(l.widths) && l.widths[i] > 0 {
			// Ancho fijo
			obj.Resize(fyne.NewSize(l.widths[i], height))
			obj.Move(fyne.NewPos(xPos, (containerSize.Height-height)/2)) // Centrado vertical
			xPos += l.widths[i] + l.separator
		} else {
			// Ancho dinámico
			obj.Resize(fyne.NewSize(dynWidth, height))
			obj.Move(fyne.NewPos(xPos, (containerSize.Height-height)/2)) // Centrado vertical
			xPos += dynWidth + l.separator
		}
	}
}
