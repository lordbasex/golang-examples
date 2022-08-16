package main

import (
    "bytes"
    "golang.org/x/crypto/openpgp"
    "encoding/base64"
    "io/ioutil"
    "log"
    "os"
)

//Create GPG Keys
//gpg --gen-key --homedir .gnupg
//gpg --no-default-keyring --homedir ./.gnupg/ --export-secret-keys > ./.gnupg/secring.gpg
//gpg --no-default-keyring --homedir ./.gnupg/ --export > ./.gnupg/pubring.gpg

const mySecretString = "hola esto es una prueba!"
const prefix, passphrase = "/root/go/src/gpg/", "LordBaseX"
const secretKeyring = prefix + ".gnupg/secring.gpg"
const publicKeyring = prefix + ".gnupg/pubring.gpg"

func encTest(secretString string) (string, error) {
    log.Println("Secret to hide:", secretString)
    log.Println("Public Keyring:", publicKeyring)

    // Read in public key
    keyringFileBuffer, _ := os.Open(publicKeyring)
    defer keyringFileBuffer.Close()
    entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
    if err != nil {
        return "", err
    }

    // encrypt string
    buf := new(bytes.Buffer)
    w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
    if err != nil {
        return "", err
    }
    _, err = w.Write([]byte(mySecretString))
    if err != nil {
        return "", err
    }
    err = w.Close()
    if err != nil {
        return "", err
    }

    // Encode to base64
    bytes, err := ioutil.ReadAll(buf)
    if err != nil {
        return "", err
    }
    encStr := base64.StdEncoding.EncodeToString(bytes)

    // Output encrypted/encoded string
    log.Println("Encrypted Secret:", encStr)

    return encStr, nil
}

func decTest(encString string) (string, error) {

    log.Println("Secret Keyring:", secretKeyring)
    log.Println("Passphrase:", passphrase)

    // init some vars
    var entity *openpgp.Entity
    var entityList openpgp.EntityList

    // Open the private key file
    keyringFileBuffer, err := os.Open(secretKeyring)
    if err != nil {
        return "", err
    }
    defer keyringFileBuffer.Close()
    entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
    if err != nil {
        return "", err
    }
    entity = entityList[0]

    // Get the passphrase and read the private key.
    // Have not touched the encrypted string yet
    passphraseByte := []byte(passphrase)
    log.Println("Decrypting private key using passphrase")
    entity.PrivateKey.Decrypt(passphraseByte)
    for _, subkey := range entity.Subkeys {
        subkey.PrivateKey.Decrypt(passphraseByte)
    }
    log.Println("Finished decrypting private key using passphrase")

    // Decode the base64 string
    dec, err := base64.StdEncoding.DecodeString(encString)
    if err != nil {
        return "", err
    }

    // Decrypt it with the contents of the private key
    md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
    if err != nil {
        return "", err
    }
    bytes, err := ioutil.ReadAll(md.UnverifiedBody)
    if err != nil {
        return "", err
    }
    decStr := string(bytes)

    return decStr, nil
}

func main() {
    encStr, err := encTest(mySecretString)
    if err != nil {
        log.Fatal(err)
    }
    decStr, err := decTest(encStr)
    if err != nil {
        log.Fatal(err)
    }
    // should be done
    log.Println("Decrypted Secret:", decStr)
}
