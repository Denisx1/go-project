package domain

type Profile struct {
	ID              string
	Version         int
	Emain           string
	UserName        string
	Password        string
	IsEmailVerified bool
	IsActive        bool
	TokenVersion    int
	PhoneNumber     *string
	RoleId          string
	RequestId       string
}

func NewProfile(
	id string,
	version int,
	email string,
	userName string,
	password string,
	isEmailVerified bool,
	isActive bool,
	tokenVersion int,
	phoneNumber *string,
	roleId string,
	requestId string,
) *Profile {
	return &Profile{
		ID:              id,
		Version:         version,
		Emain:           email,
		UserName:        userName,
		Password:        password,
		IsEmailVerified: isEmailVerified,
		IsActive:        isActive,
		TokenVersion:    tokenVersion,
		PhoneNumber:     phoneNumber,
		RoleId:          roleId,
		RequestId:       requestId,
	}
}
func NewUserUninitialized(
	email string,
	userName string,
	password string,
	phoneNumber *string,
	roleId string,
	requestId string,
) *Profile {
	return NewProfile(
		UninitializedUUID,
		UninitializedVersion,
		email,
		userName,
		password,
		UninitializedVerirfied,
		UninitializedActive,
		UninitializedVersion,
		phoneNumber,
		roleId,
		requestId,
	)
}
