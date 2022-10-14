package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type GinEndpointHandler interface {
	RegisterEndpoint(path string, handler gin.HandlerFunc)
}

type HTTPGinPlugin interface {
	Handle(*gin.Context)
	Path(path string) string
}

type HandlerFuncWrapper struct {
	handlerFunc        gin.HandlerFunc
	endpoint           string
	ginEndpointHandler GinEndpointHandler
}

func (wrapper *HandlerFuncWrapper) Name() string {
	return fmt.Sprintf("HTTP handler: %s", wrapper.endpoint)
}

func (wrapper *HandlerFuncWrapper) Startup() error {
	wrapper.ginEndpointHandler.RegisterEndpoint(wrapper.endpoint, wrapper.handlerFunc)
	return nil
}

func (wrapper *HandlerFuncWrapper) Shutdown() error {
	return nil
}

func (wrapper *HandlerFuncWrapper) Handle(ctx *gin.Context) {
	wrapper.handlerFunc(ctx)
}

func (wrapper *HandlerFuncWrapper) Path(path string) string {
	return wrapper.endpoint
}

func NewHTTPGinPlugin(
	path string,
	handlerFunc gin.HandlerFunc,
	ginEndpointHandler GinEndpointHandler) *HandlerFuncWrapper {
	return &HandlerFuncWrapper{
		endpoint:           path,
		handlerFunc:        handlerFunc,
		ginEndpointHandler: ginEndpointHandler,
	}

}
