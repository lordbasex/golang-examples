// transcoding.go
// Uso: ./transcoding input.wav output.g729
// Lee un WAV mono/8kHz/16-bit PCM y genera bitstream G.729 crudo (1 frame por cada 10 ms).
// Enlaza dinámicamente contra libbcg729.so (instalada en /usr/local/lib).

package main

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lbcg729 -Wl,-rpath,/usr/local/lib

#include <stdlib.h>
#include <stdint.h>
#include <bcg729/encoder.h>

// Envoltorios mínimos para cgo
// enableVAD: 0 = deshabilitado (siempre voz), 1 = habilitado (detecta silencios)
static bcg729EncoderChannelContextStruct* enc_new() {
    return initBcg729EncoderChannel(0);  // VAD deshabilitado para simplicidad
}

static void enc_close(bcg729EncoderChannelContextStruct* ctx) {
    closeBcg729EncoderChannel(ctx);
}

// Codifica un frame de 80 muestras (10 ms @ 8kHz).
// Devuelve en outLen la cantidad efectiva de bytes (típicamente 10 para voz o 2 para SID).
static void enc_frame(bcg729EncoderChannelContextStruct* ctx,
                      const int16_t *in80,
                      uint8_t *outBytes,
                      uint8_t *outLen) {
    bcg729Encoder(ctx, (int16_t*)in80, outBytes, outLen);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-audio/wav"
	youpywav "github.com/youpy/go-wav"
)

func main() {
	// Verificar si se solicita ayuda o versión (solo si hay al menos 2 argumentos)
	if len(os.Args) >= 2 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help" {
			showHelp()
			os.Exit(0)
		}
		if os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "version" {
			fmt.Println("WAV to G.729 Transcoder v1.0.0")
			os.Exit(0)
		}
	}

	// Verificar argumentos y mostrar ayuda si es necesario
	if len(os.Args) != 3 {
		showHelp()
		os.Exit(2)
	}

	inPath := os.Args[1]
	outPath := os.Args[2]

	if err := run(inPath, outPath); err != nil {
		log.Fatal(err)
	}
}

func showHelp() {
	fmt.Printf(`
🎵 WAV to G.729 Transcoder v1.0.0

DESCRIPTION:
  Converts WAV audio files to G.729 format.
  Optimized for VoIP telephony with excellent compression.

USAGE:
  %s input.wav output.g729

WAV FILE REQUIREMENTS:
  ✅ Format: PCM (AudioFormat = 1)
  ✅ Channels: Mono (1 channel)
  ✅ Sample Rate: 8000 Hz
  ✅ Bits per Sample: 16-bit

EXAMPLES:
  # Basic conversion
  %s audio.wav audio.g729
  
  # With Docker
  docker run --rm -v $PWD:/work wav2g729:latest input.wav output.g729
  
  # Verify conversion (convert back to WAV)
  ffmpeg -f g729 -i output.g729 -ar 8000 -ac 1 -c:a pcm_s16le output.wav

CONVERTING INCOMPATIBLE FILES:
  # Convert any audio to compatible format
  ffmpeg -i input.mp3 -ar 8000 -ac 1 -sample_fmt s16 -acodec pcm_s16le output.wav

FEATURES:
  • G.729 codec (8 kbps, excellent voice quality)
  • Frame size: 10ms (80 samples @ 8kHz)
  • VAD disabled (all frames as voice)
  • Optimized Docker image (~19MB with Alpine)

MORE INFORMATION:
  GitHub: https://github.com/lordbasex/golang-examples/wav2g729
  Documentation: See README.md

`, os.Args[0], os.Args[0])
}

func run(inPath, outPath string) error {
	// Validar y mostrar información del archivo WAV
	if err := validateAndShowWavInfo(inPath); err != nil {
		return err
	}

	// Abrir WAV para conversión
	in, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("opening WAV: %w", err)
	}
	defer in.Close()

	r := youpywav.NewReader(in)

	// Validar formato (validación adicional con youpy/go-wav)
	f, err := r.Format()
	if err != nil {
		return fmt.Errorf("reading WAV format: %w", err)
	}
	if f.AudioFormat != 1 {
		return fmt.Errorf("WAV is not PCM (AudioFormat=%d). Convert to PCM s16le first", f.AudioFormat)
	}
	if f.NumChannels != 1 {
		return fmt.Errorf("mono required (1 channel), channels=%d", f.NumChannels)
	}
	if f.SampleRate != 8000 {
		return fmt.Errorf("8000 Hz required, sampleRate=%d", f.SampleRate)
	}
	if f.BitsPerSample != 16 {
		return fmt.Errorf("16-bit PCM required, bits=%d", f.BitsPerSample)
	}

	// Crear salida
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}
	defer func() {
		_ = out.Sync()
		_ = out.Close()
	}()

	// Inicializar encoder
	ctx := C.enc_new()
	if ctx == nil {
		return errors.New("initBcg729EncoderChannel returned NULL")
	}
	defer C.enc_close(ctx)

	const samplesPerFrame = 80 // 10 ms @ 8kHz
	sampleBuf := make([]int16, 0, samplesPerFrame)

	// Leer y codificar en streaming
	for {
		samples, err := r.ReadSamples(1024) // ~buffer de lectura
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("reading WAV samples: %w", err)
		}

		for _, s := range samples {
			// Canal 0 (es mono, por eso usamos Values[0])
			// Values[0] contiene el valor de la muestra para el primer canal
			v := int16(s.Values[0])
			sampleBuf = append(sampleBuf, v)

			if len(sampleBuf) == samplesPerFrame {
				if err := encodeFrame(ctx, sampleBuf, out); err != nil {
					return err
				}
				sampleBuf = sampleBuf[:0]
			}
		}
	}

	// Si queda cola (<80), la completamos con silencio y codificamos un último frame
	if len(sampleBuf) > 0 {
		padded := make([]int16, samplesPerFrame)
		copy(padded, sampleBuf)
		if err := encodeFrame(ctx, padded, out); err != nil {
			return err
		}
	}

	return nil
}

