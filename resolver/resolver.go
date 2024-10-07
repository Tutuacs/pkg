package resolver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

func GetBody(r *http.Request, response interface{}) (err error) {

	if r.Body == nil {
		return fmt.Errorf("missing body")
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(response)

	return
}

func GetParam(r *http.Request, name string) (param string) {

	param = r.PathValue(name)
	return
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func WriteResponse(w http.ResponseWriter, status int, result interface{}) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

var Validate = validator.New()
