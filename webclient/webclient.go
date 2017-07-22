package webclient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var (
	server *url.URL
	jar, _ = cookiejar.New(nil)
	client = http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}
)

// SetServer the webclient with the correct settings
func SetServer(serverurl string) (err error) {
	server, err = url.Parse(serverurl)
	return
}

/*
type rovercodeRequest http.Request

func (req *rovercodeRequest) addcsrf() (err error) {
	if nil == csrfcookie {
		return errors.New("can not add csrf protection, no cookie received")
	}
	req.Header.Add("X-CSRFTOKEN", csrfcookie.Value)
	return nil
}
*/

func post(path string, data url.Values) (res *http.Response, err error) {
	if nil == server {
		return nil, errors.New("undefined server url")
	}
	server.Path = path
	req, err := http.NewRequest("POST", server.String(), bytes.NewBufferString(data.Encode()))
	if nil != err {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Printf("POST %s (data: %s)\n", server.String(), data.Encode())
	for _, c := range client.Jar.Cookies(server) {
		fmt.Printf(" Cookie %s: %s\n", c.Name, c.Value)
		if "csrftoken" == c.Name {
			req.Header.Add("X-CSRFTOKEN", c.Value)
		}
	}
	return client.Do(req)
}

func get(path string) (res *http.Response, err error) {
	if nil == server {
		return nil, errors.New("undefined server url")
	}
	server.Path = path
	req, err := http.NewRequest("GET", server.String(), nil)
	if nil != err {
		return
	}

	fmt.Printf("GET %s\n", server.String())
	for _, c := range client.Jar.Cookies(server) {
		fmt.Printf(" Cookie %s: %s\n", c.Name, c.Value)
		if "csrftoken" == c.Name {
			req.Header.Add("X-CSRFTOKEN", c.Value)
		}
	}
	return client.Do(req)
}

/*
func storecsrfcookie(cookies []*http.Cookie) bool {
	for _, c := range cookies {
		if "csrftoken" == c.Name {
			csrfcookie = c
			return true
		}
	}
	return false
}
*/

// RegisterRover to the rovercode webserver
func RegisterRover(rovername string) (err error) {
	return
}

// Login the user to the rovercode web server
func Login(username, password string) (err error) {
	var res *http.Response

	res, err = get("accounts/login/")
	if nil != err {
		return
	}
	defer res.Body.Close()
	//storecsrfcookie(res.Cookies())

	form := url.Values{}
	form.Add("login", username)
	form.Add("password", password)
	for _, c := range res.Cookies() {
		if "csrftoken" == c.Name {
			form.Add("csrfmiddlewaretoken", c.Value)
			break
		}
	}

	res, err = post("accounts/login/", form)
	if nil != err {
		return
	}
	// storecsrfcookie(res.Cookies())

	res, err = get("/mission-control/rovers/")
	if nil != err {
		return
	}
	fmt.Println(res)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Print("Body: ")
	fmt.Println(string(body))

	return
}
