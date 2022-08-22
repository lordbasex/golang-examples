package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var apiHostTTS string = "https://tts.iperfex.com"
var apiUserTTS string = "XXXX"
var apiPassTTS string = "XXXXXXXXX"
var cookie string
var path string = " /var/lib/asterisk/sounds/custom/tts/"

type BodyLogin struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ResponseTTS struct {
	Cache bool   `json:"cache"`
	File  string `json:"file"`
}

// login
func login() (bool, error) {
	url := apiHostTTS + "/rest/login"
	method := "POST"

	dataRaw := map[string]interface{}{
		"user": apiUserTTS,
		"pass": apiPassTTS,
	}

	jsonDataRaw, err := json.Marshal(dataRaw)
	if err != nil {
		return false, fmt.Errorf("could not marshal json: %v", err)
	}

	payload := strings.NewReader(string(jsonDataRaw))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("failed to connect with url : %v", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	//log.Println(string(body))

	cookie = res.Header.Get("Set-Cookie")

	var apiRes BodyLogin
	json.Unmarshal(body, &apiRes)

	if apiRes.Message == "Login successful. correct credentials." {
		//LOGIN OK.
		return true, nil
	} else {
		//LOGIN ERROR.
		return false, fmt.Errorf("%v", apiRes.Message)
	}
}

// TTS
func TTS(txt, voice, rate, cookie string) (ResponseTTS, error) {
	var cache bool = false
	var response ResponseTTS
	var err error
	url := apiHostTTS + "/rest/tts"
	method := "POST"

	data := map[string]interface{}{
		"txt":   txt,
		"voice": voice,
		"rate":  rate,
	}

	//hash
	tempHashTxt := md5.Sum([]byte(txt))
	tempHashVoice := md5.Sum([]byte(voice))
	tempHashRate := md5.Sum([]byte(rate))
	tempHashUser := md5.Sum([]byte(apiUserTTS))
	hash := hex.EncodeToString(tempHashTxt[:])
	hash += hex.EncodeToString(tempHashVoice[:])
	hash += hex.EncodeToString(tempHashRate[:])
	hash += hex.EncodeToString(tempHashUser[:])
	tempHashTotal := md5.Sum([]byte(hash))
	newHash := hex.EncodeToString(tempHashTotal[:])
	file := path
	file += "iperfex_" + newHash + ".wav"
	// log.Printf("totalHash: %s", hash)
	// log.Printf("newHash: %s", newHash)
	// log.Printf("file: %s", file)

	if _, err = os.Stat(file); err == nil {
		//TRUE - EXISTE
		cache = true
		response = ResponseTTS{
			Cache: cache,
			File:  file,
		}
		return response, nil
	} else {
		//FALSE - NO EXISTE.
		cache = false
		log.Printf("creando archivo....")

		//login
		ok, err := login()
		if err != nil {
			return response, fmt.Errorf("login: %s\n", err)
		}

		if ok {

			jsonData, err := json.Marshal(data)
			if err != nil {
				return response, fmt.Errorf("could not marshal json: %s\n", err)
			}

			payload := strings.NewReader(string(jsonData))

			client := &http.Client{}
			req, err := http.NewRequest(method, url, payload)

			if err != nil {
				return response, fmt.Errorf("could not connect to endpoint: %s\n", err)
			}
			req.Header.Add("Content-type", "application/json")
			req.Header.Add("Cookie", cookie)

			res, err := client.Do(req)
			if err != nil {
				return response, fmt.Errorf("could not read the response: %s\n", err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return response, fmt.Errorf("could not read the body: %s\n", err)
			}

			//create diretory
			errDIR := os.MkdirAll(path, 0755)
			if errDIR != nil {
				log.Fatal(errDIR)
				return response, fmt.Errorf("could not write directory: %s\n", errDIR)
			}

			//witefile
			err = ioutil.WriteFile(file, body, 0644)
			if err != nil {
				fmt.Println(err)
				return response, fmt.Errorf("could not write file: %s\n", err)
			}

			response = ResponseTTS{
				Cache: cache,
				File:  file,
			}

			return response, nil
		}
	}
	return response, errors.New("unexpected error")
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	txt := "hola como estas?"
	voice := "Paulina"
	rate := "170"
	res, err := TTS(txt, voice, rate, cookie)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Cache: %v | File: %s", res.Cache, res.File)

}
