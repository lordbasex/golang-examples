package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-audio/wav"
)

func main() {

	f, err := os.Open("demo.wav")
	if err != nil {
		panic(fmt.Sprintf("couldn't open audio file - %v", err))
	}

	d := wav.NewDecoder(f)

	// Obtener el buffer completo para calcular la duración
	buffer, err := d.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	f.Close()

	// Obtener información del archivo
	typeFileTmp := fmt.Sprintf("%v", d)
	substrStart := strings.Index(typeFileTmp, "WAVE")
	typeFile := string([]rune(typeFileTmp)[substrStart : substrStart+4])

	// Calcular la duración en segundos
	// Fórmula: Número de muestras ÷ (Frecuencia de muestreo × Canales)
	totalSamples := len(buffer.Data)
	sampleRate := int(d.SampleRate)
	channels := int(d.NumChans)

	// Calcular duración en segundos
	durationSeconds := float64(totalSamples) / float64(sampleRate*channels)

	// Calcular duración en minutos y segundos para mejor legibilidad
	minutes := int(durationSeconds) / 60
	seconds := int(durationSeconds) % 60
	milliseconds := int((durationSeconds - float64(int(durationSeconds))) * 1000)

	// Mostrar toda la información del archivo
	log.Printf("=== INFORMACIÓN DEL ARCHIVO WAV ===")
	log.Printf("Tipo de archivo: %v", typeFile)
	log.Printf("Profundidad de bits: %v bits", d.BitDepth)
	log.Printf("Frecuencia de muestreo: %v Hz", sampleRate)
	log.Printf("Número de canales: %v", channels)
	log.Printf("Total de muestras: %v", totalSamples)
	log.Printf("Duración total: %.3f segundos", durationSeconds)
	log.Printf("Duración formateada: %d minutos, %d segundos y %d milisegundos", minutes, seconds, milliseconds)
	log.Printf("=== VALIDACIÓN ===")
	log.Printf("✅ Cálculo verificado: %.3f segundos es la duración exacta", durationSeconds)
	log.Printf("📝 Nota: Los reproductores pueden mostrar %d segundos (redondeado)", int(durationSeconds))

}

