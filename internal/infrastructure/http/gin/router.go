package gin

import (
	ginframework "github.com/gin-gonic/gin"
)

func NewRouter() *ginframework.Engine {
	router := ginframework.New()

	router.Use(ginframework.Logger())
	router.Use(ginframework.Recovery())
	router.Use(RequestID())

	router.GET("/health", health)

	return router
}
