package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/pquerna/otp"
)

// Convierte un byte array a una cadena base64
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Genera un código QR basado en una clave TOTP existente y lo devuelve en formato base64
func GenerateQRCodeFromKey(secret, issuer, accountName string) string {
	// Generar la URL OTPAuth con la clave TOTP existente
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, accountName, secret, issuer))
	if err != nil {
		panic(err)
	}

	// Convertir la clave TOTP a una imagen PNG (QR)
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	err = png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}

	// Convertir la imagen PNG a una cadena base64
	base64Encoding := "data:image/png;base64,"
	base64Encoding += toBase64(buf.Bytes())

	return base64Encoding
}

func main() {
	// Ejemplo de uso con una clave TOTP ya existente
	secret := "ARR6LQ6EW2FVD5CI47MJIDT43HLGETEM" // Esta es la clave secreta que ya tienes
	issuer := "MyApp"                            // Nombre de la app o empresa que genera el código
	accountName := "user@example.com"            // Nombre de cuenta del usuario

	// Generar el código QR en formato base64
	base64QRCode := GenerateQRCodeFromKey(secret, issuer, accountName)

	// Mostrar el código QR en base64
	fmt.Println(base64QRCode)
}
