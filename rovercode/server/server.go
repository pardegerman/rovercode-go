package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/pardegerman/rovercode-go/rovercode"
	"github.com/tj/go-debug"
)

var dbg = debug.Debug("srv")

type server struct {
	Router         *mux.Router
	SocketIOServer *socketio.Server
	Rover          *rovercode.Rover
}

// Serve all routes and websocket of the rover server
func Serve(r *rovercode.Rover) (err error) {
	s := server{
		Router: mux.NewRouter(),
		Rover:  r,
	}
	s.SocketIOServer, err = socketio.NewServer(nil)
	if nil != err {
		return
	}

	// Socket.io websocket handlers
	s.SocketIOServer.On("connection", func(so socketio.Socket) {
		dbg(
			"new socket.io client connected, we have %d clients connected",
			s.SocketIOServer.Count(),
		)
		so.Emit("status", map[string]string{"data": "Connected"})
	})
	s.SocketIOServer.On("disconnection", func(so socketio.Socket) {
		dbg(
			"socket.io client disconnected, we have %d clients connected",
			s.SocketIOServer.Count(),
		)
		so.Emit("status", map[string]string{"data": "Not connected"})
	})
	s.SocketIOServer.On("status", func(so socketio.Socket) {
		dbg("socket.io status")
	})
	s.SocketIOServer.On("error", func(so socketio.Socket) {
		dbg("socket.io error")
	})
	s.Router.Handle("/socket.io/", s.SocketIOServer)

	// API routes
	s.Router.HandleFunc("/api/v1/sendcommand", s.sendcommand).Methods("POST")

	// Default and error handling
	s.Router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name": s.Rover.Name,
		})
	}).Methods("GET")
	s.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		dbg(
			"route for %#v, method %#v was not found, 404",
			req.URL.String(),
			req.Method,
		)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "request not found",
		})
	})

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
		dbg("START_MOTOR called")
		pin := r.FormValue("pin")
		speed, err := strconv.ParseFloat(r.FormValue("speed"), 64)
		if "" == pin || nil != err {
			log.Print("could not decode speed value")
		}
		err = s.Rover.SetPWM(pin, speed)
		if nil != err {
			log.Printf("could not set PWM at pin %s to %f", pin, speed)
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
		}
	case "STOP_MOTOR":
		dbg("STOP_MOTOR called")
	default:
		log.Printf("Undefined command %s called", cmd)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "undefined command",
		})
	}

	return
}
