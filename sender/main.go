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
	workerPool := make([]*worker, countWorkers)
	errChan := make(chan error, countWorkers)
	for i := 0; i < countWorkers; i++ {
		worker, err := newWorker(network, address, i, errChan)
		if err != nil {
			panic(err)
		}
		workerPool[i] = worker
	}
	go func(dataRound int) {
		for {
			for _, worker := range workerPool {
				go worker.Send(dataRound)
			}
			time.Sleep(time.Minute * 5)
			dataRound++
		}
	}(0)
	for err := range errChan {
		panic(err)
	}
}

type worker struct {
	Conn    *net.TCPConn
	Id      int
	ErrChan chan error
}

func newWorker(network, address string, id int, errChan chan error) (w *worker, err error) {
	tcpAddress, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return
	}
	connection, err := net.DialTCP(network, nil, tcpAddress)
	if err != nil {
		return
	}
	w = &worker{
		Conn:    connection,
		Id:      id,
		ErrChan: errChan,
	}
	return
}

func (w *worker) Send(dataRound int) {
	for i := 0; i < 1000; i++ {
		data := Data{
			Round:   dataRound,
			Id:      i,
			Sender:  w.Id,
			Message: "test",
		}
		bytes, _ := json.Marshal(data)
		fmt.Println("sending data:", string(bytes))
		_, err := w.Conn.Write([]byte(string(bytes) + "\n"))
		if err != nil {
			log.Println(err)
			w.ErrChan <- err
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
