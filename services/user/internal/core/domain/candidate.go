package domain

type PreCandiddate struct {
	Email        string
	UserName     string
	PasswordHash string
	RoleId       string
}
type Candidate struct {
	ID       string
	Email    string
	UserName string
}

func NewCandidate(
	id string,
	email string,
	userName string,
) *Candidate {
	return &Candidate{
		ID:       id,
		Email:    email,
		UserName: userName,
	}
}

func NewPreCandidate(
	email string,
	userName string,
	password string,
	roleId string,
) PreCandiddate {
	return PreCandiddate{
		Email:        email,
		UserName:     userName,
		PasswordHash: password,
		RoleId:       roleId,
	}
}
