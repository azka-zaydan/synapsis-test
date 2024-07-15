package service

type AuthService interface {
	Register()
}

type AuthServiceImpl struct {
}

func ProvideAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

func (s *AuthServiceImpl) Register() {

}
