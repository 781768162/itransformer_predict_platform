package code

import (
	"errors"
	"net/http"
)

type BizError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ToHTTP(err error) (int, BizError) {
	msg := "internal error"
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}
	switch {
	case errors.Is(err, ErrInvalidParam), errors.Is(err, ErrPassword), errors.Is(err, ErrJwtExpire):
		return http.StatusBadRequest, BizError{Status: http.StatusBadRequest, Message: msg}
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, BizError{Status: http.StatusNotFound, Message: msg}
	case errors.Is(err, ErrUserNameExists):
		return http.StatusConflict, BizError{Status: http.StatusConflict, Message: msg}
	default:
		return http.StatusInternalServerError, BizError{Status: http.StatusInternalServerError, Message: msg}
	}
}
