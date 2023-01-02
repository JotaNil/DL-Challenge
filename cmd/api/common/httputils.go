package common

import (
	"github.com/gorilla/mux"
	"net/http"
)

// GetParamFromRequest returns the requested URL param in a string format.
// if paramName is not found un the url returns common.ParamNotFoundError
func GetParamFromRequest(r *http.Request, paramName string) (string, error) {
	params := mux.Vars(r)
	param, found := params[paramName]
	if !found || param == "" {
		return "", ParamNotFoundError
	}

	return param, nil
}
