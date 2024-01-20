package package3

import (
	"fmt"
	"golang-variables/setup"
)

func ModifyChannels() {
	// Leer y modificar Channel1 y Channel2 desde package3
	fmt.Println("Package3 - Original Channel1:", setup.GetChannel1())
	fmt.Println("Package3 - Original Channel2:", setup.GetChannel2())

	setup.SetChannel1(103)
	setup.SetChannel2(203)

	fmt.Println("Package3 - 	[*] Modified Channel1:", setup.GetChannel1())
	fmt.Println("Package3 - 	[*] Modified Channel2:", setup.GetChannel2())
}
