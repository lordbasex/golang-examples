package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
)

// StatusHandler maneja la visualización de estados con spinner
type StatusHandler struct {
	currentSpinner *spinner.Spinner
	statusColors   map[string]color.Color
}

// NewStatusHandler crea una nueva instancia de StatusHandler
func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		statusColors: map[string]color.Color{
			"Progress": color.Yellow,
			"Ring":     color.Cyan,
			"Answered": color.Green,
		},
	}
}

// UpdateStatus actualiza el estado actual con un spinner o checkmark
func (sh *StatusHandler) UpdateStatus(status string) {
	// Si hay un spinner activo, detenerlo antes de crear uno nuevo
	if sh.currentSpinner != nil {
		sh.currentSpinner.Stop()
	}

	if status == "Progress" || status == "Ring" {
		// Crear nuevo spinner para el estado actual
		sh.currentSpinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		sh.currentSpinner.Prefix = "\r "
		sh.currentSpinner.Suffix = fmt.Sprintf(" Status: %s   ", sh.statusColors[status].Sprint(status))
		sh.currentSpinner.FinalMSG = ""
		sh.currentSpinner.Color("green")
		sh.currentSpinner.Start()
	} else if status == "Answered" {
		sh.currentSpinner = nil
		fmt.Printf("\r\033[K %s Status: %s   ", color.Green.Sprint("✔"), sh.statusColors[status].Sprint(status))
	}
}

func main() {
	// Imprimir el título en color azul
	color.Blue.Println("[+] Starting test call")

	// Crear instancia del manejador de estados
	statusHandler := NewStatusHandler()

	// Secuencia de estados
	statusHandler.UpdateStatus("Progress")
	time.Sleep(2 * time.Second)
	statusHandler.UpdateStatus("Ring")
	time.Sleep(4 * time.Second)
	statusHandler.UpdateStatus("Answered")
	fmt.Println() // Solo un salto de línea al final
}
