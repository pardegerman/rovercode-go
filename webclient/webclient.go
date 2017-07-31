package webclient

import (
	"github.com/levigross/grequests"

	"errors"
	"time"
)

var (
	session      websession
	registration webregistration
)

// SetServer the webclient with the correct settings
func SetServer(serverurl string) (err error) {
	session = websession{BasePath: serverurl}
	return
}

// Login the user to the rovercode web server
func Login(username, password string) (err error) {
	var res *grequests.Response

	res, err = session.Get("/accounts/login/", nil)
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not retrieve server landing page")
	}

	// TODO: Can this be made using JSON and the csrftoken in the header?
	form := make(map[string]string)
	form["login"] = username
	form["password"] = password
	for _, c := range res.RawResponse.Cookies() {
		if "csrftoken" == c.Name {
			form["csrfmiddlewaretoken"] = c.Value
			break
		}
	}

	res, err = session.Post("/accounts/login/", &grequests.RequestOptions{Data: form})
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not login")
	}

	registration = webregistration{Session: &session}
	return
}

// RegisterRover to the rovercode webserver
func RegisterRover(rovername string) (err error) {
	id, err := registration.SearchID(rovername)
	if nil != err {
		return
	}

	if 0 != id {
		// Rover is already registered, simply retrieve the data
		err = registration.Get(id)
	} else {
		// Register a new rover to rovercode-web
		err = registration.Register(rovername)
	}
	if nil != err {
		return
	}

	// Check in to rovercode-web every three seconds as a keep alive
	go func(wr *webregistration) (err error) {
		for {
			time.Sleep(3 * time.Second)
			err = wr.Update()
			if nil != err {
				break
			}
		}
		return
	}(&registration)

	return
}
