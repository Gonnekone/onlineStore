package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	Status     string    `json:"status"`
	StatusCode int       `json:"-"`
	Info       ErrorInfo `json:"info"`
}

func ErrorResponse(err Error) Response {
	return CreateResponse(err.StatusCode, err)
}

type Success struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func SuccessResponse(statusCode int, data any) Response {
	success := Success{}
	success.Status = "success"
	success.Data = data
	return CreateResponse(statusCode, success)
}

type Response struct {
	code int
	data any
}

func CreateResponse(code int, data any) Response {
	response := Response{}
	response.code = code
	response.data = data
	return response
}

func SendResponse(c *gin.Context, response Response) {
	c.JSON(response.code, response.data)
}

var RouteNotFound = Error{
	"error",
	http.StatusNotFound,
	ErrorInfo{
		404,
		"Route not found",
	}}

var MethodNotAllowed = Error{
	"error",
	http.StatusMethodNotAllowed,
	ErrorInfo{
		405,
		"Method not allowed",
	}}
