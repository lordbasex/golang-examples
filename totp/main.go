package main

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"bufio"
	"bytes"
	"encoding/base32"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"time"
)

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Printf("OTPAuth:      %s\n", key.URL())
	fmt.Println("Writing PNG to qr-code.png....")
	ioutil.WriteFile("qr-code.png", data, 0644)
}

func promptForPasscode() string {
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom parameters
func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

// cada ves que el usuario quiera tener una nueva llave.
func CreateCode(issuer, accountName string) string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		panic(err)
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	// display the QR code to the user.
	display(key, buf.Bytes())

	return key.Secret()
}

func main() {

	key := CreateCode("cnsoluciones", "lord.basex@gmail.com")

	// Now Validate that the user's successfully added the passcode.
	fmt.Println("Validating TOTP...")
	passcode := promptForPasscode()
	
	//key := "DCQHODKJCFUEC3R6QUCHSO6AMMSDYGBU"
	
	valid := totp.Validate(passcode, key)
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
