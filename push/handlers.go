package push

import (
	"log"
	"net/http"
)

func (s *Server) handleExchange(w http.ResponseWriter, r *http.Request) {
	// parse device info
	var device Device
	if err := Unmarshall([]byte(r.URL.RawQuery), &device); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// check whether device is allowed to be connected
	var cmd *ExchangeCommand

	// call hook
	if s.hook != nil {
		if hook, ok := s.hook.(ExchangeHook); ok {
			cmd = hook.OnInitialExchange(device)
		}
	}

	if cmd == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// send exchange result
	b, err := cmd.Marshall()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	log.Println("get command")
}

func (s *Server) handleCommandResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("received command response")
}

func (s *Server) handleCatchAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Method:", r.Method, r.URL.Path)
	log.Println("not found")

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
