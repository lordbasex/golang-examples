package avatarbuilderMod

/*
  We would like to express our gratitude to ShiningRush (https://github.com/ShiningRush) for the font. We borrowed their code and modified it to create an SVG without saving it to disk.
  Mod https://github.com/ShiningRush/avatarbuilder
*/

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

type FontCenterCalculator interface {
	// CalculateCenterLocation used to calculate center location in different font style
	CalculateCenterLocation(string, *AvatarBuilder) (int, int)
}

type AvatarBuilder struct {
	W        int
	H        int
	fontFile string
	fontsize float64
	bg       color.Color
	fg       color.Color
	ctx      *freetype.Context
	calc     FontCenterCalculator
}

func NewAvatarBuilder(fontFile string, calc FontCenterCalculator) *AvatarBuilder {
	ab := &AvatarBuilder{}
	ab.fontFile = fontFile
	ab.bg, ab.fg = color.White, color.Black
	ab.W, ab.H = 200, 200
	ab.fontsize = 95
	ab.calc = calc

	return ab
}

func (ab *AvatarBuilder) SetFrontgroundColor(c color.Color) {
	ab.fg = c
}

func (ab *AvatarBuilder) SetBackgroundColor(c color.Color) {
	ab.bg = c
}

func (ab *AvatarBuilder) SetFrontgroundColorHex(hex uint32) {
	ab.fg = ab.hexToRGBA(hex)
}

func (ab *AvatarBuilder) SetBackgroundColorHex(hex uint32) {
	ab.bg = ab.hexToRGBA(hex)
}

func (ab *AvatarBuilder) SetFontSize(size float64) {
	ab.fontsize = size
}

func (ab *AvatarBuilder) SetAvatarSize(w int, h int) {
	ab.W = w
	ab.H = h
}

func (ab *AvatarBuilder) GenerateImageAndSavePNG(s string, outname string) error {
	bs, err := ab.GenerateImage(s)
	if err != nil {
		return err
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(outname)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	if _, err := b.Write(bs); err != nil {
		return fmt.Errorf("write bytes to file: %w", err)
	}
	if err = b.Flush(); err != nil {
		return fmt.Errorf("flush image: %w", err)
	}

	return nil
}

func (ab *AvatarBuilder) GenerateImageAndSave(s string) (string, error) {
	bs, err := ab.GenerateImage(s)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	buf.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="` + strconv.Itoa(ab.W) + `" height="` + strconv.Itoa(ab.H) + `">`)
	buf.WriteString(`<image width="` + strconv.Itoa(ab.W) + `" height="` + strconv.Itoa(ab.H) + `" xlink:href="data:image/png;base64,` + base64.StdEncoding.EncodeToString(bs) + `" xmlns:xlink="http://www.w3.org/1999/xlink" />`)
	buf.WriteString(`</svg>`)

	return buf.String(), nil
}

func (ab *AvatarBuilder) GenerateImage(s string) ([]byte, error) {
	rgba := ab.buildColorImage()
	if ab.ctx == nil {
		if err := ab.buildDrawContext(rgba); err != nil {
			return nil, err
		}
	}

	x, y := ab.calc.CalculateCenterLocation(s, ab)
	pt := freetype.Pt(x, y)
	if _, err := ab.ctx.DrawString(s, pt); err != nil {
		return nil, fmt.Errorf("draw string: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, rgba); err != nil {
		return nil, fmt.Errorf("png encode: %w", err)
	}

	return buf.Bytes(), nil
}

func (ab *AvatarBuilder) buildColorImage() *image.RGBA {
	bg := image.NewUniform(ab.bg)
	rgba := image.NewRGBA(image.Rect(0, 0, ab.W, ab.H))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	return rgba
}

func (ab *AvatarBuilder) hexToRGBA(h uint32) *color.RGBA {
	rgba := &color.RGBA{
		R: uint8(h >> 16),
		G: uint8((h & 0x00ff00) >> 8),
		B: uint8(h & 0x0000ff),
		A: 255,
	}

	return rgba
}

func (ab *AvatarBuilder) buildDrawContext(rgba *image.RGBA) error {
	// Read the font data.
	fontBytes, err := os.ReadFile(ab.fontFile)
	if err != nil {
		return fmt.Errorf("error when open font file: %w", err)
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return fmt.Errorf("error when parse font file: %w", err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(ab.fontsize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(ab.fg))
	c.SetHinting(font.HintingNone)

	ab.ctx = c
	return nil
}

func (ab *AvatarBuilder) GetFontWidth() int {
	return int(ab.ctx.PointToFixed(ab.fontsize) >> 6)
}
