package main

import (
	"flag"
	"fmt"
	"image/color"

	"github.com/lordbasex/golang-examples/go-ui-avatars/avatarbuilderMod"

	"github.com/lordbasex/golang-examples/go-ui-avatars/calc"
)

var colors = []uint32{
	0xDDDDDD, 0x42c58e, 0x5a8de1, 0x785fe0,
}

func main() {
	flag.Parse()
	ab := avatarbuilderMod.NewAvatarBuilder("./SourceHanSansSC-Medium.ttf", &calc.SourceHansSansSCMedium{})

	ab.SetBackgroundColorHex(colors[0])
	ab.SetFrontgroundColor(color.Black)
	ab.SetFontSize(30)
	ab.SetAvatarSize(64, 64)

	svgContent, err := ab.GenerateImageAndSave("FP")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(svgContent)
}
