package domain

type Session struct {
	Candidate CandidateSession
	Progress  RegistrayionProgress
}

type RegistrayionProgress struct {
	Step     string
	Progress int64
	Message  string
}

type CandidateSession struct {
	CandidateID string
	Email       string
	UserName    string
	Role        string
}

func NewSession(
	id string,
	email string,
	userName string,
	role string,
	step string,
	progress int64,
	message string,
) Session {
	return Session{
		Candidate: NewCandidateSession(id, email, userName, role),
		Progress:  NewProgress(step, progress, message),
	}
}
func NewCandidateSession(
	id string,
	email string,
	userName string,
	role string,
) CandidateSession {
	return CandidateSession{
		CandidateID: id,
		Email:       email,
		UserName:    userName,
		Role:        role,
	}
}

func NewProgress(
	step string,
	progress int64,
	message string,
) RegistrayionProgress {
	return RegistrayionProgress{
		Step:     step,
		Progress: progress,
		Message:  message,
	}
}
