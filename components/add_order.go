package components

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"tapirus_lite/internal/domain/entities"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
	"gorm.io/gorm"
)

func NewOrderForm(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button, order *entities.Order) {
	orderWindow := fyne.CurrentApp().NewWindow("Nuevo Pedido")
	orderWindow.Resize(fyne.NewSize(800, 600))

	var clients []entities.Client
	db.Find(&clients)
	sort.Slice(clients, func(i, j int) bool { return clients[i].Name < clients[j].Name })
	clientNames := make([]string, len(clients))
	clientIDs := make(map[string]uint)
	for i, c := range clients {
		clientNames[i] = c.Name
		clientIDs[c.Name] = c.ID
	}

	var products []entities.Product
	db.Find(&products)
	sort.Slice(products, func(i, j int) bool { return products[i].Name < products[j].Name })
	productNames := make([]string, len(products))
	productIDs := make(map[string]uint)
	productPrices := make(map[string]float64)
	productUnits := make(map[string]string)
	for i, p := range products {
		productNames[i] = p.Name
		productIDs[p.Name] = p.ID
		productPrices[p.Name] = p.Price
		productUnits[p.Name] = p.Unit
	}

	idLabel := widget.NewLabel("crear")
	clientEntry := xWidget.NewCompletionEntry([]string{})
	clientEntry.SetText("Consumidor final")
	isEditComp := order != nil
	clientEntry.OnChanged = func(s string) {
		if len(s) < 2 || isEditComp {
			clientEntry.HideCompletion()
			isEditComp = false
			return
		}
		var matches []string
		sLower := strings.ToLower(s)
		for _, name := range clientNames {
			if strings.Contains(strings.ToLower(name), sLower) {
				matches = append(matches, name)
			}
		}
		if len(matches) == 0 {
			clientEntry.HideCompletion()
			return
		}
		clientEntry.SetOptions(matches)
		clientEntry.ShowCompletion()
	}
	fechaEntry := widget.NewEntry()
	fechaEntry.SetText(time.Now().Format("2006-01-02 15:04"))
	notaEntry := widget.NewMultiLineEntry()

	isEdit := order != nil
	if isEdit {
		orderWindow.SetTitle("Editar Pedido")
		idLabel.SetText(fmt.Sprintf("%07d", order.ID))
		clientEntry.SetText(order.ClientName)
		fechaEntry.SetText(order.DeliveryDate.Format("2006-01-02 15:04"))
		notaEntry.SetText(order.Note)
	}
	topItems := []fyne.CanvasObject{
		widget.NewLabel("Código"), idLabel,
		widget.NewLabel("Cliente"), clientEntry,
		widget.NewLabel("Fecha Entrega"), fechaEntry,
		widget.NewLabel("Nota"), notaEntry,
	}
	topLayout := &FormLayout{separator: 4, margin: 12, minEntryWidth: 150}
	topContainer := container.New(topLayout, topItems...)

	type ItemRow struct {
		numLabel     *widget.Label
		cantidad     *widget.Label
		unidad       *widget.Label
		producto     *widget.Label
		precioUnit   *widget.Label
		montoParcial *widget.Label
		deleteButton *widget.Button
	}
	var items []ItemRow
	itemsContainer := container.NewVBox()

	rowLayout := &ItemsRowLayout{
		separator:   6,
		widths:      []float32{40, 80, 30, 0, 100, 100, 30},
		minDynWidth: 150,
	}

	numHeader := widget.NewLabel("Nº")
	qtyHeader := widget.NewLabel("Cantidad")
	unitHeader := widget.NewLabel("")
	prodHeader := widget.NewLabel("Producto")
	priceHeader := widget.NewLabel("Precio")
	montoHeader := widget.NewLabel("Monto")
	deleteHeader := widget.NewLabel("")
	header := container.New(rowLayout, numHeader, qtyHeader, unitHeader, prodHeader, priceHeader, montoHeader, deleteHeader)

	totalLabel := widget.NewLabel("0.00")
	updateTotal := func() {
		total := 0.0
		for _, item := range items {
			if monto, err := strconv.ParseFloat(item.montoParcial.Text, 64); err == nil {
				total += monto
			}
		}
		totalLabel.SetText(fmt.Sprintf("%.2f", total))
	}

	deleteItem := func(index int) {
		if index < 0 || index >= len(items) {
			return
		}
		items = append(items[:index], items[index+1:]...)
		itemsContainer.Objects = itemsContainer.Objects[:0]
		for i, item := range items {
			item.numLabel.SetText(fmt.Sprintf("%d", i+1))
			row := container.New(rowLayout, item.numLabel, item.cantidad, item.unidad, item.producto, item.precioUnit, item.montoParcial, item.deleteButton)
			itemsContainer.Add(row)
		}
		updateTotal()
		itemsContainer.Refresh()
	}

	makeDeleteCallback := func(btn *widget.Button) func() {
		return func() {
			for i, item := range items {
				if item.deleteButton == btn {
					dialog.ShowConfirm("Confirmar", "¿Eliminar este ítem?", func(confirmed bool) {
						if confirmed {
							deleteItem(i)
						}
					}, orderWindow)
					return
				}
			}
		}
	}

	if isEdit && order != nil && len(order.Items) > 0 {
		for i, item := range order.Items {
			numLabel := widget.NewLabel(fmt.Sprintf("%d", i+1))
			cantidad := widget.NewLabel(fmt.Sprintf("%.2f", item.Quantity))
			var productName, productUnit string
			var productPrice float64
			for _, p := range products {
				if p.ID == item.ProductID {
					productName = p.Name
					productUnit = p.Unit
					productPrice = p.Price
					break
				}
			}
			unidad := widget.NewLabel(productUnit)
			producto := widget.NewLabel(productName)
			precioUnit := widget.NewLabel(fmt.Sprintf("%.2f", productPrice))
			montoParcial := widget.NewLabel(fmt.Sprintf("%.2f", item.Quantity*productPrice))
			deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			deleteButton.OnTapped = makeDeleteCallback(deleteButton)
			items = append(items, ItemRow{numLabel, cantidad, unidad, producto, precioUnit, montoParcial, deleteButton})
			row := container.New(rowLayout, numLabel, cantidad, unidad, producto, precioUnit, montoParcial, deleteButton)
			itemsContainer.Add(row)
		}
		updateTotal()
	}

	addItemButton := widget.NewButton("Agregar Ítem", func() {
		productoEntry := xWidget.NewCompletionEntry(productNames)
		productoEntry.SetPlaceHolder("Seleccione un producto...")
		cantidadEntry := widget.NewEntry()
		cantidadEntry.SetPlaceHolder("Ingrese la cantidad")

		priceUnitLabel := widget.NewLabel("Precio Unitario: -")
		unidadLabel := widget.NewLabel("-")
		montoParcialLabel := widget.NewLabel("Monto Parcial: -")

		productoEntry.OnChanged = func(s string) {
			if len(s) < 2 {
				productoEntry.HideCompletion()
				if price, ok := productPrices[s]; ok {
					priceUnitLabel.SetText(fmt.Sprintf("Precio Unitario: %.2f", price))
					unidadLabel.SetText(productUnits[s])
					if qty, err := strconv.ParseFloat(cantidadEntry.Text, 64); err == nil && qty > 0 {
						montoParcialLabel.SetText(fmt.Sprintf("Monto Parcial: %.2f", qty*price))
					} else {
						montoParcialLabel.SetText("Monto Parcial: -")
					}
				} else {
					priceUnitLabel.SetText("Precio Unitario: -")
					unidadLabel.SetText("-")
					montoParcialLabel.SetText("Monto Parcial: -")
				}
				return
			}

			var matches []string
			sLower := strings.ToLower(s)
			for _, name := range productNames {
				if strings.Contains(strings.ToLower(name), sLower) {
					matches = append(matches, name)
				}
			}

			if len(matches) == 0 {
				productoEntry.HideCompletion()
			} else {
				productoEntry.SetOptions(matches)
				productoEntry.ShowCompletion()
			}

			if price, ok := productPrices[s]; ok {
				priceUnitLabel.SetText(fmt.Sprintf("Precio Unitario: %.2f", price))
				unidadLabel.SetText(productUnits[s])
				if qty, err := strconv.ParseFloat(cantidadEntry.Text, 64); err == nil && qty > 0 {
					montoParcialLabel.SetText(fmt.Sprintf("Monto Parcial: %.2f", qty*price))
				} else {
					montoParcialLabel.SetText("Monto Parcial: -")
				}
			} else {
				priceUnitLabel.SetText("Precio Unitario: -")
				unidadLabel.SetText("-")
				montoParcialLabel.SetText("Monto Parcial: -")
			}
		}

		cantidadEntry.OnChanged = func(s string) {
			if qty, err := strconv.ParseFloat(s, 64); err == nil && qty > 0 {
				if price, ok := productPrices[productoEntry.Text]; ok {
					montoParcialLabel.SetText(fmt.Sprintf("Monto Parcial: %.2f", qty*price))
				}
			} else {
				montoParcialLabel.SetText("Monto Parcial: -")
			}
		}

		expVBoxLayout := &ExpandeHbox{dynIndex: 1}
		content := container.NewVBox(
			container.New(
				expVBoxLayout,
				widget.NewLabel("Producto:"),
				productoEntry,
			),
			container.New(
				expVBoxLayout,
				widget.NewLabel("Cantidad:"),
				cantidadEntry,
				unidadLabel,
			),
			priceUnitLabel,
			montoParcialLabel,
		)

		customDialog := dialog.NewCustomConfirm(
			"Agregar Ítem",
			"Aceptar",
			"Cancelar",
			content,
			func(confirm bool) {
				if confirm {
					qty, err := strconv.ParseFloat(cantidadEntry.Text, 64)
					if err != nil || qty <= 0 || productoEntry.Text == "" {
						dialog.ShowError(fmt.Errorf("cantidad inválida o producto no seleccionado"), orderWindow)
						return
					}
					price, ok := productPrices[productoEntry.Text]
					if !ok {
						dialog.ShowError(fmt.Errorf("producto no válido"), orderWindow)
						return
					}
					numLabel := widget.NewLabel(fmt.Sprintf("%d", len(items)+1))
					cantidad := widget.NewLabel(fmt.Sprintf("%.2f", qty))
					unidad := widget.NewLabel(productUnits[productoEntry.Text])
					producto := widget.NewLabel(productoEntry.Text)
					precioUnit := widget.NewLabel(fmt.Sprintf("%.2f", price))
					montoParcial := widget.NewLabel(fmt.Sprintf("%.2f", qty*price))
					deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
					deleteButton.OnTapped = makeDeleteCallback(deleteButton)
					items = append(items, ItemRow{numLabel, cantidad, unidad, producto, precioUnit, montoParcial, deleteButton})
					row := container.New(rowLayout, numLabel, cantidad, unidad, producto, precioUnit, montoParcial, deleteButton)
					itemsContainer.Add(row)
					updateTotal()
					itemsContainer.Refresh()
				}
			},
			orderWindow,
		)

		customDialog.Resize(fyne.NewSize(400, 350))
		customDialog.Show()
	})

	itemsWithHeader := container.NewVBox(header, widget.NewSeparator(), itemsContainer)
	itemsScroll := container.NewScroll(itemsWithHeader)
	itemsScroll.SetMinSize(fyne.NewSize(0, 300))

	totalRow := container.New(rowLayout,
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel("Total:"),
		totalLabel,
		widget.NewLabel(""),
	)

	saveButton := widget.NewButton("Guardar", nil)
	cancelButton := widget.NewButton("Cancelar", nil)
	deleteButton := widget.NewButton("Eliminar", nil)
	deleteButton.Disable()
	if isEdit {
		deleteButton.Enable()
	}
	buttons := container.NewHBox(addItemButton, layout.NewSpacer(), deleteButton, saveButton, cancelButton)

	bottomContainer := container.NewVBox(totalRow, buttons)
	content := container.NewBorder(topContainer, bottomContainer, nil, nil, itemsScroll)
	orderWindow.SetContent(content)

	// Función para verificar si hay cambios
	hasChanges := func() bool {
		if !isEdit {
			// Modo creación: hay cambios si se ingresó algo
			return len(items) > 0 || clientEntry.Text != "" || notaEntry.Text != ""
		}

		// Modo edición: comparar con el pedido original
		if clientEntry.Text != order.ClientName || notaEntry.Text != order.Note ||
			fechaEntry.Text != order.DeliveryDate.Format("2006-01-02 15:04") {
			return true
		}

		// Comparar ítems
		if len(items) != len(order.Items) {
			return true
		}
		for i, item := range items {
			origItem := order.Items[i]
			qty, err := strconv.ParseFloat(item.cantidad.Text, 64)
			if err != nil || qty != origItem.Quantity {
				return true
			}
			prodID, ok := productIDs[item.producto.Text]
			if !ok || prodID != origItem.ProductID {
				return true
			}
		}
		return false
	}

	// Función para manejar el cierre con confirmación
	confirmClose := func() bool {
		if hasChanges() {
			confirmed := false
			dialog.ShowConfirm("Confirmar", "¿Desea salir sin guardar los cambios?", func(conf bool) {
				confirmed = conf
				if confirmed {
					orderWindow.Close()
				}
			}, orderWindow)
			return confirmed
		}
		return true
	}

	// Asignar el manejador al botón Cancelar
	cancelButton.OnTapped = func() {
		if confirmClose() {
			orderWindow.Close()
		}
	}

	// Interceptar el cierre de la ventana (X)
	orderWindow.SetCloseIntercept(func() {
		if confirmClose() {
			orderWindow.Close()
		}
	})

	saveButton.OnTapped = func() {
		fecha, err := time.Parse("2006-01-02 15:04", fechaEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("fecha inválida"), orderWindow)
			return
		}
		clientText := clientEntry.Text
		if clientText == "" {
			dialog.ShowError(fmt.Errorf("debe ingresar un cliente"), orderWindow)
			return
		}

		var orderItems []entities.OrderItem
		totalAmount := 0.0
		for _, item := range items {
			qty, err := strconv.ParseFloat(item.cantidad.Text, 64)
			if err != nil || qty <= 0 || item.producto.Text == "" {
				continue
			}
			prodID, ok := productIDs[item.producto.Text]
			if !ok {
				continue
			}
			price := productPrices[item.producto.Text]
			orderItems = append(orderItems, entities.OrderItem{
				ProductID: prodID,
				Quantity:  qty,
			})
			totalAmount += qty * price
		}

		if len(orderItems) == 0 {
			dialog.ShowError(fmt.Errorf("debe ingresar al menos un ítem"), orderWindow)
			return
		}

		clientID, ok := clientIDs[clientText]
		if !ok {
			dialog.ShowError(fmt.Errorf("cliente no encontrado"), orderWindow)
			return
		}
		var client entities.Client
		db.First(&client, clientID)

		if isEdit {
			order.ClientID = clientID
			order.ClientName = client.Name
			order.DeliveryDate = fecha
			order.Note = notaEntry.Text
			order.Amount = totalAmount
			db.Save(order)
			db.Where("order_id = ?", order.ID).Delete(&entities.OrderItem{})
			for i := range orderItems {
				orderItems[i].OrderID = order.ID
			}
			db.Create(&orderItems)
		} else {
			newOrder := entities.Order{
				ClientID:     clientID,
				ClientName:   client.Name,
				DeliveryDate: fecha,
				Note:         notaEntry.Text,
				Completed:    false,
				Amount:       totalAmount,
				Items:        orderItems,
			}
			db.Create(&newOrder)
			idLabel.SetText(fmt.Sprintf("%07d", newOrder.ID))
		}

		dialog.ShowInformation("Éxito", "Pedido guardado", orderWindow)
		OrderList(db, w, nuevoBoton)
		orderWindow.Close()
	}

	if isEdit {
		deleteButton.OnTapped = func() {
			dialog.ShowConfirm("Confirmar", "¿Eliminar este pedido?", func(confirmed bool) {
				if confirmed {
					db.Delete(order)
					dialog.ShowInformation("Éxito", "Pedido eliminado", orderWindow)
					OrderList(db, w, nuevoBoton)
					orderWindow.Close()
				}
			}, orderWindow)
		}
	}

	orderWindow.Show()
}
