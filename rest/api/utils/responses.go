package utils

import (
	"encoding/json"
	"net/http"
	"vl_sa/models"
)

var errorResponseMap = map[int]func(http.ResponseWriter, error){
	http.StatusBadRequest:          BuildBadRequest,
	http.StatusNotFound:            BuildNotFound,
	http.StatusInternalServerError: BuildInternalServerError,
}

//BuildResponse ..
func BuildResponse(w http.ResponseWriter, responseStatus int, body interface{}) {
	errorResponseBuilder, ok := errorResponseMap[responseStatus]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseStatus)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	} else {
		errorResponseBuilder(w, body.(error))
	}
}

//BuildInternalServerError ..
func BuildInternalServerError(w http.ResponseWriter, err error) {
	buildResponseInternal(w, http.StatusInternalServerError, models.BadRequest{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

//BuildBadRequest ..
func BuildBadRequest(w http.ResponseWriter, err error) {
	buildResponseInternal(w, http.StatusBadRequest, models.BadRequest{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}

//BuildNotFound ..
func BuildNotFound(w http.ResponseWriter, err error) {
	buildResponseInternal(w, http.StatusNotFound, models.NotFound{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	})
}

func buildResponseInternal(w http.ResponseWriter, responseStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	json.NewEncoder(w).Encode(body)
}
