package web

import (
	"net/url"

	"github.com/levigross/grequests"
)

type websession struct {
	BasePath string
	server   *url.URL
	sess     *grequests.Session
}

// GET a resource from rovercode-web
func (ws *websession) Get(resource string, ro *grequests.RequestOptions) (res *grequests.Response, err error) {
	if nil == ws.server || nil == ws.sess {
		err = initws(ws)
		if nil != err {
			return nil, err
		}
	}
	ws.server.Path = resource
	url := ws.server.String()

	err = ws.addcsrftoken(ro)
	if nil != err {
		return nil, err
	}

	dbg("GETTING from %s", ws.server.String())
	if nil != ro {
		dbg(" data: %#v", ro.Data)
		dbg(" header: %#v", ro.Headers)
	}
	dbg(" cookies: %#v", ws.sess.HTTPClient.Jar.Cookies(ws.server))

	res, err = ws.sess.Get(url, ro)

	if nil != res.RawResponse {
		dbg("RESPONSE header: %#v", res.RawResponse.Header)
	}

	return res, err
}

// POST to a resource on rovercode-web
func (ws *websession) Post(resource string, ro *grequests.RequestOptions) (res *grequests.Response, err error) {
	if nil == ws.server || nil == ws.sess {
		err = initws(ws)
		if nil != err {
			return nil, err
		}
	}
	ws.server.Path = resource

	// Merge data to add csrftoken to the form
	err = ws.addcsrftoken(ro)
	if nil != err {
		return nil, err
	}

	dbg("POSTING to %s", ws.server.String())
	if nil != ro {
		dbg(" data: %#v", ro.Data)
		dbg(" header: %#v", ro.Headers)
	}
	dbg(" cookies: %#v", ws.sess.HTTPClient.Jar.Cookies(ws.server))

	res, err = ws.sess.Post(ws.server.String(), ro)

	if nil != res.RawResponse {
		dbg("RESPONSE header: %#v", res.RawResponse.Header)
	}

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

	err = ws.addcsrftoken(ro)
	if nil != err {
		return nil, err
	}

	dbg("PUTTING from %s", ws.server.String())
	if nil != ro {
		dbg(" data: %#v", ro.Data)
		dbg(" header: %#v", ro.Headers)
	}
	dbg(" cookies: %#v", ws.sess.HTTPClient.Jar.Cookies(ws.server))

	res, err = ws.sess.Put(ws.server.String(), ro)

	if nil != res.RawResponse {
		dbg("RESPONSE header: %#v", res.RawResponse.Header)
	}

	return res, err
}

// HasSessionID returns true if a logged in session has been established
func (ws *websession) HasSessionID() bool {
	for _, c := range ws.sess.HTTPClient.Jar.Cookies(ws.server) {
		if "sessionid" == c.Name && "" != c.Value {
			return true
		}
	}
	return false
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

func (ws *websession) addcsrftoken(ro *grequests.RequestOptions) (err error) {
	for _, c := range ws.sess.HTTPClient.Jar.Cookies(ws.server) {
		if "csrftoken" == c.Name {
			if nil == ro {
				ro = &grequests.RequestOptions{}
			}
			ro.Headers = map[string]string{
				"X-CSRFTOKEN": c.Value,
			}
			if 0 < len(ro.Data) {
				ro.Data["csrfmiddlewaretoken"] = c.Value
			}
			break
		}
	}

	return nil
}
