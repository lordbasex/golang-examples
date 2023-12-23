package main

/*
	//test socket
	watch -n1 "lsof -i :4573" //macOS
	watch -n1 "netstat -putan | grep '4573'" //linux

	//test server FastAGI
	telnet localhost 4573


*/

import (
	"log"
	"time"

	"serverFastAGI/service/agi"
)

var err error
var agiServer *agi.AGIServer
var countSecond int = 1

func main() {
	// Iniciar servidor AGI
	go StartFastAGI()

	// Esperar a que el servidor AGI se establezca
	time.Sleep(2 * time.Second) // Ajusta este tiempo según sea necesario

	// Detener servidor AGI después de 30 segundos
	go StopFastAGI(agiServer)

	// Bucle principal
	for {
		log.Printf("countSecond: %d", countSecond)
		time.Sleep(1 * time.Second)
		countSecond++
	}
}

func handler(a *agi.AGI) {
	defer a.Close()

	a.Answer()
	err := a.Set("MYVAR", "foo")
	if err != nil {
		log.Printf("failed to set variable MYVAR")
	}
	a.Hangup()
}

func StopFastAGI(agiServer *agi.AGIServer) {
	time.Sleep(30 * time.Second)
	log.Print("StopFastAGI")
	agiServer.Close()
}

func StartFastAGI() {
	log.Print("StartFastAGI")

	// Iniciar servidor AGI
	agiServer, err = agi.ListenMod(":4573", handler)
	if err != nil {
		log.Fatal(err)
	}

}
