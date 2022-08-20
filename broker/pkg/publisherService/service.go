package publisherService

import (
	"encoding/json"
	"graph-broker/broker/models"
	"log"
	"net"
)

type Service interface {
	Send(dataChunc []*models.Data) (failedDatas []*models.Data)
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

func (s *service) Send(dataChunc []*models.Data) (failedDatas []*models.Data) {
	failedDatas = dataChunc
	tcpAddress, err := net.ResolveTCPAddr(s.ConsumerHost, s.ConsumerPort)
	if err != nil {
		return
	}
	connection, err := net.DialTCP(s.ConsumerHost, nil, tcpAddress)
	if err != nil {
		return
	}
	failedDatas = make([]*models.Data, 0)
	for _, data := range dataChunc {
		bytes, _ := json.Marshal(data)
		_, err := connection.Write([]byte(string(bytes) + "\n"))
		if err != nil {
			log.Println(err)
			failedDatas = append(failedDatas, data)
		}
	}
	return failedDatas
}
