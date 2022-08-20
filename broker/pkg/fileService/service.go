package fileService

import (
	"bufio"
	"fmt"
	"graph-broker/broker/models"
	"log"
	"os"
	"strconv"
	"strings"
)

type Service interface {
	Save(file *os.File, dataChunc []*models.Data) error
	// reads from file and returns all data, and delete the file
	Get(file *os.File) ([]*models.Data, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) Save(file *os.File, dataChunc []*models.Data) error {
	chunc := ""
	for _, data := range dataChunc {
		chunc = chunc + fmt.Sprintf("%d %d %d %s\n", data.Round, data.Sender, data.Id, data.Message)
	}
	_, err := file.WriteString(chunc)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *service) Get(file *os.File) ([]*models.Data, error) {
	datas := make([]*models.Data, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		detailes := strings.Split(scanner.Text(), " ")
		data := models.Data{
			Round:   atoi(detailes[0]),
			Sender:  atoi(detailes[1]),
			Id:      atoi(detailes[2]),
			Message: detailes[3],
		}
		datas = append(datas, &data)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return datas, nil
}

func atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
