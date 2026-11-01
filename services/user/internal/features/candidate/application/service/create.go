package candidateservice

import (
	"context"
	"fmt"
	error_type "go_project/infrastructure/error"
	"user/internal/core/domain"
	portin "user/internal/features/candidate/application/port/in"
)

func (s *CandidateService) CreateCandidate(
	ctx context.Context,
	candidate portin.CreateCandidateCommand,
) (string, error) {

	if err := s.isUnique(ctx, candidate.UserName, candidate.Email); err != nil {
		return "", fmt.Errorf("create candidate: %w", err)
	}

	role, err := s.getRole(ctx, candidate.RoleName)
	if err != nil {
		return "", fmt.Errorf("failed to get candidate role role=%s :%w", candidate.RoleName, err)
	}
	hashedPassword, err := s.hashPassword(ctx, candidate.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hahs password: %w", err)
	}

	candidateDomain := domain.NewPreCandidate(candidate.Email, candidate.UserName, hashedPassword, role.ID)

	newCandidate, err := s.createCandidate(ctx, candidate.RoleName, candidateDomain)
	if err != nil {
		return "", fmt.Errorf("candidate create: %w", err)
	}

	if err := s.createSession(ctx, candidate.RoleName, newCandidate); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	return newCandidate.ID, nil
}

func (s *CandidateService) createSession(
	ctx context.Context,
	role string,
	candidate domain.Candidate,
) error {
	candidateSession := domain.NewSession(
		candidate.ID,
		candidate.Email,
		candidate.UserName,
		role,
		"AWAITING_CONFIRMATION",
		10,
		"waiting to confirm",
	)
	fmt.Println(candidateSession)
	if err := s.cacheRepository.Set(ctx, candidateSession); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil

}

func (s *CandidateService) createCandidate(
	ctx context.Context,
	role string,
	preCandidate domain.PreCandiddate,
) (domain.Candidate, error) {
	var candidate domain.Candidate
	err := s.txManaget.WithTx(ctx, func(ctx context.Context) error {
		newCandidate, err := s.candidateRepository.Create(ctx, preCandidate)
		if err != nil {
			return fmt.Errorf("create candidate: %w", err)
		}
		candidateCreated := domain.NewCandidateCreated(
			newCandidate.ID,
			candidate.Email,
			newCandidate.UserName,
			role)
		if err := s.outboxRepository.Publish(ctx, candidateCreated); err != nil {
			return fmt.Errorf("create outbox: %w", err)
		}
		candidate = newCandidate
		return nil
	})
	if err != nil {
		return domain.Candidate{}, fmt.Errorf("create candidate: %w", err)
	}
	return candidate, nil

}

func (s *CandidateService) getRole(
	ctx context.Context,
	roleName string,
) (domain.Role, error) {
	role, err := s.roleRepository.GetByName(ctx, roleName)
	if err != nil {
		return domain.Role{}, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

func (s *CandidateService) hashPassword(
	ctx context.Context,
	password string,
) (string, error) {
	hashedPassword, err := s.bcrypt.Hash(ctx, password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return hashedPassword, nil
}

func (s *CandidateService) isUnique(
	ctx context.Context,
	userName, email string,
) error {

	existingSession, err := s.cacheRepository.ExistsSession(ctx, userName, email)
	if err != nil {
		return fmt.Errorf("check unique: %w", err)
	}
	if existingSession {
		return error_type.ErrAlreadyExist
	}
	existsCandidate, err := s.candidateRepository.ExistsCandidate(ctx, userName, email)
	if err != nil {
		return fmt.Errorf("check unique: %w", err)
	}

	if existsCandidate {
		return error_type.ErrAlreadyExist
	}

	existingProfile, err := s.profileRepository.ExistsProfile(ctx, userName, email)
	if err != nil {
		return fmt.Errorf("check unique: %w", err)
	}

	if existingProfile {
		return error_type.ErrAlreadyExist

	}
	return nil
}
