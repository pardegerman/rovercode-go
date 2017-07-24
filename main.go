package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/pardegerman/rovercode-go/webclient"
)

type params struct {
	baseURL   string
	username  string
	password  string
	rovername string
}

func getenv(key, defaultvalue string) (val string, err error) {
	val = os.Getenv(key)
	if len(val) == 0 {
		if len(defaultvalue) == 0 {
			err = errors.New("Value not set and no default")
		}
		val = defaultvalue
	}
	return
}

func loadparams() (p params, err error) {
	godotenv.Load()
	p.baseURL, err = getenv("ROVERCODE_WEB_URL", "https://rovercode.com/")
	p.username, err = getenv("ROVERCODE_WEB_USER_NAME", "")
	p.password, err = getenv("ROVERCODE_WEB_USER_PASS", "")
	p.rovername, err = getenv("ROVER_NAME", "Curiosity")
	return
}

func main() {
	p, err := loadparams()
	if nil != err {
		log.Fatal(err)
	}

	fmt.Println("Loaded parameters:")
	fmt.Printf(
		" Base URL: %s, Username: %s, Password: %s, Rover name: %s\n",
		p.baseURL, p.username, p.password, p.rovername,
	)

	err = webclient.SetServer(p.baseURL)
	if nil != err {
		log.Fatal(err)
	}

	err = webclient.Login(p.username, p.password)
	if nil != err {
		log.Fatal(err)
	}

	err = webclient.RegisterRover(p.rovername)
	if nil != err {
		log.Fatal(err)
	}
}
