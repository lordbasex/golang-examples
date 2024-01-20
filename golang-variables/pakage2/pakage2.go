package package2

import (
	"fmt"
	"golang-variables/setup"
)

func ModifyChannels() {
	// Leer y modificar Channel1 y Channel2 desde package2
	fmt.Println("Package2 - Original Channel1:", setup.GetChannel1())
	fmt.Println("Package2 - Original Channel2:", setup.GetChannel2())

	setup.SetChannel1(102)
	setup.SetChannel2(202)

	fmt.Println("Package2 - 	[*] Modified Channel1:", setup.GetChannel1())
	fmt.Println("Package2 - 	[*] Modified Channel2:", setup.GetChannel2())
}
