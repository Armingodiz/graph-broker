package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	consumerNetwork := os.Getenv("CONSUMER_NETWORK")
	consumerPort := os.Getenv("CONSUMER_PORT")
	if consumerNetwork == "" {
		consumerNetwork = "tcp"
	}
	if consumerPort == "" {
		consumerPort = ":81"
	}
	listener, err := createListener(consumerNetwork, consumerPort)
	if err != nil {
		log.Println(err)
		return
	}
	errChan := make(chan error, 1)
	go func() {
		for {
			connection, err := listenningForConnection(listener, consumerPort)
			if err != nil {
				log.Println(err)
				return
			}
			go handleConnection(connection, errChan)
		}
	}()
	for err := range errChan {
		panic(err)
	}
}

func handleConnection(conn net.Conn, errChan chan error) {
	dataChunc := make([]*Data, 0)
	for {
		buffer := make([]byte, 1638400)
		n, err := conn.Read(buffer[0:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				errChan <- err
			}

		}
		all := string(buffer[:n])
		for _, line := range strings.Split(all, "\n") {
			if line == "" {
				continue
			}
			//fmt.Println(line)
			var dt Data
			err = json.Unmarshal([]byte(line), &dt)
			if err != nil {
				errChan <- err
			}
			dataChunc = append(dataChunc, &dt)
		}
	}
	fmt.Println("received " + strconv.Itoa(len(dataChunc)) + " data")
}

type Data struct {
	Round   int    `json:"round"`
	Sender  int    `json:"sender"`
	Id      int    `json:"id"`
	Message string `json:"message"`
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
