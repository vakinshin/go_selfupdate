package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Simple Fyne App")
	w.Resize(fyne.NewSize(420, 240))

	label := widget.NewLabel("Hello from Go + Fyne!")
	button := widget.NewButton("Click me", func() {
		label.SetText("Button clicked")
	})

	content := container.NewVBox(
		widget.NewLabel("Minimal desktop app"),
		label,
		button,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
