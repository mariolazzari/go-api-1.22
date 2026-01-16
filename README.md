# How to create your Golang APIs from now on with v1.22

## Path value

```go
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

	router.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("User ID: " + userId))
	})


    // catch all -> GET
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
```

### Route handlers

```go
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

	// put
	router.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("PUT User ID: " + userId))
	})

	// no verbs -> GET
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("Catch all User ID: " + userId))
	})

	server := http.Server{
		Addr:    a.addr,
		Handler: router,
	}

	log.Printf("Server started on %s", a.addr)

	return server.ListenAndServe()
}
```

### Middlewares

```go
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

	// put
	router.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("PUT User ID: " + userId))
	})

	// no verbs -> GET
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("Catch all User ID: " + userId))
	})

	server := http.Server{
		Addr:    a.addr,
		Handler: AuthMiddleware(RequestLoggerMiddleware(router)),
	}

	log.Printf("Server started on %s", a.addr)

	return server.ListenAndServe()
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method %s - Path %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if user is logged
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := range len(middlewares) {
			next = middlewares[i](next)
		}

		return next.ServeHTTP
	}
}
```

### Prefix

```go
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

	// put
	router.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("PUT User ID: " + userId))
	})

	// no verbs -> GET
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		w.Write([]byte("Catch all User ID: " + userId))
	})

	// prefix
	v1 := http.NewServeMux()
	v1.Handle("/api/v1", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    a.addr,
		Handler: AuthMiddleware(RequestLoggerMiddleware(router)),
	}

	log.Printf("Server started on %s", a.addr)

	return server.ListenAndServe()
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method %s - Path %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if user is logged
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := range len(middlewares) {
			next = middlewares[i](next)
		}

		return next.ServeHTTP
	}
}
```
