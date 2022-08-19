package main

import (
	"fmt"
	"log"
	"os/exec"
)

func AsteriskCMD(command string) (string, error) {
	prg := "asterisk"

	arg1 := "-rx"
	arg2 := command

	cmd := exec.Command(prg, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(stdout), nil
}

func main() {
	out, err := AsteriskCMD("dialplan reload")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("OUT: %s", out)
}
