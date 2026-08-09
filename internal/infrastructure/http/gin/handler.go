package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	"github.com/ivanbatistao/recommendations-service/internal/application/dto"
	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type Handler struct {
	getRecommendationsHandler    *queries.GetRecommendationsHandler
	processEventHandler         *commands.ProcessEventHandler
	generateRecommendationsHandler *commands.GenerateRecommendationsHandler
}

func NewHandler(
	getRecommendationsHandler *queries.GetRecommendationsHandler,
	processEventHandler *commands.ProcessEventHandler,
	generateRecommendationsHandler *commands.GenerateRecommendationsHandler,
) *Handler {
	return &Handler{
		getRecommendationsHandler:    getRecommendationsHandler,
		processEventHandler:         processEventHandler,
		generateRecommendationsHandler: generateRecommendationsHandler,
	}
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) GetRecommendations(c *gin.Context) {
	userID := c.Param("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	result, err := h.getRecommendationsHandler.Execute(
		c.Request.Context(),
		queries.GetRecommendationsQuery{
			UserID: userID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dtos := dto.FromDomainSlice(result)

	c.JSON(http.StatusOK, gin.H{
		"recommendations": dtos,
	})
}

func (h *Handler) ProcessEvent(c *gin.Context) {
	var eventDTO dto.EventDTO

	if err := c.ShouldBindJSON(&eventDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	event, err := dto.ToEventDomain(eventDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid event data",
		})
		return
	}

	err = h.processEventHandler.Execute(
		c.Request.Context(),
		commands.ProcessEventCommand{
			Event: event,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "event processed",
	})
}

func (h *Handler) GenerateRecommendations(c *gin.Context) {
	var request struct {
		UserID string          `json:"user_id" binding:"required"`
		Events []dto.EventDTO  `json:"events" binding:"required"`
		Limit  int             `json:"limit" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	events := make([]event.Event, len(request.Events))
	for i, eventDTO := range request.Events {
		e, err := dto.ToEventDomain(eventDTO)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid event data",
			})
			return
		}
		events[i] = e
	}

	result, err := h.generateRecommendationsHandler.Execute(
		c.Request.Context(),
		commands.GenerateRecommendationsCommand{
			UserID: request.UserID,
			Events: events,
			Limit:  request.Limit,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dtos := dto.FromDomainSlice(result)

	c.JSON(http.StatusOK, gin.H{
		"recommendations": dtos,
	})
}
