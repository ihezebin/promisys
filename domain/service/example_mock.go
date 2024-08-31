package service

import (
	"time"

	"github.com/ihezebin/jwt"
	"github.com/ihezebin/promisys/component/constant"
	"github.com/ihezebin/promisys/domain/entity"
	"github.com/ihezebin/promisys/domain/repository"
)

type exampleDomainServiceMock struct {
	exampleRepository repository.ExampleRepository
}

func (e *exampleDomainServiceMock) ValidateExample(example *entity.Example) (bool, string) {
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

func (e *exampleDomainServiceMock) GenerateToken(example *entity.Example) (string, error) {
	token := jwt.Default(jwt.WithOwner(example.Id), jwt.WithExpire(time.Hour))
	tokenStr, err := token.Signed(constant.TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func NewExampleServiceMock(exampleRepository repository.ExampleRepository) ExampleDomainService {
	return &exampleDomainServiceMock{
		exampleRepository: exampleRepository,
	}
}

var _ ExampleDomainService = (*exampleDomainServiceMock)(nil)
