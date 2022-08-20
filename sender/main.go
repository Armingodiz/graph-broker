package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	network := os.Getenv("SERVER_NETWORK")
	address := os.Getenv("SERVER_ADDRESS")
	countWorkersEnv := os.Getenv("COUNT_WORKERS")
	countWorkers, err := strconv.Atoi(countWorkersEnv)
	if err != nil {
		countWorkers = 10
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
				tcpAddress, err := net.ResolveTCPAddr(network, address)
				if err != nil {
					return
				}
				connection, err := net.DialTCP(network, nil, tcpAddress)
				if err != nil {
					return
				}
				go worker(connection, dataRound, i, errChan)
			}
			time.Sleep(time.Second * 5)
			dataRound++
		}
	}(0)
	for err := range errChan {
		panic(err)
	}
}

func worker(conn *net.TCPConn, dataRound, id int, errChan chan error) {
	for i := 0; i < 1000; i++ {
		data := Data{
			Round:   dataRound,
			Id:      i,
			Sender:  id,
			Message: "test",
		}
		bytes, _ := json.Marshal(data)
		fmt.Println("sending data:", string(bytes))
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
