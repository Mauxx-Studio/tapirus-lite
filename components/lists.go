package components

import (
	"fmt"
	"math"
	"sort"
	"tapirus_lite/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func OrderList(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button) {

	var orders []models.Order
	db.Preload("Items.Product").Find(&orders)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].DeliveryDate.Before(orders[j].DeliveryDate)
	})

	dataTable := widget.NewTable(
		func() (int, int) { return len(orders) + 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(orders) {
				switch id.Col {
				case 0:
					label.SetText(fmt.Sprintf("%05d", orders[id.Row].ID))
					label.Alignment = fyne.TextAlignCenter
				case 1:
					label.SetText(orders[id.Row].ClientName)
				case 2:
					label.SetText(orders[id.Row].DeliveryDate.Format("2006-01-02 15:04"))
					label.Alignment = fyne.TextAlignCenter
				case 3:
					label.SetText(fmt.Sprintf("$%.2f", orders[id.Row].Amount))
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		},
	)

	dataTable.SetColumnWidth(0, 80)  // Código
	dataTable.SetColumnWidth(2, 150) // Fecha de Entrega
	dataTable.SetColumnWidth(3, 130) // Monto

	headerTable := widget.NewTable(
		func() (int, int) { return 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("Código")
			case 1:
				label.SetText("Cliente")
			case 2:
				label.SetText("Fecha de Entrega")
			case 3:
				label.SetText("Monto")
			}
			label.Alignment = fyne.TextAlignCenter
		},
	)

	headerTable.SetColumnWidth(0, 80)  // Código
	headerTable.SetColumnWidth(2, 150) // Fecha de Entrega
	headerTable.SetColumnWidth(3, 130) // Monto

	windowWidth := w.Content().Size().Width
	fixedWidths := float32(80 + 150 + 130 + 12)
	clientWidth := float32(math.Max(float64(windowWidth-fixedWidths), 100))
	dataTable.SetColumnWidth(1, clientWidth)
	headerTable.SetColumnWidth(1, clientWidth)

	dataTable.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(orders) {
			NewOrderForm(db, w, nuevoBoton, &orders[id.Row]) // Modo edición
			dataTable.Unselect(id)
		}
	}

	headerWithSeparator := container.NewVBox(
		widget.NewSeparator(),
		headerTable,
		widget.NewSeparator(),
	)

	w.Content().(*fyne.Container).Objects[0] = container.NewBorder(headerWithSeparator, nil, nil, nil, container.NewScroll(dataTable))
	w.Content().Refresh()
	nuevoBoton.SetText("Nuevo Pedido")
	nuevoBoton.OnTapped = func() { NewOrderForm(db, w, nuevoBoton, nil) } // Modo creación
}

func ProductList(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button) {
	var products []models.Product
	db.Find(&products)

	dataTable := widget.NewTable(
		func() (int, int) { return len(products) + 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(products) {
				switch id.Col {
				case 0: // Código
					label.SetText(fmt.Sprintf("%05d", products[id.Row].ID))
					label.Alignment = fyne.TextAlignCenter
				case 1: // Nombre
					label.SetText(products[id.Row].Name)
				case 2: // Precio
					label.SetText(fmt.Sprintf("$%.2f", products[id.Row].Price))
					label.Alignment = fyne.TextAlignCenter
				case 3: // Stock
					label.SetText(fmt.Sprintf("%d %s", products[id.Row].Stock, products[id.Row].Unit))
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		},
	)

	dataTable.SetColumnWidth(0, 80)  // Código
	dataTable.SetColumnWidth(2, 100) // Precio
	dataTable.SetColumnWidth(3, 100) // Stock

	headerTable := widget.NewTable(
		func() (int, int) { return 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("Código")
			case 1:
				label.SetText("Nombre")
			case 2:
				label.SetText("Precio")
			case 3:
				label.SetText("Stock")
			}
			label.Alignment = fyne.TextAlignCenter
		},
	)

	headerTable.SetColumnWidth(0, 80)  // Código
	headerTable.SetColumnWidth(2, 100) // Precio
	headerTable.SetColumnWidth(3, 100) // Stock

	windowWidth := w.Content().Size().Width
	fixedWidths := float32(80 + 100 + 100 + 12)
	nameWidth := float32(math.Max(float64(windowWidth-fixedWidths), 100))
	dataTable.SetColumnWidth(1, nameWidth)
	headerTable.SetColumnWidth(1, nameWidth)

	dataTable.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(products) {
			NewProductForm(db, w, nuevoBoton, &products[id.Row]) // Pasamos el producto para edición
			dataTable.Unselect(id)
		}
	}

	headerWithSeparator := container.NewVBox(
		widget.NewSeparator(),
		headerTable,
		widget.NewSeparator(),
	)

	w.Content().(*fyne.Container).Objects[0] = container.NewBorder(headerWithSeparator, nil, nil, nil, container.NewScroll(dataTable))
	w.Content().Refresh()
	nuevoBoton.SetText("Nuevo Producto")
	nuevoBoton.OnTapped = func() { NewProductForm(db, w, nuevoBoton, nil) } // Nil para modo creación
}

func ClientList(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button) {
	var clients []models.Client
	db.Find(&clients)

	dataTable := widget.NewTable(
		func() (int, int) { return len(clients) + 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(clients) {
				switch id.Col {
				case 0: // Código
					label.SetText(fmt.Sprintf("%05d", clients[id.Row].ID))
					label.Alignment = fyne.TextAlignCenter
				case 1: // Cliente (sin centrar)
					label.SetText(clients[id.Row].Name)
				case 2: // Teléfono (sin centrar)
					label.SetText(clients[id.Row].Phone)
				case 3: // Pedidos
					label.SetText(fmt.Sprintf("%d / %d", clients[id.Row].PendingOrders, clients[id.Row].TotalOrders))
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		},
	)

	dataTable.SetColumnWidth(0, 80)  // Código
	dataTable.SetColumnWidth(2, 150) // Teléfono
	dataTable.SetColumnWidth(3, 100) // Pedidos

	headerTable := widget.NewTable(
		func() (int, int) { return 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("Código")
			case 1:
				label.SetText("Cliente")
			case 2:
				label.SetText("Teléfono")
			case 3:
				label.SetText("Pedidos")
			}
			label.Alignment = fyne.TextAlignCenter
		},
	)

	headerTable.SetColumnWidth(0, 80)  // Código
	headerTable.SetColumnWidth(2, 150) // Teléfono
	headerTable.SetColumnWidth(3, 100) // Pedidos

	windowWidth := w.Content().Size().Width
	fixedWidths := float32(80 + 150 + 100 + 12)
	nameWidth := float32(math.Max(float64(windowWidth-fixedWidths), 100))
	dataTable.SetColumnWidth(1, nameWidth)
	headerTable.SetColumnWidth(1, nameWidth)

	dataTable.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(clients) {
			NewClientForm(db, w, nuevoBoton, &clients[id.Row]) // Modo edición
			dataTable.Unselect(id)
		}
	}

	// Encabezado con separadores arriba y abajo
	headerWithSeparator := container.NewVBox(
		widget.NewSeparator(),
		headerTable,
		widget.NewSeparator(),
	)

	w.Content().(*fyne.Container).Objects[0] = container.NewBorder(headerWithSeparator, nil, nil, nil, container.NewScroll(dataTable))
	w.Content().Refresh()
	nuevoBoton.SetText("Nuevo Cliente")
	nuevoBoton.OnTapped = func() { NewClientForm(db, w, nuevoBoton, nil) } // Modo creación
}
