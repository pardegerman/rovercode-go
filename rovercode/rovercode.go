package rovercode

import (
	"errors"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all" // Load all platforms supported
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("gpio")

// Rover describe all properties of a registered rover
type Rover struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Owner int    `json:"owner"`
	IP    string `json:"local_ip"`

	// Assigned pins
	LeftForward   string `json:"left_forward_pin"`
	LeftBackward  string `json:"left_backward_pin"`
	RightForward  string `json:"right_forward_pin"`
	RightBackward string `json:"right_backward_pin"`
	LeftEye       string `json:"left_eye_pin"`
	RightEye      string `json:"right_eye_pin"`

	pwmpins map[string]embd.PWMPin
}

const (
	cycle = 20 // [us]
)

// SetPWM of the specified pin
func (r *Rover) SetPWM(pinname string, duty float64) (err error) {
	if 0 > duty || 1.0 < duty {
		err = errors.New("duty out of range, must be between 0 and 1")
		return
	}
	pin, found := r.pwmpins[pinname]
	if !found {
		err = r.initpin(pinname)
		if nil != err {
			// Could not init pin, let's mock it
			r.pwmpins[pinname] = nil
		}
	}
	ns := (int)(1e3 * cycle * duty)
	dbg("setting duty on pin %s to %d ns (%f)", pinname, ns, duty)
	if nil != pin {
		err = pin.SetDuty(ns)
	}

	return
}

func (r *Rover) initpin(pin string) (err error) {
	dbg("creating pin %s", pin)
	if nil == r.pwmpins {
		r.pwmpins = make(map[string]embd.PWMPin)
	}
	r.pwmpins[pin], err = embd.NewPWMPin(pin)
	if nil != err {
		dbg("NewPWMPin: %s", err.Error())
		return
	}
	dbg("setting PWM period on pin %s to %d us", pin, cycle)
	err = r.pwmpins[pin].SetMicroseconds(cycle)
	if nil != err {
		dbg("SetMicroseconds: %s", err.Error())
		return
	}

	return
}

func (r *Rover) close() (err error) {
	for _, pin := range r.pwmpins {
		dbg("Closing pin %#v", pin)
		pin.Close()
	}
	embd.CloseGPIO()

	return
}
