package domain

type Event interface {
}

type CandidateCreated struct {
	CandidateId string
	Email       string
	UserName    string
	Role        string
	ActioType   string
}

func NewCandidateCreated(
	id string,
	email string,
	userName string,
	role string,
) CandidateCreated {
	return CandidateCreated{
		CandidateId: id,
		Email:       email,
		UserName:    userName,
		Role:        role,
		ActioType:   "REGISTRATION",
	}
}
