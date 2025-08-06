package settings

import (
	"encoding/json"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	fyne.Theme
	TextSize      float32 // SizeNameText
	HeaderSize    float32 // SizeNameHeadingText
	SubheaderSize float32 // SizeNameSubHeadingText
	Variant       fyne.ThemeVariant
}

type Config struct {
	TextSize      float32           `json:"text_size"`
	HeaderSize    float32           `json:"header_size"`
	SubheaderSize float32           `json:"subheader_size"`
	ThemeVariant  fyne.ThemeVariant `json:"theme_variant"`
}

func NewCustomTheme() *CustomTheme {
	return &CustomTheme{
		Theme:         theme.DefaultTheme(),
		TextSize:      14,
		HeaderSize:    16,
		SubheaderSize: 15,
		Variant:       theme.VariantLight,
	}
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, t.Variant)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return t.TextSize
	case theme.SizeNameHeadingText:
		return t.HeaderSize
	case theme.SizeNameSubHeadingText:
		return t.SubheaderSize
	default:
		return t.Theme.Size(name)
	}
}

func (t *CustomTheme) SetTextSize(size float32) {
	t.TextSize = size
}

func (t *CustomTheme) SetHeaderSize(size float32) {
	t.HeaderSize = size
}

func (t *CustomTheme) SetSubheaderSize(size float32) {
	t.SubheaderSize = size
}

func (t *CustomTheme) SetVariant(variant fyne.ThemeVariant) {
	t.Variant = variant
}

func (t *CustomTheme) GetTextSize() float32 {
	return t.TextSize
}

func (t *CustomTheme) GetHeaderSize() float32 {
	return t.HeaderSize
}

func (t *CustomTheme) GetSubheaderSize() float32 {
	return t.SubheaderSize
}

func (t *CustomTheme) GetVariant() fyne.ThemeVariant {
	return t.Variant
}

func LoadConfig() *CustomTheme {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return NewCustomTheme()
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return NewCustomTheme()
	}
	theme := NewCustomTheme()
	theme.SetTextSize(config.TextSize)
	theme.SetHeaderSize(config.HeaderSize)
	theme.SetSubheaderSize(config.SubheaderSize)
	theme.SetVariant(config.ThemeVariant)
	return theme
}

func SaveConfig(theme *CustomTheme) {
	config := Config{
		TextSize:      theme.GetTextSize(),
		HeaderSize:    theme.GetHeaderSize(),
		SubheaderSize: theme.GetSubheaderSize(),
		ThemeVariant:  theme.GetVariant(),
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile("config.json", data, 0644)
}
