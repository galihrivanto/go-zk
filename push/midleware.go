package push

import (
	"log"
	"net/http"
)

// PanicTrap ensure middleware pipe line flow catch panic
func PanicTrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PANIC ", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
