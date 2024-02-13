package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"syscall/js"

	"github.com/go-audio/wav"
)

var debugMode bool = false //defauilt false
var quit = make(chan struct{})

func customLog(msg string, a ...interface{}) {
	if debugMode {
		log.Printf("%v%v", msg, fmt.Sprintln(a...))
	}
}

func main() {
	//fmt.SetFlags(fmt.LstdFlags | fmt.Lshortfile)
	registerCallbacks()

	// Esperar hasta que se envíe una señal para salir
	<-quit
}

func registerCallbacks() {
	js.Global().Set("start", js.FuncOf(start))
	js.Global().Set("debug", js.FuncOf(debug))
}

func debug(this js.Value, p []js.Value) interface{} {
	if len(p) > 0 && p[0].Type() == js.TypeBoolean {
		debugMode = p[0].Bool()
	} else {
		debugMode = true
	}

	if debugMode {
		customLog("Debug mode enabled.")
	} else {
		customLog("Debug mode disabled.")
	}

	return nil
}

func start(this js.Value, p []js.Value) interface{} {
	customLog("===== Start =====")
	fileInput := js.Global().Get("document").Call("getElementById", "fileInput")
	files := fileInput.Get("files")

	if files.Length() > 0 {
		file := files.Index(0)
		fileName := file.Get("name").String()
		fileSize := file.Get("size").Int()
		customLog("FileName: ", fileName)
		customLog("FileSize: ", fileSize)

		// Llamar a la función para procesar el contenido del archivo
		readFileContents(file)
	}

	// Devolver nil si no se seleccionó ningún archivo
	return nil
}

func readFileContentsAA(file js.Value) {
	reader := js.Global().Get("FileReader").New()

	// Configurar la función de retorno de llamada onload
	reader.Set("onload", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		// Obtener el array buffer del Blob
		arrayBuffer := reader.Get("result")

		// Crear un Uint8Array desde el array buffer
		uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)

		// Obtener el tamaño del array
		length := uint8Array.Get("length").Int()

		// Crear un slice de bytes y copiar los datos
		data := make([]byte, length)
		js.CopyBytesToGo(data, uint8Array)

		// Procesar el contenido del archivo de manera asíncrona
		go func() {
			// Procesar el contenido del archivo
			ok := GetBytesContentType(data)

			var resultString string

			if ok {
				resultString = "El archivo WAV es válido"
			} else {
				resultString = "No es un WAV válido"
			}

			//fmt.Printf("result: %s", resultString)

			// Llamar a la función JavaScript con el resultado
			js.Global().Call("handleResult", resultString)
		}()

		return nil
	}))

	// Leer el contenido del archivo como ArrayBuffer
	reader.Call("readAsArrayBuffer", file)
}

func readFileContents(file js.Value) bool {
	reader := js.Global().Get("FileReader").New()

	// Configurar la función de retorno de llamada onload
	reader.Set("onload", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		// Obtener el array buffer del Blob
		arrayBuffer := reader.Get("result")

		// Crear un Uint8Array desde el array buffer
		uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)

		// Obtener el tamaño del array
		length := uint8Array.Get("length").Int()

		// Crear un slice de bytes y copiar los datos
		data := make([]byte, length)
		js.CopyBytesToGo(data, uint8Array)

		// Procesar el contenido del archivo de manera asíncrona
		go func() {
			// Procesar el contenido del archivo y obtener el resultado como booleano
			ok := GetBytesContentType(data)

			// Llamar a la función JavaScript con el resultado
			js.Global().Call("handleResult", ok)
		}()

		return nil
	}))

	// Leer el contenido del archivo como ArrayBuffer
	reader.Call("readAsArrayBuffer", file)

	// Devolver el resultado
	return true
}

func GetBytesContentType(data []byte) bool {
	// Para detectar el tipo de contenido, se usan solo los primeros
	// 512 bytes.
	const sniffSize = 512

	// Verificar que hay suficientes bytes para leer
	if len(data) < sniffSize {
		return false
	}

	// Tomar los primeros 512 bytes
	buf := data[:sniffSize]

	// La función que realmente hace el trabajo
	contentType := http.DetectContentType(buf)

	customLog("contentType: ", contentType)

	if contentType == "audio/wave" {
		//wavBytes := []byte{}
		wavBytes := data

		// Crear un lector de bytes
		reader := bytes.NewReader(wavBytes)

		// Decodificar el archivo WAV desde los bytes
		decoder := wav.NewDecoder(reader)
		decoder.ReadInfo()

		// Obtener la información necesaria
		bitDepth := decoder.BitDepth
		sampleRate := decoder.SampleRate
		channels := decoder.NumChans

		// Imprimir la información obtenida
		customLog("BitDepth: ", bitDepth)
		customLog("SampleRate: ", sampleRate)
		customLog("Channels: ", channels)
		customLog("===== End =====")

		if bitDepth == 16 && sampleRate == 8000 && channels == 1 {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}
