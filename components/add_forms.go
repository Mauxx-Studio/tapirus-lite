package components

import (
	"fmt"
	"strconv"
	"strings"
	"tapirus_lite/internal/domain/entities"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func NewProductForm(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button, product *entities.Product) {
	// Campos del formulario
	idLabel := widget.NewLabel("")
	nombreEntry := widget.NewEntry()
	descripcionEntry := widget.NewMultiLineEntry()
	precioEntry := widget.NewEntry()
	disponibilidadEntry := widget.NewEntry()
	unidadEntry := widget.NewEntry()
	totalSoldLabel := widget.NewLabel("")

	// Botones
	saveButton := widget.NewButton("Guardar", nil)
	cancelButton := widget.NewButton("Cancelar", nil)
	deleteButton := widget.NewButton("Eliminar", nil)
	deleteButton.Disable()

	// Modo edición o creación
	isEdit := product != nil
	if isEdit {
		idLabel.SetText(fmt.Sprintf("%05d", product.ID))
		nombreEntry.SetText(product.Name)
		descripcionEntry.SetText(product.Description)
		precioEntry.SetText(fmt.Sprintf("%.2f", product.Price))
		disponibilidadEntry.SetText(fmt.Sprintf("%d", product.Stock))
		unidadEntry.SetText(product.Unit)
		totalSoldLabel.SetText(fmt.Sprintf("%.2f", product.TotalSold))
		deleteButton.Enable()
	} else {
		idLabel.SetText("crear")
		totalSoldLabel.SetText("0.00")
	}

	// Ítems del formulario
	formItems := []fyne.CanvasObject{
		widget.NewLabel("Código"), idLabel,
		widget.NewLabel("Nombre"), nombreEntry,
		widget.NewLabel("Descripción"), descripcionEntry,
		widget.NewLabel("Precio"), precioEntry,
		widget.NewLabel("Disponibilidad"), disponibilidadEntry,
		widget.NewLabel("Unidad (kg, un, lts)"), unidadEntry,
		widget.NewLabel("Total Vendido"), totalSoldLabel,
	}

	// Crear contenedor con layout personalizado
	formLayout := &FormLayout{
		separator:     4,
		margin:        12,
		minEntryWidth: 150,
	}
	formContainer := container.New(formLayout, formItems...)

	// Botones centrados
	buttons := container.NewHBox(deleteButton, cancelButton, saveButton)
	centeredButtons := container.NewCenter(buttons)

	// Contenido completo del formulario con título centrado
	formContent := container.NewVBox(
		container.NewCenter(widget.NewLabelWithStyle("Producto", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})), //titulo centrdo y en negrita
		formContainer,
		centeredButtons,
	)

	// Crear el popup modal
	popup := widget.NewModalPopUp(formContent, w.Canvas())
	popup.Resize(fyne.NewSize(500, 0))

	// Lógica del botón Guardar con validaciones simplificadas
	saveButton.OnTapped = func() {
		// Validaciones mínimas
		nombre := strings.TrimSpace(nombreEntry.Text)
		unidad := strings.TrimSpace(unidadEntry.Text)
		if len(nombre) <= 2 {
			dialog.ShowError(fmt.Errorf("el nombre debe tener más de 2 letras"), w)
			return
		}
		if unidad == "" {
			dialog.ShowError(fmt.Errorf("la unidad no puede estar vacía"), w)
			return
		}

		// Parseo con valores por defecto
		precio, _ := strconv.ParseFloat(precioEntry.Text, 64)
		disponibilidad, _ := strconv.Atoi(disponibilidadEntry.Text)

		// Guardar o crear el producto
		if isEdit {
			product.Name = nombre
			product.Description = descripcionEntry.Text
			product.Price = precio
			product.Stock = disponibilidad
			product.Unit = unidad
			db.Save(product)
		} else {
			newProduct := entities.Product{
				Name:        nombre,
				Description: descripcionEntry.Text,
				Price:       precio,
				Stock:       disponibilidad,
				Unit:        unidad,
				TotalSold:   0,
			}
			db.Create(&newProduct)
			idLabel.SetText(fmt.Sprintf("%05d", newProduct.ID))
		}
		ProductList(db, w, nuevoBoton)
		popup.Hide()
	}

	// Lógica del botón Cancelar
	cancelButton.OnTapped = func() {
		popup.Hide()
	}

	// Lógica del botón Eliminar
	if isEdit {
		deleteButton.OnTapped = func() {
			dialog.ShowConfirm("Confirmar", "¿Eliminar este producto?", func(confirmed bool) {
				if confirmed {
					db.Delete(product)
					dialog.ShowInformation("Éxito", "Producto eliminado", w)
					ProductList(db, w, nuevoBoton)
					popup.Hide()
				}
			}, w)
		}
	}

	// Mostrar el popup
	popup.Show()
}

