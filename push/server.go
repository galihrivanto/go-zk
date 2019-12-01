package push

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/galihrivanto/runner"
	"github.com/gorilla/mux"
)

// ServerHook define available callbacks for push server
type ServerHook interface {
	// called on initial exchange.return true allowed exchange,
	// thus further operation is allowed, false otherwise
	OnInitialExchange(d Device) bool
}

// ServerOption define push server setting
type ServerOption struct {
	Name    string
	Address string

	// tls setting
	CertFile string
	KeyFile  string
}

// Server is http server which
// implement push service protocol
type Server struct {
	option *ServerOption

	// hook define callback to extend
	// push service function
	hook ServerHook

	// flag to ensure server started
	// before issuing command to device
	started bool
}

func (s *Server) handleExchange(w http.ResponseWriter, r *http.Request) {
	// parse device info
	var device Device
	if err := Unmarshall([]byte(r.URL.RawQuery), &device); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// check whether device is allowed to be connected
	allowed := true

	// call hook
	if s.hook != nil {
		allowed = s.hook.OnInitialExchange(device)
	}

	if !allowed {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// send exchange result

}

func (s *Server) registerAPI(router *mux.Router) {
	router.HandleFunc("/iclock/cdata", s.handleExchange).Methods("GET")
}

// Start run http server which
// server api endpoints that implement
// push protocol
// this method is blocking
func (s *Server) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)

	router := mux.NewRouter()
	s.registerAPI(router)

	srv := http.Server{
		Addr:    s.option.Address,
		Handler: router,
	}

	enableTLS := s.option.CertFile != "" && s.option.KeyFile != ""

	// enable TLS
	if enableTLS {
		srv.TLSConfig = &tls.Config{
			ServerName: s.option.Name,
		}
	}

	runner.
		Run(runCtx, func(ctx context.Context) error {
			s.started = true

			if enableTLS {
				log.Println("Starting HTTPS service")
				log.Printf("HTTPS service is started on %s\n", s.option.Address)

				return srv.ListenAndServeTLS(s.option.CertFile, s.option.KeyFile)
			}

			log.Println("Starting HTTP service")
			log.Printf("HTTP service is started on %s\n", s.option.Address)

			return srv.ListenAndServe()
		}).
		Handle(func(sig os.Signal) {
			if sig == syscall.SIGHUP {
				return
			}

			log.Println("Shutting down...")
			cancel()
		})

}

// NewServer create and start push service
func NewServer(option *ServerOption, hook ...ServerHook) *Server {
	var h ServerHook
	if len(hook) > 0 {
		h = hook[0]
	}

	s := &Server{option: option, hook: h}

	return s
}
