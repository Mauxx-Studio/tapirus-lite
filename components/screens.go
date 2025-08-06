package components

import (
	"fmt"
	"tapirus_lite/internal/domain/entities"

	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func ShowMainScreen(db *gorm.DB, w fyne.Window, nuevoBoton *widget.Button) {
	nuevoBoton.SetText("Nuevo Pedido")
	nuevoBoton.OnTapped = func() { NewOrderForm(db, w, nuevoBoton, nil) }
	w.Content().(*fyne.Container).Objects[0] = container.NewScroll(MainSummary(db, w))
	w.Content().Refresh()
}

func MainSummary(db *gorm.DB, w fyne.Window) *fyne.Container {
	now := time.Now()
	morningStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	afternoonStart := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location())
	nextDay := morningStart.Add(24 * time.Hour)

	morningOrders := getOrdersBetween(db, morningStart, afternoonStart)
	afternoonOrders := getOrdersBetween(db, afternoonStart, nextDay)

	morningSummary := summarizeOrders(db, morningOrders)
	afternoonSummary := summarizeOrders(db, afternoonOrders)

	morningLabel := widget.NewLabel("Pedidos de la mañana:")
	morningProducts := widget.NewLabel(morningSummary.products)
	morningNotes := widget.NewLabel(morningSummary.notes)
	afternoonLabel := widget.NewLabel("Pedidos de la tarde:")
	afternoonProducts := widget.NewLabel(afternoonSummary.products)
	afternoonNotes := widget.NewLabel(afternoonSummary.notes)

	return container.NewVBox(
		morningLabel, morningProducts, morningNotes,
		widget.NewSeparator(),
		afternoonLabel, afternoonProducts, afternoonNotes,
	)
}

type summary struct {
	products string
	notes    string
}

func getOrdersBetween(db *gorm.DB, start, end time.Time) []entities.Order {
	var orders []entities.Order
	db.Preload("Items.Product").
		Where("delivery_date >= ? AND delivery_date < ? AND completed = ?", start, end, false).
		Find(&orders)
	return orders
}

func summarizeOrders(db *gorm.DB, orders []entities.Order) summary {
	productTotals := make(map[uint]struct {
		qty  float64 // Cambiado de int a float64
		unit string
	})
	var notes string
	for _, order := range orders {
		for _, item := range order.Items {
			total := productTotals[item.ProductID]
			total.qty += item.Quantity
			total.unit = item.Product.Unit
			productTotals[item.ProductID] = total
		}
		if order.Note != "" {
			notes += fmt.Sprintf("%s - %s\n", order.ClientName, order.Note)
		}
	}

	var products string
	for id, total := range productTotals {
		var p entities.Product
		if err := db.First(&p, id).Error; err != nil {
			products += fmt.Sprintf("Producto ID %d no encontrado: %.2f %s\n", id, total.qty, total.unit)
			continue
		}
		if total.qty == float64(int(total.qty)) {
			products += fmt.Sprintf("%s: %d %s\n", p.Name, int(total.qty), total.unit)
		} else {
			products += fmt.Sprintf("%s: %.2f %s\n", p.Name, total.qty, total.unit)
		}
	}
	return summary{products: products, notes: notes}
}
