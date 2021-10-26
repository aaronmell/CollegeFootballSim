package utils

import "net/http"

type Code struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	Success           = &Code{http.StatusOK, 200, "OK"}
	RequestParamError = &Code{http.StatusBadRequest, 400, "Bad Request Parameters"}
)
