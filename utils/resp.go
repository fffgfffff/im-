package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Row   interface{} `json:"row,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total"`
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Row:  data,
		Msg:  msg,
	}

	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(ret)
}
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}
func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(ret)
}
func RespListOk(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
func RespListFail(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, -1, data, total)
}
