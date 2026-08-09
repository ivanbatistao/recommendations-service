package gin

import (
	ginframework "github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *ginframework.Engine {
	router := ginframework.New()

	router.Use(ginframework.Logger())
	router.Use(ginframework.Recovery())
	router.Use(RequestID())

	router.GET("/health", health)

	router.GET("/recommendations/:userId", handler.GetRecommendations)
	router.POST("/events", handler.ProcessEvent)
	router.POST("/recommendations/generate", handler.GenerateRecommendations)

	return router
}
