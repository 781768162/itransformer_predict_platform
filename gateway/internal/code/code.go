package code

import "errors"

var (
	ErrUnknown      = errors.New("unknown error")
	ErrInvalidParam = errors.New("invalid parameter")
	ErrEncrypt      = errors.New("encrypt error")
	ErrPassword     = errors.New("password wrong")
	ErrSaveFile     = errors.New("save file error")
	ErrParseFile    = errors.New("parse file error")
	ErrJwtCreate    = errors.New("create token error")
	ErrJwtExpire    = errors.New("expired token")
	ErrMessageQueue = errors.New("mq error")
	ErrJsonMarshal  = errors.New("json marshal error")

	ErrDatabase       = errors.New("database error")
	ErrUserNameExists = errors.New("username already exists")
	ErrNotFound       = errors.New("not found")
)
