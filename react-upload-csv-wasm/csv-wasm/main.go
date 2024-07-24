package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"strings"
	"syscall/js"
)

// CSVRecord representa un registro en el CSV
type CSVRecord struct {
	Headers []string            `json:"headers"`
	Data    []map[string]string `json:"data"`
	Status  string              `json:"status"`
}

// toLowerCase convierte todas las cabeceras a minúsculas
func toLowerCase(headers []string) []string {
	for i, header := range headers {
		headers[i] = strings.ToLower(header)
	}
	return headers
}

func processCSV(this js.Value, p []js.Value) interface{} {
	data := p[0].String()
	r := csv.NewReader(bytes.NewBufferString(data))

	headers, err := r.Read()
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Convertir las cabeceras a minúsculas
	lowerCaseHeaders := toLowerCase(headers)

	var records []map[string]string
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		row := make(map[string]string)
		for j, header := range lowerCaseHeaders {
			row[header] = record[j]
		}
		records = append(records, row)

		// Limitar a los primeros 10 registros
		if len(records) >= 10 {
			break
		}
	}

	response := CSVRecord{
		Headers: lowerCaseHeaders,
		Data:    records,
		Status:  "200",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": err.Error(),
		})
	}

	log.Print("=== HEADERS ===")
	log.Print(lowerCaseHeaders)

	log.Print("=== JSON ===")
	log.Print(string(jsonResponse))

	return js.ValueOf(string(jsonResponse))
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("processCSV", js.FuncOf(processCSV))
	<-c
}
