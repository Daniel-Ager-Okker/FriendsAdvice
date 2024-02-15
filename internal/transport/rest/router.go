package rest

import (
	"FriendsAdvice/internal/transport"

	"github.com/gorilla/mux"
)

type Router struct {
	MuxRouter  *mux.Router
	controller transport.IController
}

// Router register necessary routes and returns an instance of a router.
func CreateRouter(c transport.IController) *Router {
	router := Router{MuxRouter: mux.NewRouter(), controller: c}
	initRouter(&router)
	return &router
}

// Connect spceial REST funcs with right handlers
func initRouter(r *Router) {
	// Checking functionality and readiness
	r.MuxRouter.HandleFunc("/probes/liveness", Liveness).Methods("GET")
	r.MuxRouter.HandleFunc("/probes/readiness", Readiness(r.controller)).Methods("GET")

	// Reading objects from storage
	r.MuxRouter.HandleFunc("/objects/{key}", Put(r.controller)).Methods("PUT")
	r.MuxRouter.HandleFunc("/objects/{key}", Get(r.controller)).Methods("GET")
}
