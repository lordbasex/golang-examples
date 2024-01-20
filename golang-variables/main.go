// main.go
package main

import (
	"fmt"
	package1 "golang-variables/pakage1"
	package2 "golang-variables/pakage2"
	package3 "golang-variables/pakage3"
	"golang-variables/setup"
)

func main() {
	// Inicializar las variables en el paquete main
	setup.SetChannel1(0)
	setup.SetChannel2(0)

	// Llamar a las funciones ModifyChannels de los paquetes package1, package2 y package3
	package1.ModifyChannels()
	package2.ModifyChannels()
	package3.ModifyChannels()

	// Verificar que las variables se han modificado
	fmt.Println("Main - Final Channel1:", setup.GetChannel1())
	fmt.Println("Main - Final Channel2:", setup.GetChannel2())
}
