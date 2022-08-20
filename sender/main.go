package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	network := os.Getenv("SENDER_NETWORK")
	address := os.Getenv("SENDER_ADDRESS")
	countWorkersEnv := os.Getenv("COUNT_WORKERS")
	intervalEnv := os.Getenv("INTERVAL")
	countWorkers, err := strconv.Atoi(countWorkersEnv)
	if err != nil {
		countWorkers = 10
	}
	interval, err := strconv.Atoi(intervalEnv)
	if err != nil {
		interval = 5
	}
	if network == "" {
		network = "tcp"
	}
	if address == "" {
		address = "localhost:80"
	}
	errChan := make(chan error, countWorkers)
	go func(dataRound int) {
		for {
			for i := 0; i < countWorkers; i++ {
				connection, err := connectToServer(network, address)
				if err != nil {
					errChan <- err
				}
				go worker(connection, dataRound, i, errChan)
			}
			time.Sleep(time.Minute * time.Duration(interval))
			dataRound++
		}
	}(0)
	for err := range errChan {
		panic(err)
	}
}
func connectToServer(network, address string) (connection *net.TCPConn, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	connection, err = net.DialTCP(network, nil, tcpAddress)
	return
}

func worker(conn *net.TCPConn, dataRound, id int, errChan chan error) {
	defer conn.Close()
	for i := 0; i < 1000; i++ {
		data := Data{
			Round:   dataRound,
			Id:      i,
			Sender:  id,
			Message: "test",
		}
		bytes, _ := json.Marshal(data)
		//fmt.Println("sending data:", string(bytes))
		_, err := conn.Write([]byte(string(bytes) + "\n"))
		if err != nil {
			log.Println(err)
			errChan <- err
			break
		}
	}
}

type Data struct {
	Round   int    `json:"round"`
	Sender  int    `json:"sender"`
	Id      int    `json:"id"`
	Message string `json:"message"`
}
