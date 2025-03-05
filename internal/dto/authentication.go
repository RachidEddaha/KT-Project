package dto

type Login struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type CreateUserRequest struct {
	Login
}

type JWTTokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
