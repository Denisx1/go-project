package candidategrpchandler

import (
	"context"
	"fmt"
	userv1 "go_project/contracts/gen/go/user"
	portin "user/internal/features/candidate/application/port/in"
)

func (h *CandidateGrpcHandler) CreateCandidate(
	ctx context.Context,
	req *userv1.CreateRequest,
) (*userv1.CreateResponse, error) {
	command := portin.CreateCandidateCommand{
		Email:    req.Email,
		UserName: req.UserName,
		Password: req.Password,
		RoleName: req.Role,
	}
	candidateId, err := h.CandidateService.CreateCandidate(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to create candidate :%w", err)
	}
	return &userv1.CreateResponse{
		Id: candidateId,
	}, nil

}
