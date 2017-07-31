package webclient

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

	return ws.sess.Get(ws.server.String(), ro)
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
	res, err = ws.sess.Post(ws.server.String(), ro)

	// A POST might update the CSRFTOKEN, store that when that happens
	for _, c := range res.RawResponse.Cookies() {
		if "csrftoken" == c.Name {
			ws.sess.RequestOptions.Headers = map[string]string{
				"X-CSRFTOKEN": c.Value,
			}
			break
		}
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

	return ws.sess.Put(ws.server.String(), ro)
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
