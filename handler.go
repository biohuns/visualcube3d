package main

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-chi/render"
)

const (
	ContentTypeGltf = "model/gltf+json"
)

type request struct {
	Algorithm []string
}

var algRegex = regexp.MustCompile("[UDFBLRudlrfb'2 ]*")

func getCubeHandler(w http.ResponseWriter, r *http.Request) {
	req, err := bindGetCubeHandlerRequest(r.URL.Query())
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	data, err := generateCube(req.Algorithm)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	}

	w.Header().Set("Content-Type", ContentTypeGltf)
	_, _ = w.Write(data)
}

var allowedAlg = [18]string{
	"U", "D", "F", "B", "L", "R",
	"U'", "D'", "F'", "B'", "L'", "R'",
	"U2", "D2", "F2", "B2", "L2", "R2",
}

func bindGetCubeHandlerRequest(urlValues url.Values) (*request, error) {
	req := new(request)

	if alg := urlValues.Get("alg"); alg != "" {
		if !algRegex.Match([]byte(alg)) {
			return nil, errors.New(`alg must be the following pattern: "[UDFBLRudlrfb'2 ]*"`)
		}

		algSlice := strings.Split(alg, " ")
		for _, a := range algSlice {
			ok := false
			upper := strings.ToUpper(a)
			for _, allowed := range allowedAlg {
				if upper == allowed {
					ok = true
				}
			}
			if !ok {
				return nil, errors.New(`alg must only use "U D F B L R U' D' F' B' L' R' U2 D2 F2 B2 L2 R2"`)
			}
		}

		req.Algorithm = algSlice
	}

	return req, nil
}
