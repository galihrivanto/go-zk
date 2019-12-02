package push

import (
	"io/ioutil"
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

	// put on registered devices, if not exists
	if _, ok := s.deviceCommands.Load(device.SN); !ok {
		s.RegisterDevice(device.SN)
	}

	// send exchange result
	b, err := cmd.Marshall()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// TODO: remove this log
	log.Println(string(b))

	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("get info")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	log.Println("get command")

	sn := r.URL.Query().Get("SN")

	// get command queued to target device
	queue, err := s.getCommandQueue(sn)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}

	log.Println("queue", queue)

	// bulk command write
	w.WriteHeader(http.StatusOK)
	if len(queue) > 0 {
		for _, cmd := range queue {
			log.Println("cmd", cmd.ID)
			if b, err := cmd.Marshal(); err == nil {
				w.Write(b)
			}
		}
	} else {
		w.Write([]byte("OK"))
	}

	// flush command
	s.flushCommandQueue(sn)
}

func (s *Server) handleCommandResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("received command response")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var response CommandResponse
	if err := Unmarshall(body, &response); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// get registered callback
	cmd, err := s.getCommandCallback(response.ID)
	if err == nil {

		// trigger callback
		cmd.Callback(response)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleCatchAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Method:", r.Method, r.URL.Path)
	log.Println("not found")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
