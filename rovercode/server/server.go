package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pardegerman/rovercode-go/rovercode"
)

type server struct {
	Router *mux.Router
	Rover  *rovercode.Rover
}

// Serve all routes and websocket of the rover server
func Serve(r *rovercode.Rover) (err error) {
	s := server{
		Router: mux.NewRouter(),
		Rover:  r,
	}

	// API routes
	s.Router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "The rover %s is running its service at this address", s.Rover.Name)
	}).Methods("GET")
	s.Router.HandleFunc("/api/v1/sendcommand", s.sendcommand).Methods("POST")

	srv := &http.Server{
		Handler:      s.Router,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go log.Fatal(srv.ListenAndServe())

	return
}

func (s *server) sendcommand(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("command")
	switch cmd {
	case "START_MOTOR":
		pin := r.FormValue("pin")
		speed, err := strconv.ParseFloat(r.FormValue("speed"), 64)
		if "" == pin || nil != err {
			log.Print("Could not decode command\n")
		}
		s.Rover.SetPWM(pin, speed)
	case "STOP_MOTOR":
		fmt.Println("STOP_MOTOR called")
	default:
		fmt.Printf("Undefined command %s called", cmd)
		w.WriteHeader(http.StatusInternalServerError)
	}

	return
}