func encodeFrame(ctx *C.bcg729EncoderChannelContextStruct, pcm80 []int16, w io.Writer) error {
	if len(pcm80) != 80 {
		return fmt.Errorf("encodeFrame: expected 80 samples, got %d", len(pcm80))
	}

	var out [10]byte     // tamaño máximo típico de un frame de voz G.729
	var outLen C.uint8_t // longitud efectiva devuelta por la lib

	// Pasamos punteros a los buffers C
	C.enc_frame(
		ctx,
		(*C.int16_t)(&pcm80[0]),
		(*C.uint8_t)(&out[0]),
		&outLen,
	)

	n := int(outLen)
	if n <= 0 || n > len(out) {
		return fmt.Errorf("encoder returned invalid length: %d", n)
	}

	_, err := w.Write(out[:n])
	return err
}

func validateAndShowWavInfo(filePath string) error {
	// Abrir archivo para análisis
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("couldn't open audio file - %v", err)
	}
	defer f.Close()

	// Crear decoder
	d := wav.NewDecoder(f)

	// Obtener el buffer completo para calcular la duración
	buffer, err := d.FullPCMBuffer()
	if err != nil {
		return fmt.Errorf("couldn't decode audio file - %v", err)
	}

	// Obtener información del archivo
	typeFileTmp := fmt.Sprintf("%v", d)
	substrStart := strings.Index(typeFileTmp, "WAVE")
	var typeFile string
	if substrStart >= 0 {
		typeFile = string([]rune(typeFileTmp)[substrStart : substrStart+4])
	} else {
		typeFile = "UNKNOWN"
	}

	// Calcular la duración en segundos
	totalSamples := len(buffer.Data)
	sampleRate := int(d.SampleRate)
	channels := int(d.NumChans)
	bitDepth := int(d.BitDepth)

	// Calcular duración en segundos
	durationSeconds := float64(totalSamples) / float64(sampleRate*channels)

	// Calcular duración en minutos y segundos para mejor legibilidad
	minutes := int(durationSeconds) / 60
	seconds := int(durationSeconds) % 60
	milliseconds := int((durationSeconds - float64(int(durationSeconds))) * 1000)

	// Mostrar toda la información del archivo
	log.Printf("=== WAV FILE INFORMATION ===")
	log.Printf("File type: %v", typeFile)
	log.Printf("Bit depth: %v bits", bitDepth)
	log.Printf("Sample rate: %v Hz", sampleRate)
	log.Printf("Number of channels: %v", channels)
	log.Printf("Total samples: %v", totalSamples)
	log.Printf("Total duration: %.3f seconds", durationSeconds)
	log.Printf("Formatted duration: %d minutes, %d seconds and %d milliseconds", minutes, seconds, milliseconds)

	// Validar requisitos para G.729
	log.Printf("=== VALIDATION ===")

	// Validar formato PCM (go-audio/wav no expone AudioFormat directamente)
	// Verificamos que sea un archivo WAV válido
	if typeFile != "WAVE" {
		return fmt.Errorf("❌ Invalid file format: %s (required: WAVE). Convert to WAV first", typeFile)
	}
	log.Printf("✅ Audio format: WAVE")

	// Validar canales mono
	if channels != 1 {
		return fmt.Errorf("❌ Invalid channels: %d (required: mono = 1). Convert to mono first", channels)
	}
	log.Printf("✅ Channels: Mono")

	// Validar frecuencia de muestreo
	if sampleRate != 8000 {
		return fmt.Errorf("❌ Invalid sample rate: %d Hz (required: 8000 Hz). Convert to 8000 Hz first", sampleRate)
	}
	log.Printf("✅ Sample rate: 8000 Hz")

	// Validar profundidad de bits
	if bitDepth != 16 {
		return fmt.Errorf("❌ Invalid bit depth: %d bits (required: 16 bits). Convert to 16-bit first", bitDepth)
	}
	log.Printf("✅ Bit depth: 16 bits")

	log.Printf("✅ File is compatible with G.729 encoding")
	log.Printf("📝 Note: Players may show %d seconds (rounded)", int(durationSeconds))

	return nil
}
