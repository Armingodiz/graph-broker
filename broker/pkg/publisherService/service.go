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
	ConsumerAddress string
	ConsumerHost    string
}

func NewService(consumerHost, consumerAddress string) Service {
	return &service{
		ConsumerAddress: consumerAddress,
		ConsumerHost:    consumerHost,
	}
}

func (s *service) Send(dataChunc []*models.Data) (failedDatas []*models.Data) {
	failedDatas = dataChunc
	connection, _ := connectToServer(s.ConsumerHost, s.ConsumerAddress)
	defer connection.Close()
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

func connectToServer(network, address string) (connection *net.TCPConn, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	connection, err = net.DialTCP(network, nil, tcpAddress)
	return
}
