package models

import "time"

type Response struct {
	Message string `json:"message"`

	Data any `json:"data"`

	ErrMsg string `json:"err_msg"`

	Time string `json:"time"`
}

// NewDataResponse : Create a new response with data
func NewDataResponse(data any) Response {
	return Response{
		Data:   data,
		ErrMsg: "None",
		Time:   time.Now().String(),
	}
}

// NewDataResponseWithMessage : Craete new response with string message and data
func NewDataResponseWithMessage(data any, msg string) Response {
	return Response{
		Message: msg,
		Data:    data,
		Time:    time.Now().String(),
	}
}

// NewBadResponse : Create a bad response with err
func NewBadResponse(data any, err error) Response {
	return Response{
		Data:   data,
		ErrMsg: err.Error(),
		Time:   time.Now().String(),
	}
}
