package event

import (
	"context"
	"fmt"
	"go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/router"
)

type CandidateService interface{}

type EventHandler struct {
	candidateService CandidateService
}

func NewCandidateEventHandler(candidateService CandidateService) *EventHandler {

	return &EventHandler{candidateService: candidateService}
}

func (h *EventHandler) Register(r *router.Router) {
	r.Handle("candidate.created", h.onCreated)
	r.Handle("candidate.delete", h.onDeleted)
}
func (h *EventHandler) onCreated(ctx context.Context, msg messaging.Message) error {
	fmt.Println(msg)
	return nil
}

func (h *EventHandler) onDeleted(ctx context.Context, msg messaging.Message) error {
	fmt.Println(msg)
	return nil
}
