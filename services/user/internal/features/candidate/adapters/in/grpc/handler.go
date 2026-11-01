package candidategrpchandler

import (
	userv1 "go_project/contracts/gen/go/user"
	portin "user/internal/features/candidate/application/port/in"
)

type CandidateGrpcHandler struct {
	CandidateService portin.CandidateService
	userv1.CandidateServiceServer
}

func NewCandidateGrpcHandler(candidateService portin.CandidateService) *CandidateGrpcHandler {
	return &CandidateGrpcHandler{
		CandidateService: candidateService,
	}
}
