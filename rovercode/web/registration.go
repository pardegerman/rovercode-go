package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/levigross/grequests"
	"github.com/pardegerman/rovercode-go/rovercode"
)

type webregistration struct {
	Session *websession
	r       rovercode.Rover
}

// Search for a named rover and return its ID
// Will return zero if rover is not found
func (wr *webregistration) SearchID(rovername string) (id int, err error) {
	if nil == wr.Session {
		err = errors.New("cannot search rovers; no session")
		return
	}

	res, err := wr.Session.Get(
		"/mission-control/rovers/",
		map[string]string{
			"name": rovername,
		},
	)
	if nil != err {
		return
	} else if !res.Ok {
		err = errors.New("could not search for rover; " + res.String())
		return
	}

	rovers := make([]rovercode.Rover, 1)
	err = json.Unmarshal(res.Bytes(), &rovers)
	if nil != err {
		return
	}
	if 1 == len(rovers) {
		id = rovers[0].ID
	} else {
		fmt.Printf("Could not retrieve rover %s: %s\n", rovername, res.String())
	}

	return
}

// Get registration data from rovercode-web
func (wr *webregistration) Get(id int) (r *rovercode.Rover, err error) {
	if nil == wr.Session {
		err = errors.New("cannot get rover data; no session")
		return
	}

	res, err := wr.Session.Get(
		"/mission-control/rovers/"+strconv.Itoa(id)+"/",
		nil,
	)
	if nil != err {
		return
	} else if !res.Ok {
		err = errors.New("could not get rover data; " + res.String())
		return
	}

	err = res.JSON(&wr.r)
	return &wr.r, err
}

// Register a new rover to rovercode-web
func (wr *webregistration) Register(rovername string) (r *rovercode.Rover, err error) {
	if nil == wr.Session {
		err = errors.New("cannot register rover; no session")
		return
	}

	ipaddr, err := getip()
	if nil != err {
		return
	}

	res, err := wr.Session.Post(
		"/mission-control/rovers/",
		&grequests.RequestOptions{
			JSON: map[string]string{
				"name":     rovername,
				"local_ip": ipaddr,
			},
		},
	)
	if nil != err {
		return
	} else if !res.Ok {
		err = errors.New("could not register rover: " + res.String())
		return
	}

	err = res.JSON(&wr.r)
	return &wr.r, err
}

// Store updated data for an existing rover registration
func (wr *webregistration) Update() (err error) {
	if nil == wr.Session {
		return errors.New("cannot update rover data; no session")
	}

	if 0 == wr.r.ID {
		return errors.New("cannot update rover data; no data available")
	}

	// Update ip just in case it has changed
	wr.r.IP, err = getip()
	if nil != err {
		return
	}

	res, err := session.Put(
		"/mission-control/rovers/"+strconv.Itoa(wr.r.ID)+"/",
		&grequests.RequestOptions{
			JSON: wr.r,
		},
	)
	if nil != err {
		return
	} else if !res.Ok {
		return errors.New("could not update rover data; " + res.String())
	}

	return
}

// Return our ip address, for now simply do the first non-loopback IPv4 address
func getip() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return
}
