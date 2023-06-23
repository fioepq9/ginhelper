package main

import (
	"encoding/json"
)

type ResponseCode int

const (
	ResponseCodeSuccess ResponseCode = 200000 + iota
)

const (
	ResponseCodeBadRequest ResponseCode = 400000 + iota
	ResponseCodeUnauthorized
	ResponseCodeForbidden
	ResponseCodeNotFound
)

const (
	ResponseCodeInternalError ResponseCode = 500000 + iota
)

func (rc ResponseCode) String() string {
	switch rc {
	case ResponseCodeSuccess:
		return "success"
	case ResponseCodeBadRequest:
		return "bad request"
	case ResponseCodeUnauthorized:
		return "unauthorized"
	case ResponseCodeForbidden:
		return "forbidden"
	case ResponseCodeNotFound:
		return "not found"
	case ResponseCodeInternalError:
		return "internal error"
	default:
		return "unknown"
	}
}

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
}

func NewResponse(code ResponseCode, data any) Response {
	return Response{
		Code:    code,
		Message: code.String(),
		Data:    data,
	}
}

func (r Response) Error() string {
	buf, _ := json.Marshal(r)
	return string(buf)
}
