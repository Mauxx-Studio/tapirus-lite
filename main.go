package main

import (
	"tapirus_lite/components"
	"tapirus_lite/db"
	"tapirus_lite/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var currentTheme *settings.CustomTheme

func main() {
	db := db.DBSetup()

	a := app.New()
	currentTheme = settings.LoadConfig()
	a.Settings().SetTheme(currentTheme)

	w := a.NewWindow("Tapirus Lite")
	w.Resize(fyne.NewSize(1000, 700))

	nuevoBoton := widget.NewButton("Nuevo Pedido", nil)
	nuevoBoton.OnTapped = func() { components.NewOrderForm(db, w, nuevoBoton, nil) }

	botonPedidos := widget.NewButton("Pedidos", func() { components.OrderList(db, w, nuevoBoton) })
	botonProductos := widget.NewButton("Productos", func() { components.ProductList(db, w, nuevoBoton) })
	botonClientes := widget.NewButton("Clientes", func() { components.ClientList(db, w, nuevoBoton) })
	botonInicio := widget.NewButton("Inicio", func() { components.ShowMainScreen(db, w, nuevoBoton) })
	botonSetup := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		components.ShowSetupWindow(a, w, currentTheme, settings.SaveConfig)
	})

	// Barra superior con espaciado
	appBar := container.NewHBox(
		botonPedidos,
		botonProductos,
		botonClientes,
		layout.NewSpacer(),
		nuevoBoton,
		botonInicio,
		botonSetup,
	)

	centerContent := components.MainSummary(db, w)
	w.SetContent(container.NewBorder(appBar, nil, nil, nil, centerContent))
	w.ShowAndRun()
}
