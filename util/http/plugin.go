package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GinEndpointHandler is an interface for an object that can deal with
// GIN endpoints (i.e. register them with an engine).
type GinEndpointHandler interface {
	RegisterEndpoint(path string, handler gin.HandlerFunc)
}

// HTTPGinPlugin is an interface to allow for GIN HTTP handlers to be
// used as plugins.
type HTTPGinPlugin interface {
	Handle(*gin.Context)
	Path(path string) string
}

// HandlerFuncWrapper wraps a handlerFunc in an HTTPGinPlugin object
// interface -- it adapts handler functions to the interface.
type HandlerFuncWrapper struct {
	handlerFunc        gin.HandlerFunc
	endpoint           string
	ginEndpointHandler GinEndpointHandler
}

// Name returns the name of the endpoint.
func (wrapper *HandlerFuncWrapper) Name() string {
	return fmt.Sprintf("HTTP handler: %s", wrapper.endpoint)
}

// Startup starts the endpoint --and registers it with the
// GinEndpointHandler.
func (wrapper *HandlerFuncWrapper) Startup() error {
	wrapper.ginEndpointHandler.RegisterEndpoint(wrapper.endpoint, wrapper.handlerFunc)
	return nil
}

// Shutdown shuts the endpoint down (no-op)
func (wrapper *HandlerFuncWrapper) Shutdown() error {
	return nil
}

// Handle handles the request by calling the wrapped function.
func (wrapper *HandlerFuncWrapper) Handle(ctx *gin.Context) {
	wrapper.handlerFunc(ctx)
}

// Path returns the path of the endpoint.
func (wrapper *HandlerFuncWrapper) Path(path string) string {
	return wrapper.endpoint
}

// NewHTTPGinPlugin returns a HandlerFuncWrapper (which is a plugin),
// that delegates gin requests at the path to the given endpoint.
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
