package webclient

import (
	"net/url"

	"github.com/levigross/grequests"

	"errors"
	"fmt"
)

var (
	server    *url.URL
	sess      = grequests.NewSession(nil)
	csrftoken string
)

// SetServer the webclient with the correct settings
func SetServer(serverurl string) (err error) {
	server, err = url.Parse(serverurl)
	return
}

func post(resource string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	if nil == server {
		return nil, errors.New("undefined server url")
	}
	server.Path = resource

	fmt.Print("POST: ")
	fmt.Println(server.String())
	fmt.Println(ro.Data)

	return sess.Post(server.String(), ro)
}

func get(resource string) (*grequests.Response, error) {
	if nil == server {
		return nil, errors.New("undefined server url")
	}
	server.Path = resource

	return sess.Get(server.String(), nil)
}

// RegisterRover to the rovercode webserver
func RegisterRover(rovername string) (err error) {
	var res *grequests.Response

	for _, c := range sess.HTTPClient.Jar.Cookies(server) {
		fmt.Printf("Cookie: %s\n", c.Name)
	}

	res, err = post(
		"/mission-control/rovers/",
		&grequests.RequestOptions{
			Data: map[string]string{
				"name":     rovername,
				"owner":    "-1",
				"local_ip": "192.168.0.1",
				/*
					"left_forward_pin":   "",
					"left_backward_pin":  "",
					"right_forward_pin":  "",
					"right_backward_pin": "",
					"left_eye_pin":       "",
					"right_eye_pin":      "",
				*/
				"csrfmiddlewaretoken": csrftoken,
			},
		},
	)
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not register rover: " + res.String())
	}

	res, err = get("/mission-control/rovers/")
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not retrieve list of rovers")
	}
	fmt.Println(res.String())

	return
}

// Login the user to the rovercode web server
func Login(username, password string) (err error) {
	var res *grequests.Response

	res, err = get("/accounts/login/")
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not retrieve server landing page")
	}

	form := make(map[string]string)
	form["login"] = username
	form["password"] = password
	for _, c := range res.RawResponse.Cookies() {
		if "csrftoken" == c.Name {
			form["csrfmiddlewaretoken"] = c.Value
			break
		}
	}

	res, err = post("/accounts/login/", &grequests.RequestOptions{Data: form})
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not login")
	}
	for _, c := range res.RawResponse.Cookies() {
		if "csrftoken" == c.Name {
			csrftoken = c.Value
			break
		}
	}

	return
}