func NewClientForm(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button, client *entities.Client) {
	// Campos del formulario
	idLabel := widget.NewLabel("")
	nombreEntry := widget.NewEntry()
	telefonoEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	cuitEntry := widget.NewEntry()
	direccionEntry := widget.NewMultiLineEntry()
	pedidosLabel := widget.NewLabel("") // Nuevo campo para Pedidos

	// Botones
	saveButton := widget.NewButton("Guardar", nil)
	cancelButton := widget.NewButton("Cancelar", nil)
	deleteButton := widget.NewButton("Eliminar", nil)
	deleteButton.Disable()

	// Modo edición o creación
	isEdit := client != nil
	if isEdit {
		idLabel.SetText(fmt.Sprintf("%05d", client.ID))
		nombreEntry.SetText(client.Name)
		telefonoEntry.SetText(client.Phone)
		emailEntry.SetText(client.Email)
		cuitEntry.SetText(client.CUIT)
		direccionEntry.SetText(client.Address)
		pedidosLabel.SetText(fmt.Sprintf("%d / %d", client.PendingOrders, client.TotalOrders)) // Valor desde el cliente
		deleteButton.Enable()
	} else {
		idLabel.SetText("crear")
		pedidosLabel.SetText("0 / 0") // Valor por defecto para nuevo cliente
	}

	// Ítems del formulario
	formItems := []fyne.CanvasObject{
		widget.NewLabel("Código"), idLabel,
		widget.NewLabel("Nombre"), nombreEntry,
		widget.NewLabel("Teléfono"), telefonoEntry,
		widget.NewLabel("Email"), emailEntry,
		widget.NewLabel("CUIT"), cuitEntry,
		widget.NewLabel("Dirección"), direccionEntry,
		widget.NewLabel("Pedidos"), pedidosLabel, // Nuevo ítem al final
	}

	// Crear contenedor con layout personalizado
	formLayout := &FormLayout{
		separator:     4,
		margin:        12,
		minEntryWidth: 150,
	}
	formContainer := container.New(formLayout, formItems...)

	// Botones centrados
	buttons := container.NewHBox(deleteButton, cancelButton, saveButton)
	centeredButtons := container.NewCenter(buttons)

	// Contenido completo del formulario con título centrado y en negrita
	formContent := container.NewVBox(
		container.NewCenter(widget.NewLabelWithStyle("Cliente", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		formContainer,
		centeredButtons,
	)

	// Crear el popup modal
	popup := widget.NewModalPopUp(formContent, w.Canvas())
	popup.Resize(fyne.NewSize(500, 0))

	// Lógica del botón Guardar con validación mínima
	saveButton.OnTapped = func() {
		nombre := strings.TrimSpace(nombreEntry.Text)
		if len(nombre) <= 2 {
			dialog.ShowError(fmt.Errorf("el nombre debe tener más de 2 letras"), w)
			return
		}

		if isEdit {
			client.Name = nombre
			client.Phone = telefonoEntry.Text
			client.Email = emailEntry.Text
			client.CUIT = cuitEntry.Text
			client.Address = direccionEntry.Text
			db.Save(client)
		} else {
			newClient := entities.Client{
				Name:    nombre,
				Phone:   telefonoEntry.Text,
				Email:   emailEntry.Text,
				CUIT:    cuitEntry.Text,
				Address: direccionEntry.Text,
			}
			db.Create(&newClient)
			idLabel.SetText(fmt.Sprintf("%05d", newClient.ID))
		}
		ClientList(db, w, nuevoBoton)
		popup.Hide()
	}

	// Lógica del botón Cancelar
	cancelButton.OnTapped = func() {
		popup.Hide()
	}

	// Lógica del botón Eliminar
	if isEdit {
		deleteButton.OnTapped = func() {
			dialog.ShowConfirm("Confirmar", "¿Eliminar este cliente?", func(confirmed bool) {
				if confirmed {
					db.Delete(client)
					dialog.ShowInformation("Éxito", "Cliente eliminado", w)
					ClientList(db, w, nuevoBoton)
					popup.Hide()
				}
			}, w)
		}
	}

	// Mostrar el popup
	popup.Show()
}
