package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/go-mp3"
	"github.com/youpy/go-wav"
)

var debugMode bool = true // default false
var quit = make(chan struct{})

func customLog(msg string, a ...interface{}) {
	if debugMode {
		log.Printf("%v%v", msg, fmt.Sprintln(a...))
	}
}

func main() {
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
	customLog(">>>>> Start <<<<<")
	fileInput := js.Global().Get("document").Call("getElementById", "fileInput")
	files := fileInput.Get("files")

	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		if files.Length() > 0 {
			file := files.Index(0)
			fileName := file.Get("name").String()
			fileSize := file.Get("size").Int()
			fileType := file.Get("type").String()
			customLog("===== MP3 =====")
			customLog("FileName: ", fileName)
			customLog("FileSize: ", fileSize)
			customLog("FileType: ", fileType)

			// Llamar a la función para procesar el contenido del archivo
			readFileContents(file, resolve, reject)
		} else {
			reject.Invoke("No file selected")
		}
		return nil
	}))

	// Devolver la promesa
	return promise
}

func readFileContents(file js.Value, resolve js.Value, reject js.Value) {
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
			// Convertir el archivo MP3 a WAV
			wavData, err := convertMp3ToWav(data)
			if err != nil {
				log.Printf("Error converting MP3 to WAV: %v", err)
				reject.Invoke("Error converting MP3 to WAV")
				return
			}

			// Convertir el resultado a Uint8Array y resolver la promesa
			uint8Array := js.Global().Get("Uint8Array").New(len(wavData))
			js.CopyBytesToJS(uint8Array, wavData)
			resolve.Invoke(uint8Array)
		}()

		return nil
	}))

	// Leer el contenido del archivo como ArrayBuffer
	reader.Call("readAsArrayBuffer", file)
}

func processWavFile(data []byte) error {
	reader := wav.NewReader(bytes.NewReader(data))

	// Obtener el formato del archivo WAV
	format, err := reader.Format()
	if err != nil {
		return fmt.Errorf("failed to read WAV format: %w", err)
	}

	// Mostrar la información del formato del archivo WAV
	customLog("===== WAV =====")
	customLog("Audio Format: ", format.AudioFormat)
	customLog("Number of Channels: ", format.NumChannels)
	customLog("Sample Rate: ", format.SampleRate)
	customLog("Bits per Sample: ", format.BitsPerSample)
	customLog("===============")

	return nil
}

func convertMp3ToWav(data []byte) ([]byte, error) {
	// Leer el archivo MP3
	reader := bytes.NewReader(data)
	dec, err := mp3.NewDecoder(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create mp3 decoder: %v", err)
	}

	customLog("Length: ", dec.Length())
	customLog("Sample Rate: ", dec.SampleRate())

	// Crear un buffer para almacenar los datos decodificados
	var pcmData []int16
	for {
		var sample int16
		err := binary.Read(dec, binary.LittleEndian, &sample)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read mp3 data: %v", err)
		}
		pcmData = append(pcmData, sample)
	}

	// Crear un archivo WAV con la misma información de audio
	wavFile := &bytes.Buffer{}
	writer := wav.NewWriter(wavFile, uint32(len(pcmData)), 1, 8000, 16)

	// Escribir los datos de audio en el archivo WAV
	for _, sample := range pcmData {
		if err := binary.Write(writer, binary.LittleEndian, sample); err != nil {
			return nil, fmt.Errorf("failed to write wav data: %v", err)
		}
	}

	// Procesar el contenido del archivo WAV
	wavFileData := wavFile.Bytes()
	processWavFile(wavFileData)

	return wavFile.Bytes(), nil
}
