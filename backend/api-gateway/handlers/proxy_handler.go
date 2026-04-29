package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(target, stripPrefix string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		trimmedPath := strings.TrimPrefix(req.URL.Path, stripPrefix)
		if trimmedPath == "" {
			trimmedPath = "/"
		}
		if !strings.HasPrefix(trimmedPath, "/") {
			trimmedPath = "/" + trimmedPath
		}
		req.URL.Path = trimmedPath
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
