package models

import "time"

type Response struct {
	Data any `json:"data"`

	ErrMsg string `json:"err_msg"`

	Time string `json:"time"`
}

func NewDataResponse(data any) Response {
	return Response{
		Data:   data,
		ErrMsg: "None",
		Time:   time.Now().String(),
	}
}

func NewBadResponse(data any, err error) Response {
	return Response{
		Data:   data,
		ErrMsg: err.Error(),
		Time:   time.Now().String(),
	}
}
