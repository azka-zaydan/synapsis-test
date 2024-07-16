package dto

type RegisterDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

func NewJWTResponse(token string) JWTResponse {
	return JWTResponse{
		Token: token,
	}
}
