package router

import (
	"encoding/json"
	"github.com/biluping/go-router/entity"
	"log"
	"net/http"
)

// 响应数据，以json格式
func response(write http.ResponseWriter, statusCode int, msg string, data interface{}) {
	resEntity := &entity.ResEntity{
		Code: statusCode,
		Msg:  msg,
		Data: data,
	}
	jsonData, err := json.Marshal(resEntity)
	if err != nil {
		log.Println(err)
		return
	}
	write.Header().Add("Content-Type", "application/json; charset=utf-8")
	_, err = write.Write(jsonData)
	if err != nil {
		log.Println(err)
	}
}

func responseCode(write http.ResponseWriter, statusCode int, data interface{}) {
	response(write, statusCode, http.StatusText(statusCode), data)
}

func responseOk(write http.ResponseWriter, data interface{}) {
	response(write, http.StatusOK, http.StatusText(http.StatusOK), data)
}

func responseNotFound(write http.ResponseWriter) {
	response(write, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
}

func ResponseBadRequest(write http.ResponseWriter, msg string) {
	if msg == "" {
		msg = http.StatusText(http.StatusBadRequest)
	}
	response(write, http.StatusBadRequest, msg, nil)
}
