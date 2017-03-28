package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		wrappedHandler := Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(wrappedHandler)
	}

	return router
}

var routes = Routes{
	Route{
		"Player",
		"GET",
		"/player/{player_id}",
		PlayerHandler,
	},
	Route{
		"Leaders",
		"GET",
		"/leaders/{category}",
		LeadersHandler,
	},
}
