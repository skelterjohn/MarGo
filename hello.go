package main

import (
	"encoding/json"
	"net/http"
)

func init() {
	http.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
		json.NewEncoder(rw).Encode(map[string]string{"hello": req.FormValue("hello")})
	})
}
