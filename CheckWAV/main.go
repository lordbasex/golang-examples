package main

import (
	"fmt"
	"math"
	"os"

	"github.com/youpy/go-wav"
)

func read_wav_file(input_file string, number_of_samples uint32) ([]float64, uint16, uint32, uint16) {

	if number_of_samples == 0 {
		number_of_samples = math.MaxInt32
	}

	blockAlign := 2
	file, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}

	reader := wav.NewReader(file)
	wavformat, err_rd := reader.Format()
	if err_rd != nil {
		panic(err_rd)
	}

	if wavformat.AudioFormat != wav.AudioFormatPCM {
		panic("Audio format is invalid ")
	}

	if int(wavformat.BlockAlign) != blockAlign {
		fmt.Println("Block align is invalid ", wavformat.BlockAlign)
	}

	samples, err := reader.ReadSamples(number_of_samples)

	wav_samples := make([]float64, 0)

	for _, curr_sample := range samples {
		wav_samples = append(wav_samples, reader.FloatValue(curr_sample, 0))
	}

	return wav_samples, wavformat.BitsPerSample, wavformat.SampleRate, wavformat.NumChannels
}

func main() {

	input_audio := "demo1.wav"

	audio_samples, bits_per_sample, input_audio_sample_rate, num_channels := read_wav_file(input_audio, 0)

	fmt.Println("num samples ", len(audio_samples)/int(num_channels))
	fmt.Println("bit depth   ", bits_per_sample)
	fmt.Println("sample rate ", input_audio_sample_rate)
	fmt.Println("num channels", num_channels)
}
