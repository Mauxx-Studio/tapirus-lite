package main

import (
	"tapirus_lite/components"
	"tapirus_lite/db"
	"tapirus_lite/internal/domain/services"
	"tapirus_lite/internal/infrastructure/repository"
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
	//Inicializar Base de Datos
	db := db.DBSetup()

	//Inicializar Repositorios
	productRepo := repository.NewProductRepository(db)
	clientRepo := repository.NewClientRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	//Inicializar Servicios
	productService := services.NewProductService(productRepo)
	clientService := services.NewClientService(clientRepo)
	orderService := services.NewOrderService(orderRepo, productService, clientService)

	a := app.New()
	currentTheme = settings.LoadConfig()
	a.Settings().SetTheme(currentTheme)

	w := a.NewWindow("Tapirus Lite")
	w.Resize(fyne.NewSize(1000, 700))

	nuevoBoton := widget.NewButton("Nuevo Pedido", nil)
	nuevoBoton.OnTapped = func() { components.NewOrderForm(orderService, w, nuevoBoton, nil) }

	botonPedidos := widget.NewButton("Pedidos", func() {
		w.Content().(*fyne.Container).Objects[0] = components.OrderList(orderService, w, nuevoBoton)
	})
	botonProductos := widget.NewButton("Productos", func() {
		w.Content().(*fyne.Container).Objects[0] = components.ProductList(productService, w, nuevoBoton)
	})
	botonClientes := widget.NewButton("Clientes", func() {
		w.Content().(*fyne.Container).Objects[0] = components.ClientList(clientService, w, nuevoBoton)
	})
	botonInicio := widget.NewButton("Inicio", func() { components.ShowMainScreen(orderService, w, nuevoBoton) })
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

/*
	w.Content().(*fyne.Container).Objects[0] = container.NewBorder(headerWithSeparator, nil, nil, nil, container.NewScroll(dataTable))
	w.Content().Refresh()
*/
