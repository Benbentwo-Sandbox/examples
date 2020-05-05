// Package main launches the calculator example directly
package main

import (
	"fyne.io/fyne/app"
	"github.com/Benbentwo-Sandbox/examples/ben"
	"github.com/fyne-io/examples/img/icon"
)

func main() {
	app := app.New()
	app.SetIcon(icon.CalculatorBitmap)

	ben.Show(app)
	app.Run()
}
