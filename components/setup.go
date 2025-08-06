package components

import (
	"fmt"
	"tapirus_lite/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowSetupWindow(a fyne.App, parent fyne.Window, themeObj any, saveFunc func(*settings.CustomTheme)) {
	customTheme, ok := themeObj.(*settings.CustomTheme)
	if !ok {
		return
	}

	setupWindow := a.NewWindow("Setup")
	setupWindow.Resize(fyne.NewSize(400, 300))

	// Text size
	textLabel := widget.NewLabel(fmt.Sprintf("Texto: %.0f", customTheme.TextSize))
	textSlider := widget.NewSlider(10, 24)
	textSlider.Value = float64(customTheme.TextSize)
	textSlider.OnChanged = func(value float64) {
		customTheme.SetTextSize(float32(value))
		textLabel.SetText(fmt.Sprintf("Texto: %.0f", value))
		a.Settings().SetTheme(customTheme)
	}

	// Header size
	headerLabel := widget.NewLabel(fmt.Sprintf("Título: %.0f", customTheme.HeaderSize))
	headerSlider := widget.NewSlider(10, 24)
	headerSlider.Value = float64(customTheme.HeaderSize)
	headerSlider.OnChanged = func(value float64) {
		customTheme.SetHeaderSize(float32(value))
		headerLabel.SetText(fmt.Sprintf("Título: %.0f", value))
		a.Settings().SetTheme(customTheme)
	}

	// Subheader size
	subheaderLabel := widget.NewLabel(fmt.Sprintf("Subtítulo: %.0f", customTheme.SubheaderSize))
	subheaderSlider := widget.NewSlider(10, 24)
	subheaderSlider.Value = float64(customTheme.SubheaderSize)
	subheaderSlider.OnChanged = func(value float64) {
		customTheme.SetSubheaderSize(float32(value))
		subheaderLabel.SetText(fmt.Sprintf("Subtítulo: %.0f", value))
		a.Settings().SetTheme(customTheme)
	}

	// Theme selector
	themeOptions := []string{"Claro", "Oscuro"}
	themeSelect := widget.NewSelect(themeOptions, func(value string) {
		switch value {
		case "Claro":
			customTheme.SetVariant(theme.VariantLight)
		case "Oscuro":
			customTheme.SetVariant(theme.VariantDark)
		}
		a.Settings().SetTheme(customTheme)
	})
	switch customTheme.GetVariant() {
	case theme.VariantLight:
		themeSelect.SetSelected("Claro")
	case theme.VariantDark:
		themeSelect.SetSelected("Oscuro")
	default:
		themeSelect.SetSelected("Oscuro")
	}

	// Botones
	saveButton := widget.NewButton("Guardar", func() {
		saveFunc(customTheme)
		setupWindow.Close()
	})
	closeButton := widget.NewButton("Cerrar", func() {
		setupWindow.Close()
	})

	formLayout := &FormLayout{
		separator:     4,
		margin:        12,
		minEntryWidth: 160,
	}

	formItems := []fyne.CanvasObject{
		textLabel, textSlider,
		headerLabel, headerSlider,
		subheaderLabel, subheaderSlider,
	}

	content := container.NewVBox(
		widget.NewLabel("Configuración de Tamaños de Fuente"),
		container.New(formLayout, formItems...),
		widget.NewLabel("Tema"),
		themeSelect,
		container.NewHBox(saveButton, closeButton),
	)
	setupWindow.SetContent(content)
	setupWindow.Show()
}
