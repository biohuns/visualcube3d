package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

	r.Get("/cube", generateCube)

	fmt.Println("listening...")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln("listen error has occurred")
	}
}

func generateCube(w http.ResponseWriter, r *http.Request) {
	req, err := bindGenerateCubeRequest(r.Context(), r.URL.Query())
	if err != nil {
		ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}.Write(w)
		return
	}

	body, err := Generate(req.Algorithm, req.IsBinary, req.IsUnlit)
	if err != nil {
		panic(err)
	}

	if req.IsBinary {
		w.Header().Set("Content-Type", "model/gltf-binary")
	} else {
		w.Header().Set("Content-Type", "model/gltf+json")
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}

type (
	Request struct {
		Algorithm []string
		IsBinary  bool
		IsUnlit   bool
	}

	ErrorResponse struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	}
)

func (r ErrorResponse) Write(w http.ResponseWriter) {
	body, err := json.Marshal(r)
	if err != nil {
		return
	}
	w.WriteHeader(r.StatusCode)
	_, _ = w.Write(body)
}

var faces = [18]string{
	"U", "D", "F", "B", "L", "R",
	"U'", "D'", "F'", "B'", "L'", "R'",
	"U2", "D2", "F2", "B2", "L2", "R2",
}

func bindGenerateCubeRequest(ctx context.Context, values url.Values) (*Request, error) {
	req := new(Request)

	if format, ok := ctx.Value(middleware.URLFormatCtxKey).(string); ok {
		if format != "gltf" && format != "glb" {
			return nil, errors.New("provided format is not supported. supported: .gltf, .glb")
		}
		req.IsBinary = format == "glb"
	}

	if alg := values.Get("alg"); alg != "" {
		algs := strings.Split(alg, " ")
		for _, a := range algs {
			ok := false
			for _, f := range faces {
				if a == f {
					ok = true
				}
			}
			if !ok {
				return nil, errors.New(`alg must only use "U D F B L R U' D' F' B' L' R' U2 D2 F2 B2 L2 R2"`)
			}
		}
		req.Algorithm = algs
	}

	if isUnlitStr := values.Get("is_unlit"); isUnlitStr != "" {
		if isUnlitStr != "true" && isUnlitStr != "false" {
			return nil, errors.New("is_unlit must be true or false")
		}
		req.IsUnlit = isUnlitStr == "true"
	}

	return req, nil
}
