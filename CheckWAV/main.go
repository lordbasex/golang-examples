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
        _, err = d.FullPCMBuffer()

        if err != nil {
                panic(err)
        }

        f.Close()

        typeFileTmp := fmt.Sprintf("%v", d)
        substrStart := strings.Index(typeFileTmp, "WAVE")
        typeFile := string([]rune(typeFileTmp)[substrStart : substrStart+4])

        log.Printf("TypeFile: %v", typeFile)
        log.Printf("BitDepth: %v", d.BitDepth)
        log.Printf("SampleRate: %v", d.SampleRate)
        log.Printf("Channels: %v", d.NumChans)

}
