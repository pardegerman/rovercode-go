package web

import (
	"fmt"
	"net/url"

	"github.com/levigross/grequests"
)

// PostArgs is the arguments used for a POST to the session
type PostArgs struct {
	Data map[string]string
	JSON map[string]string
}

type websession struct {
	BasePath string
	server   *url.URL
	sess     *grequests.Session
}

// GET a resource from rovercode-web
func (ws *websession) Get(resource string, params map[string]string) (res *grequests.Response, err error) {
	if nil == ws.server || nil == ws.sess {
		err = initws(ws)
		if nil != err {
			return nil, err
		}
	}
	ws.server.Path = resource
	url := ws.server.String()

	ro := grequests.RequestOptions{
		Params: params,
	}

	fmt.Printf("GETTING to %s with\n\tparams: ", ws.server.String())
	fmt.Println(ro.Params)
	fmt.Printf("\tcookies: ")
	fmt.Println(ws.sess.HTTPClient.Jar.Cookies(ws.server))

	res, err = ws.sess.Get(url, &ro)

	fmt.Print("GET response header: ")
	fmt.Println(res.RawResponse.Header)

	// A GET might update the CSRFTOKEN, store that when that happens
	err = ws.setcsrftoken(res)

	return res, err
}

// POST to a resource on rovercode-web
func (ws *websession) Post(resource string, args PostArgs) (res *grequests.Response, err error) {
	if nil == ws.server || nil == ws.sess {
		err = initws(ws)
		if nil != err {
			return nil, err
		}
	}
	ws.server.Path = resource

	ro := grequests.RequestOptions{
		Data: args.Data,
		JSON: args.JSON,
	}

	// Merge data to add csrftoken to the form
	if 0 < len(args.Data) {
		for n, v := range ws.sess.RequestOptions.Data {
			ro.Data[n] = v
		}
	}

	fmt.Printf("POSTING to %s with\n\tdata: ", ws.server.String())
	fmt.Println(ro.Data)
	fmt.Printf("\tcookies: ")
	fmt.Println(ws.sess.HTTPClient.Jar.Cookies(ws.server))

	res, err = ws.sess.Post(ws.server.String(), &ro)

	fmt.Print("POST response header: ")
	fmt.Println(res.RawResponse.Header)

	// A POST might update the CSRFTOKEN, store that when that happens
	err = ws.setcsrftoken(res)

	return res, err
}

// PUT to a resource on rovercode-web
func (ws *websession) Put(resource string, ro *grequests.RequestOptions) (res *grequests.Response, err error) {
	if nil == ws.server || nil == ws.sess {
		err = initws(ws)
		if nil != err {
			return nil, err
		}
	}
	ws.server.Path = resource

	res, err = ws.sess.Put(ws.server.String(), ro)

	// A PUT might update the CSRFTOKEN, store that when that happens
	err = ws.setcsrftoken(res)

	return res, err
}

func initws(ws *websession) (err error) {
	if "" == ws.BasePath {
		ws.BasePath = "https://rovercode.com/"
	}
	if nil == ws.server {
		ws.server, err = url.Parse(ws.BasePath)
		if nil != err {
			return
		}
	}
	if nil == ws.sess {
		ws.sess = grequests.NewSession(nil)
	}
	return
}

func (ws *websession) setcsrftoken(res *grequests.Response) (err error) {
	for _, c := range res.RawResponse.Cookies() {
		fmt.Printf("Storing cookie %s: %s\n", c.Name, c.Value)
		if "csrftoken" == c.Name {
			ws.sess.RequestOptions.Headers = map[string]string{
				"X-CSRFTOKEN": c.Value,
			}
			ws.sess.RequestOptions.Data = map[string]string{
				"csrfmiddlewaretoken": c.Value,
			}
		}
	}

	return nil
}
