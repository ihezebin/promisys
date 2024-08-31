package service

import (
	"time"

	"github.com/ihezebin/jwt"
	"github.com/ihezebin/promisys/component/constant"
	"github.com/ihezebin/promisys/domain/entity"
	"github.com/ihezebin/promisys/domain/repository"
)

type exampleDomainServiceImpl struct {
	exampleRepository repository.ExampleRepository
}

func (svc *exampleDomainServiceImpl) ValidateExample(example *entity.Example) (bool, string) {
	if example.Username != "" && !example.ValidateUsernameRule() {
		return false, "账号格式不正确"
	}
	if example.Password != "" && !example.ValidatePasswordRule() {
		return false, "密码格式不正确"
	}
	if example.Email != "" && !example.ValidateEmailRule() {
		return false, "邮箱格式不正确"
	}

	return true, ""
}

func (svc *exampleDomainServiceImpl) GenerateToken(example *entity.Example) (string, error) {
	token := jwt.Default(jwt.WithOwner(example.Id), jwt.WithExpire(time.Hour))
	tokenStr, err := token.Signed(constant.TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func NewExampleServiceImpl() ExampleDomainService {
	return &exampleDomainServiceImpl{
		exampleRepository: repository.GetExampleRepository(),
	}
}

var _ ExampleDomainService = (*exampleDomainServiceImpl)(nil)
