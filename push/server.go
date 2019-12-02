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
	GetMiddlewares(*ServerOption) []Middleware
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
}

func (s *Server) registerAPI(router *mux.Router) {
	// apply custom middleware if defined
	var mws = make([]Middleware, 0)
	if s.hook != nil {
		if hook, ok := s.hook.(MiddlewareProvider); ok {
			mws = append(mws, hook.GetMiddlewares(s.option)...)
		}
	}

	router.Handle("/iclock/cdata", DecorateHandler(http.HandlerFunc(s.handleExchange), mws...)).Methods("GET")
	router.Handle("/iclock/getrequest", DecorateHandler(http.HandlerFunc(s.handleCommand), mws...)).Methods("GET")
	router.Handle("/iclock/cdata", DecorateHandler(http.HandlerFunc(s.handleCommandResponse), mws...)).Methods("POST")

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

// NewServer create and start push service
func NewServer(option *ServerOption, hook ...ServerHook) *Server {
	var h ServerHook
	if len(hook) > 0 {
		h = hook[0]
	}

	s := &Server{option: option, hook: h}

	return s
}
