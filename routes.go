package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Files",
		"GET",
		"/files",
		Files,
	},
	Route{
		"Search",
		"GET",
		"/search/{word}",
		Search,
	},
	Route{
		"Upload",
		"POST",
		"/push",
		Upload,
	},
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
}
