package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tracing/logging"
	"tracing/tracelib"

	"github.com/gin-gonic/gin"
)

type ApiController1 struct {

}

func NewApiController1() *ApiController1{
	return &ApiController1{}
}

func (a *ApiController1) CallApi2(context *gin.Context) {

	defaultLogger := logging.GetDefaultLogger(context).With("method_name", "CallApi2").With("class_name", "ApiController1")
	defaultLogger.Info("inside api 1", "key", "value")
		
	defaultLogger.Info("call to api 2 started")

	resp, err := tracelib.HTTPClient(context.Request.Context(), "GET", "http://go-api2:8081/pong", nil)
	if err != nil {
		defaultLogger.Error("error occured while calling api 2 :", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}else {
		defaultLogger.Info("call to api 2 finished")
	}


	var response tracelib.Response
	unMarshallErr := json.Unmarshal(resp, &response)
	if unMarshallErr != nil {
		fmt.Print("error occured: ", unMarshallErr)
	}

	defaultLogger.Info("call to api 4 started")
	resp2, err2 := tracelib.HTTPClient(context.Request.Context(), "GET", "http://go-api4:8083/dong", nil)
	if err2 != nil {
		defaultLogger.Error("error occured while calling api 4 :", err2)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}else {
		defaultLogger.Info("call to api 4 finished")
	}

	var response2 tracelib.Response
	json.Unmarshal(resp2, &response2)
	response.Message = response.Message + " : " + response2.Message
	defaultLogger.Info("exiting api 1")
	context.JSON(http.StatusOK, response)
}