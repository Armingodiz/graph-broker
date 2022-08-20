package main

import (
	"encoding/json"
	"fmt"
	"graph-broker/broker/models"
	"graph-broker/broker/pkg/fileService"
	"graph-broker/broker/pkg/publisherService"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	network := os.Getenv("SERVER_NETWORK")
	port := os.Getenv("SERVER_PORT")
	if network == "" {
		network = "tcp"
	}
	if port == "" {
		port = ":80"
	}
	consumerNetwork := os.Getenv("CONSUMER_NETWORK")
	ConsumerAddress := os.Getenv("CONSUMER_ADDRESS")
	if consumerNetwork == "" {
		consumerNetwork = "tcp"
	}
	if ConsumerAddress == "" {
		ConsumerAddress = "localhost:81"
	}
	listener, err := createListener(network, port)
	if err != nil {
		log.Println(err)
		return
	}
	errChan := make(chan error, 1)
	retryChan := make(chan []*models.Data, 10)
	dataHandler := &DataHandler{
		FileService:      fileService.NewService(),
		PublisherService: publisherService.NewService(consumerNetwork, ConsumerAddress),
	}
	go dataHandler.handlePersistant(retryChan, errChan)
	go func() {
		for {
			connection, err := listenningForConnection(listener, port)
			if err != nil {
				log.Println(err)
				return
			}
			go dataHandler.handleConnection(connection, errChan, retryChan)
		}
	}()
	go dataHandler.handleRetriers(retryChan, errChan)
	for err := range errChan {
		panic(err)
	}
}

func createListener(network, address string) (listener *net.TCPListener, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	listener, err = net.ListenTCP(network, tcpAddress)
	return
}

func listenningForConnection(listener *net.TCPListener, port string) (conn net.Conn, err error) {
	fmt.Println("listenning on port " + port + " ...")
	conn, err = listener.Accept()
	if err != nil {
		return
	}
	return
}

type DataHandler struct {
	FileService      fileService.Service
	PublisherService publisherService.Service
}

func (dt *DataHandler) handleConnection(conn net.Conn, errChan chan error, retryChan chan []*models.Data) {
	defer conn.Close()
	dataChunc := make([]*models.Data, 0)
	for {
		done := false
		buffer := make([]byte, 1638400)
		n, err := conn.Read(buffer[0:])
		if err != nil {
			errChan <- err
		}
		all := string(buffer[:n])
		for _, line := range strings.Split(all, "\n") {
			if line == "" {
				continue
			}
			var dt models.Data
			err = json.Unmarshal([]byte(line), &dt)
			if err != nil {
				log.Println(line)
				//errChan <- err
			}
			dataChunc = append(dataChunc, &dt)
			if dt.Id == 999 {
				done = true
				break
			}
		}
		if done {
			break
		}
	}
	fmt.Println("received " + strconv.Itoa(len(dataChunc)) + " data")
	failedDatas := dt.PublisherService.Send(dataChunc)
	if len(failedDatas) > 0 {
		retryChan <- failedDatas
	}
}

func (dt *DataHandler) handleRetriers(retryChan chan []*models.Data, errCchan chan error) {
	for datas := range retryChan {
		f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			errCchan <- err
		}
		fmt.Println("TSSSSSsss")
		err = dt.FileService.Save(f, datas)
		if err != nil {
			errCchan <- err
		}
	}
}

func (dt *DataHandler) handlePersistant(retryChan chan []*models.Data, errChan chan error) {
	for {
		f, err := os.Open("data.txt")
		if err != nil {
			errChan <- err
		}
		datas, err := dt.FileService.Get(f)
		if err != nil {
			errChan <- err
		}
		if err := os.Truncate("data.txt", 0); err != nil {
			log.Printf("Failed to truncate: %v", err)
			errChan <- err
		}
		if len(datas) > 0 {
			failedDatas := dt.PublisherService.Send(datas)
			if len(failedDatas) > 0 {
				retryChan <- failedDatas
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
