package package1

import (
	"fmt"
	"golang-variables/setup"
)

func ModifyChannels() {
	// Leer y modificar Channel1 y Channel2 desde package1
	fmt.Println("Package1 - Original Channel1:", setup.GetChannel1())
	fmt.Println("Package1 - Original Channel2:", setup.GetChannel2())

	setup.SetChannel1(101)
	setup.SetChannel2(201)

	fmt.Println("Package1 - 	[*] Modified Channel1:", setup.GetChannel1())
	fmt.Println("Package1 - 	[*] Modified Channel2:", setup.GetChannel2())
}
