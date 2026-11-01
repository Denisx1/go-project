package candidateservice

import "user/internal/features/candidate/application/port/out"

type CandidateService struct {
	candidateRepository out.CandidateRepository
	profileRepository   out.ProfileRepository
	outboxRepository    out.OutboxRepository
	roleRepository      out.RoleRepository
	cacheRepository     out.CacheRepository
	txManaget           out.TxManager
	bcrypt              out.Bcrypt
}

func NewCandidateService(
	candidateRepository out.CandidateRepository,
	profileRepository out.ProfileRepository,
	outboxRepository out.OutboxRepository,
	roleRepository out.RoleRepository,
	cacheRepository out.CacheRepository,
	txManager out.TxManager,
	bcrypt out.Bcrypt,
) *CandidateService {
	return &CandidateService{
		candidateRepository: candidateRepository,
		profileRepository:   profileRepository,
		outboxRepository:    outboxRepository,
		roleRepository:      roleRepository,
		cacheRepository:     cacheRepository,
		txManaget:           txManager,
		bcrypt:              bcrypt,
	}
}
