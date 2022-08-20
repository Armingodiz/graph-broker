package publisherService

import (
	"graph-broker/broker/models"
)

type Service interface {
	Send(dataChunc []*models.Data) (failedDatas []*models.Data, err error)
}

type service struct {
	ConsumerHost string
	ConsumerPort string
}

func NewService(consumerHost, consumerPort string) Service {
	return &service{
		ConsumerHost: consumerHost,
		ConsumerPort: consumerPort,
	}
}

func (s *service) Send(dataChunc []*models.Data) (failedDatas []*models.Data, err error) {
	return dataChunc, nil
}
