package transfers

import (
	"github.com/daniel1sender/Desafio-API/pkg/domain/login"
	"github.com/daniel1sender/Desafio-API/pkg/domain/transfers"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	useCase      transfers.UseCase
	loginUseCase login.UseCase
	logger       *logrus.Entry
}

func NewHandler(useCase transfers.UseCase, loginUseCase login.UseCase, logger *logrus.Entry) Handler {
	return Handler{
		useCase:      useCase,
		loginUseCase: loginUseCase,
		logger:       logger,
	}
}
