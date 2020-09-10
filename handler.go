package main

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const (
	ContentTypeGltf = "model/gltf+json"
	ContentTypeGlb  = "model/gltf-binary"
)

type request struct {
	Algorithm []string
	IsBinary  bool
	IsUnlit   bool
}

var algRegex = regexp.MustCompile("[UDFBLRudlrfb'2 ]*")

func getCubeHandler(w http.ResponseWriter, r *http.Request) {
	req, err := bindGetCubeHandlerRequest(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	data, err := generateCube(
		req.Algorithm,
		req.IsBinary,
		req.IsUnlit,
	)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	}

	if req.IsBinary {
		w.Header().Set("Content-Type", ContentTypeGlb)
	} else {
		w.Header().Set("Content-Type", ContentTypeGltf)
	}

	render.Data(w, r, data)
}

var allowedAlg = [18]string{
	"U", "D", "F", "B", "L", "R",
	"U'", "D'", "F'", "B'", "L'", "R'",
	"U2", "D2", "F2", "B2", "L2", "R2",
}

func bindGetCubeHandlerRequest(r *http.Request) (*request, error) {
	req := new(request)

	if format, ok := r.Context().Value(middleware.URLFormatCtxKey).(string); ok {
		if format != "gltf" && format != "glb" {
			return nil, errors.New("provided format is not supported. supported: .gltf, .glb")
		}
		req.IsBinary = format == "glb"
	}

	urlValues := r.URL.Query()

	if alg := urlValues.Get("alg"); alg != "" {
		if algRegex.Match([]byte(alg)) {
			return nil, errors.New(`alg must be the following pattern: "[UDFBLRudlrfb'2 ]"`)
		}

		algSlice := strings.Split(alg, " ")
		for _, a := range algSlice {
			ok := false
			for _, allowed := range allowedAlg {
				if a == allowed {
					ok = true
				}
			}
			if !ok {
				return nil, errors.New(`alg must only use "U D F B L R U' D' F' B' L' R' U2 D2 F2 B2 L2 R2"`)
			}
		}

		req.Algorithm = algSlice
	}

	if isUnlit := urlValues.Get("is_unlit"); isUnlit != "" {
		if isUnlit != "true" && isUnlit != "false" {
			return nil, errors.New("is_unlit must be true or false")
		}
		req.IsUnlit = isUnlit == "true"
	}

	return req, nil
}
