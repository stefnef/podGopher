package handler

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method string
	Path   string
}

type Handler interface {
	GetRoute() *Route
	Handle(context *gin.Context)
}
