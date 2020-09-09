package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/c20820/visualcube3d/statik"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//go:generate statik -src gltf

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("$PORT must be set")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(10 * time.Second))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln("listen error has occurred")
	}
}
