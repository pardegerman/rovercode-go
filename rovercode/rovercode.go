package rovercode

import (
	"fmt"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all" // Load all platforms supported
)

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
func (r *Rover) SetPWM(pin string, duty float64) (err error) {
	fmt.Printf("Setting PWM of %s to %f\n", pin, duty)

	if _, ok := r.pwmpins[pin]; !ok {
		fmt.Printf("-- Initializing pin %s\n", pin)
		if 0 == len(r.pwmpins) {
			fmt.Printf("-- Init embd.GPIO\n")
			err = embd.InitGPIO()
			if nil != err {
				return
			}
		}
		r.pwmpins[pin], err = embd.NewPWMPin(pin)
		if nil != err {
			return
		}
		err = r.pwmpins[pin].SetMicroseconds(cycle)
		if nil != err {
			return
		}
	}
	ns := (int)(1e3 * cycle * duty)
	fmt.Printf("-- Setting duty to %d ns (%f) at pin %s\n", ns, duty, pin)
	err = r.pwmpins[pin].SetDuty(ns)

	return
}

func (r *Rover) close() (err error) {
	for _, pin := range r.pwmpins {
		pin.Close()
	}
	embd.CloseGPIO()

	return
}
