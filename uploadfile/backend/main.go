package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/youpy/go-wav"
)

const MAX_UPLOAD_SIZE = 1024 * 3072 // 3MB

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		//fmt.Println("DONE!")
		return
	}

	//fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("func uploadHandler")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Method not allowed"})
		return
	}

	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
		return
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {

		//log.Print(fileHeader.Filename)

		if fileHeader.Size > MAX_UPLOAD_SIZE {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 3MB in size", fileHeader.Filename)})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "No file sent."})
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
			return
		}

		filetype := http.DetectContentType(buff)
		log.Print("audio/wave")
		if filetype != "audio/wave" && filetype != "audio/wav" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "The provided file format is not allowed. Please upload a WAV"})
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
			return
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
			return
		}

		f, err := os.Create(fmt.Sprintf("./uploads/%v", fileHeader.Filename))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
			return
		}

		defer f.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
			return
		}

		//validate wav
		samples, bits, frequency, channels := read_wav_file(fmt.Sprintf("./uploads/%v", fileHeader.Filename), 0)

		// Validar el tipo de audio
		if bits != 16 {
			errDelete := deleteFile(fmt.Sprintf("./uploads/%v", fileHeader.Filename))
			if !errDelete {
				log.Print("no se pudo eliminar el archivo")
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "The file must be 16-bit PCM format"})
			return

		}
		if channels != 1 {
			errDelete := deleteFile(fmt.Sprintf("./uploads/%v", fileHeader.Filename))
			if !errDelete {
				log.Print("no se pudo eliminar el archivo")
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "The file wav must be single channel"})
			return
		}
		if frequency != 8000 {
			errDelete := deleteFile(fmt.Sprintf("./uploads/%v", fileHeader.Filename))
			if !errDelete {
				log.Print("no se pudo eliminar el archivo")
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "The sample rate must be 8000 Hz"})
			return
		}

		log.Printf("channels: %v", channels)
		log.Printf("frequency: %v", frequency)
		log.Printf("bit: %v", bits)
		log.Printf("samples: %v", len(samples)/int(channels))

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Upload successful"})
}

func read_wav_file(input_file string, number_of_samples uint32) ([]float64, uint16, uint32, uint16) {

	if number_of_samples == 0 {
		number_of_samples = math.MaxInt32
	}

	blockAlign := 2
	file, err := os.Open(input_file)
	if err != nil {
		log.Print(err)
	}

	reader := wav.NewReader(file)

	wavformat, err_rd := reader.Format()
	if err_rd != nil {
		log.Print(err_rd)
	}

	if wavformat.AudioFormat != wav.AudioFormatPCM {
		log.Print("Audio format is invalid ")
	}

	if int(wavformat.BlockAlign) != blockAlign {
		log.Print("Block align is invalid ", wavformat.BlockAlign)
	}

	samples, err := reader.ReadSamples(number_of_samples) // must supply num samples w/o defaults to 2048
	//                                                    // just supply a HUGE number then actual num is returned
	wav_samples := make([]float64, 0)

	for _, curr_sample := range samples {
		wav_samples = append(wav_samples, reader.FloatValue(curr_sample, 0))
	}

	return wav_samples, wavformat.BitsPerSample, wavformat.SampleRate, wavformat.NumChannels
}

func deleteFile(filename string) bool {
	err := os.Remove(filename)
	if err != nil {
		//fmt.Printf("Error deleting file: %v\n", err)
		return false
	}
	//fmt.Printf("File %s deleted successfully\n", filename)
	return true
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)

	handler := cors.Default().Handler(mux)

	log.Println("Servidor iniciado en http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", handler))

}
