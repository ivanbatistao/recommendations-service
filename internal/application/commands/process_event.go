package commands

import (
	"context"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type ProcessEventCommand struct {
	Event event.Event
}

type ProcessEventHandler struct {
	service *recommendation.Service
}

func NewProcessEventHandler(
	service *recommendation.Service,
) *ProcessEventHandler {
	return &ProcessEventHandler{
		service: service,
	}
}

func (h *ProcessEventHandler) Execute(
	ctx context.Context,
	command ProcessEventCommand,
) error {
	return h.service.ProcessEvent(ctx, command.Event)
}