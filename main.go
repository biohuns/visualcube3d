package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/c20820/visualcube3d/statik"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//go:generate statik -src gltf

var (
	port = os.Getenv("PORT")
)

func main() {
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

	r.Get("/cube", getCubeHandler)

	fmt.Println("listening...")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln("listen error has occurred")
	}
}
