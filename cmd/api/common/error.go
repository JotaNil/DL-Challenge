package common

import (
	"errors"
	"fmt"
	"net/http"
)

var ParamNotFoundError = fmt.Errorf("requested parameter was not found %w", ErrorBadRequest)

// http use errors

var ErrorNotFound = errors.New("not found")
var ErrorBadRequest = errors.New("bad request")
var ErrorInternalServer = errors.New("internal server error")

func HandlerErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrorNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, ErrorBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, ErrorInternalServer):
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
