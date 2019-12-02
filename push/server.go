package push

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/galihrivanto/runner"
	"github.com/gorilla/mux"
)

// ServerHook define available callbacks for push server
type ServerHook interface{}

// ExchangeHook define callback upon device initial exchange
type ExchangeHook interface {
	// called on initial exchange. return command for initiate exchange
	// nil command considered blocked for further communication
	OnInitialExchange(d Device) *ExchangeCommand
}

// Middleware defines the callable that can be chained
type Middleware interface {
	// Handler will execute `next` appropriately after examining the current request
	// and often enhance the `context.Context` object passed to `http.Request`
	// to provide scoped values
	Handler(next http.Handler) http.Handler
}

// MiddlewareFunc provides the adapter for function middleware
type MiddlewareFunc func(http.Handler) http.Handler

// Handler implement Middleware interface
func (m MiddlewareFunc) Handler(next http.Handler) http.Handler {
	return m(next)
}

// MiddlewareProvider define custom middleware pipelines
// which will be executed before actual handler
type MiddlewareProvider interface {
	// GetMiddlewares returns the middleware factories that will be
	// applied for all the endpoints during initialization
	Middlewares(*ServerOption) []Middleware
}

// DecorateHandler returns the new handler by chaining middlewares against
// the handler
func DecorateHandler(handler http.Handler, mw ...Middleware) http.Handler {

	if len(mw) == 0 {
		return handler
	}

	n := len(mw) - 1
	h := handler

	// Call the right most middleware first
	for i := n; i >= 0; i-- {
		h = mw[i].Handler(h)
	}

	return PanicTrap(h)
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

	// deviceCommands map, only allowed devices should be listed here.
	// those device which "sync'ed" on initial exchange
	// commands which target device require check on allowed devices
	// return error device not found if not in list
	deviceCommands sync.Map

	// lisf of in-fligh commands which "sent" to device
	// upon receiveing response, corresponding command callback
	// will be triggered.
	// device command id will be used as identity
	commandCallbacks sync.Map
}

func (s *Server) getCommandQueue(sn string) ([]Command, error) {
	// get command queued to target device
	v, ok := s.deviceCommands.Load(sn)
	if !ok || v == nil {
		return nil, ErrDeviceNotRegistered
	}

	cmds, ok := v.([]Command)
	if !ok {
		return nil, ErrDeviceNotRegistered
	}

	return cmds, nil
}

func (s *Server) flushCommandQueue(sn string) error {
	// get command queued to target device
	v, ok := s.deviceCommands.Load(sn)
	if !ok || v == nil {
		return ErrDeviceNotRegistered
	}

	s.deviceCommands.Store(sn, make([]Command, 0))

	return nil
}

func (s *Server) putCommandQueue(sn string, cmds ...Command) error {
	queue, err := s.getCommandQueue(sn)
	if err != nil {
		return err
	}

	// update queue
	queue = append(queue, cmds...)
	s.deviceCommands.Store(sn, queue)

	return nil
}

func (s *Server) registerCommandCallback(id string, cmd Command) {
	s.commandCallbacks.Store(id, cmd)
}

func (s *Server) getCommandCallback(id string) (Command, error) {
	v, ok := s.commandCallbacks.Load(id)
	if !ok {
		return Command{}, errors.New("Callback not found")
	}

	cmd, ok := v.(Command)
	if !ok {
		return Command{}, errors.New("Callback not found")
	}

	return cmd, nil
}

func (s *Server) removeCommandCallback(sn string) {
	s.commandCallbacks.Delete(sn)
}

// RegisterDevice add device to registered device
// without waiting initial exchange
func (s *Server) RegisterDevice(sn string) {
	s.deviceCommands.Store(sn, make([]Command, 0))
}

// DoBackground send single or multiple command to target device
// no need to wait response from device, as should define on each
// command callback
func (s *Server) DoBackground(target string, cmds ...Command) error {
	// ensure command id is generated
	for _, cmd := range cmds {
		cmd.ID = randomCommandID()

		// put in command queue
		if err := s.putCommandQueue(target, cmd); err != nil {
			return err
		}

		// put in callback list
		s.registerCommandCallback(cmd.ID, cmd)
	}

	return nil
}

// Do execute single command and wait until received response
func (s *Server) Do(target string, cmd Command) (CommandResponse, error) {
	cmd.ID = randomCommandID()

	// put in command queue
	if err := s.putCommandQueue(target, cmd); err != nil {
		return CommandResponse{}, err
	}

	// replace original callback
	waitc := make(chan CommandResponse, 0)
	defer close(waitc)

	go func() {
		cmd.Callback = func(resp CommandResponse) {
			waitc <- resp
		}

		// put in callback list
		s.registerCommandCallback(cmd.ID, cmd)
	}()

	return <-waitc, nil
}

func (s *Server) registerAPI(router *mux.Router) {
	// apply custom middleware if defined
	var mws = make([]Middleware, 0)
	if s.hook != nil {
		if hook, ok := s.hook.(MiddlewareProvider); ok {
			mws = append(mws, hook.Middlewares(s.option)...)
		}
	}

	router.Handle("/iclock/cdata", DecorateHandler(http.HandlerFunc(s.handleExchange), mws...)).
		Methods("GET")

	router.Handle("/iclock/getrequest", DecorateHandler(http.HandlerFunc(s.handleInfo), mws...)).
		Methods("GET").
		Queries("INFO", "{.+}")

	router.Handle("/iclock/getrequest", DecorateHandler(http.HandlerFunc(s.handleCommand), mws...)).
		Methods("GET")

	router.Handle("/iclock/devicecmd", DecorateHandler(http.HandlerFunc(s.handleCommandResponse), mws...)).
		Methods("POST")

	router.Handle("/{path:.*}", DecorateHandler(http.HandlerFunc(s.handleCatchAll), mws...)).Methods("GET", "POST", "PUT", "DELETE")
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
			MinVersion: tls.VersionTLS12,
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

// Ready wait until server ready
func (s *Server) Ready() <-chan struct{} {
	// check if server started
	waitc := make(chan struct{}, 1)

	go func() {
		defer close(waitc)

		for {
			if s.started {
				return
			}

			<-time.After(time.Second)
		}
	}()

	return waitc
}

// NewServer create and start push service
func NewServer(option *ServerOption, hook ...ServerHook) *Server {
	var h ServerHook
	if len(hook) > 0 {
		h = hook[0]
	}

	return &Server{
		option: option,
		hook:   h,
	}
}
