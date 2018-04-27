package gosharexserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

type reverseProxyRouter struct {
	realRouter         *mux.Router
	reverseProxyHeader string
}

// ServeHTTP is the implementation of the http.Handler function which modifies the request to adjust the remote address.
func (reverseProxyRouter *reverseProxyRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// set remote address
	req.RemoteAddr = req.Header[reverseProxyRouter.reverseProxyHeader][0]
	reverseProxyRouter.realRouter.ServeHTTP(writer, req)
}

// WrapRouterToReverseProxyRouter wraps the given router to a one which adjusts incoming requests fitting to the reverse
// proxy real header ip settings.
func WrapRouterToReverseProxyRouter(router *mux.Router, reverseProxyHeader string) http.Handler {
	return &reverseProxyRouter{
		realRouter:         router,
		reverseProxyHeader: reverseProxyHeader,
	}
}
