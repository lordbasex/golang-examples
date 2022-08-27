package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func Roundf(x float64) float64 {
	return roundFloat(x, 1)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func main() {
  //RAM
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	var ramTotal float64 = Roundf(float64(memory.Total) / 1024 / 1024 / 1024)
	var ramUsed float64 = Roundf(float64(memory.Used) / 1024 / 1024 / 1024)

	fmt.Printf("memory total: %v Gibibytes (GiB)\n", ramTotal)
	fmt.Printf("memory used: %.1f Gibibytes (GiB)\n", ramUsed)
  //RAM
  
	//CPU
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	total := float64(after.Total - before.Total)

	var cpuUsed float64 = float64(after.User-before.User) / total * 100

	fmt.Printf("cpu total: %f %%\n", float64(100-cpuUsed))
	fmt.Printf("cpu user: %f %%\n", float64(after.User-before.User)/total*100)
  //CPU

}
