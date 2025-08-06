package components

import (
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CustomSelectEntry struct {
	*widget.SelectEntry
}

func NewCustomSelectEntry(options []string) *CustomSelectEntry {
	base := widget.NewSelectEntry(options)
	entry := &CustomSelectEntry{base}
	base.Enable() // Asegurarse de que esté habilitado desde el inicio
	return entry
}

func (c *CustomSelectEntry) ShowDropDown() {
	v := reflect.ValueOf(c.SelectEntry).Elem()
	dropDownField := v.FieldByName("dropDown")
	if !dropDownField.IsValid() || dropDownField.IsNil() {
		return
	}
	showMethod := dropDownField.MethodByName("Show")
	if showMethod.IsValid() {
		showMethod.Call([]reflect.Value{})
	}
}

func (c *CustomSelectEntry) HideDropDown() {
	v := reflect.ValueOf(c.SelectEntry).Elem()
	dropDownField := v.FieldByName("dropDown")
	if !dropDownField.IsValid() || dropDownField.IsNil() {
		return
	}
	hideMethod := dropDownField.MethodByName("Hide")
	if hideMethod.IsValid() {
		hideMethod.Call([]reflect.Value{})
	}
}

//El siguiente widget se basa en un entry para hacer el autoconpletado
/*
package main

import (
        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/widget"
        "fyne.io/fyne/v2/canvas"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/container"
)
*/
// AutocompleteEntry es un widget Entry personalizado con autocompletado.
type AutocompleteEntry struct {
	widget.Entry              // Hereda la funcionalidad de Entry
	suggestions  *widget.List // Lista de sugerencias
	data         []string     // Datos para las sugerencias
	limite       int          // Límite de caracteres para mostrar sugerencias
	window       fyne.Window  // Ventana principal para CanvasOverlay
}

// NewAutocompleteEntry crea un nuevo AutocompleteEntry.
// Recibe el límite de caracteres y la ventana principal como argumentos.
func NewAutocompleteEntry(lim int, w fyne.Window) *AutocompleteEntry {
	entry := &AutocompleteEntry{
		limite: lim,
		window: w,
	}
	entry.ExtendBaseWidget(entry) // Extiende la funcionalidad de Entry

	// Crea la lista de sugerencias.
	entry.suggestions = widget.NewList(
		func() int { // Retorna la cantidad de elementos en la lista.
			return len(entry.data)
		},
		func() fyne.CanvasObject { // Crea un nuevo Label para cada elemento.
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) { // Establece el texto del Label.
			item.(*widget.Label).SetText(entry.data[id])
		},
	)
	entry.suggestions.Hide() // Oculta la lista al inicio.

	// Maneja el evento OnChanged del Entry.
	entry.OnChanged = func(s string) {
		entry.updateSuggestions(s)
	}

	// Maneja el evento OnSelected de la lista de sugerencias.
	entry.suggestions.OnSelected = func(id widget.ListItemID) {
		entry.SetText(entry.data[id]) // Establece el texto del Entry con la sugerencia seleccionada.
		entry.suggestions.Hide()      // Oculta la lista.
	}

	return entry
}

// updateSuggestions actualiza la lista de sugerencias según el texto ingresado.
func (e *AutocompleteEntry) updateSuggestions(s string) {
	// Lógica de búsqueda de sugerencias (simulada).
	e.data = []string{"Opción 1" + s, "Opción 2" + s, "Opción 3" + s}

	// Muestra la lista si se alcanza el límite de caracteres y hay sugerencias.
	if len(s) >= e.limite && len(e.data) > 0 {
		e.suggestions.Refresh()                                               // Actualiza la lista.
		e.window.Canvas().Overlays().Add(e.suggestions)                       // Agrega la lista al CanvasOverlay.
		e.suggestions.Resize(fyne.NewSize(e.Size().Width, 150))               // Ajusta el tamaño de la lista.
		e.suggestions.Move(e.Position().Add(fyne.NewPos(0, e.Size().Height))) // Mueve la lista debajo del Entry.
	} else {
		e.window.Canvas().Overlays().Remove(e.suggestions) // Oculta la lista.
	}
}

/*
func main() {
        a := app.New()
        w := a.NewWindow("Autocomplete Entry")

        entry := NewAutocompleteEntry(2, w) // Crea el AutocompleteEntry con límite 2.

        // Crea un formulario con otros widgets.
        form := container.NewVBox(
                widget.NewLabel("Otros widgets aquí"),
                entry,
                widget.NewLabel("Más widgets"),
        )

        w.SetContent(form) // Establece el contenido de la ventana.
        w.ShowAndRun()      // Muestra la ventana y ejecuta la aplicación.
}
*/
