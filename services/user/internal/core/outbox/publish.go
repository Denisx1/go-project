package useroutbox

import (
	"context"
	"errors"
	"fmt"
	orderevents "go_project/contracts/events"
	error_type "go_project/infrastructure/error"
	"go_project/infrastructure/grpc/server/interceptors"
	postgres_pool "go_project/infrastructure/repository/postgres/pool"
	"user/internal/core/domain"

	userevents "go_project/contracts/events"
	usereventsv1 "go_project/contracts/gen/go/user/events"

	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/proto"
)

func (o *UserOutbox) Publish(
	ctx context.Context,
	e domain.Event,
) error {
	requestId := ctx.Value(interceptors.RequestIdKey{}).(string)
	eventType, msg, orderId, err := toContract(e)
	if err != nil {
		return err
	}
	payload, err := sonic.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	query := `
		INSERT INTO users.outbox (aggregate_type, aggregate_id, event_type, payload, request_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = o.pool.Exec(
		ctx,
		query,
		orderevents.AggregateType,
		orderId,
		eventType,
		string(payload),
		requestId,
	)
	if err != nil {
		switch {
		case errors.Is(err, postgres_pool.ErrViolatesUnique):
			// Уникальность нарушена -> это конфликт в домене
			return fmt.Errorf("%w: order already exists: %w", error_type.ErrConflict, err)
		case errors.Is(err, postgres_pool.ErrViolatesNotNull),
			errors.Is(err, postgres_pool.ErrViolatesCheck):
			return fmt.Errorf("%w: invalid order data: %w", error_type.ErrInvalidArgument, err)
		case errors.Is(err, postgres_pool.ErrNoRows):
			return fmt.Errorf("%w: order not found: %w", error_type.ErrNotFound, err)

		default:
			return fmt.Errorf("scan error: %w", err)
		}
	}
	return nil
}

func toContract(e domain.Event) (string, proto.Message, string, error) {
	switch ev := e.(type) {
	case domain.CandidateCreated:
		return userevents.EventTypeCreated, &usereventsv1.CandidateCreated{
			CandidateId: ev.CandidateId,
			Email:       ev.Email,
			Role:        ev.Role,
			ActionType:  ev.ActioType,
		}, ev.CandidateId, nil
	default:
		return "", nil, "", fmt.Errorf("unknown event type %T", e)

	}
}
