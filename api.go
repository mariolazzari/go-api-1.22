package main

import (
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	// db...
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr,
	}
}

func (a *APIServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("User ID: " + userId))
	})

	server := http.Server{
		Addr:    a.addr,
		Handler: router,
	}

	log.Printf("Server started on %s", a.addr)

	return server.ListenAndServe()
}
