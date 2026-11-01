package portin

import (
	"context"
)

type CreateCandidateCommand struct {
	Email    string
	UserName string
	Password string
	RoleName string
}
type CandidateService interface {
	CreateCandidate(
		ctx context.Context,
		candidate CreateCandidateCommand,
	) (string, error)
}
