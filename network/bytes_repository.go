package network

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type BytesRepository interface {
	GetTx() int
	GetRx() int
}

func NewBytesRepository(interfaceName string) BytesRepository {
	repository := bytesRepository{
		interfaceName: interfaceName,
	}
	return &repository
}

type bytesRepository struct {
	interfaceName string
}

func (repository *bytesRepository) GetTx() int {
	return readIntFromFile(repository.pathFor("tx_bytes"))
}

func (repository *bytesRepository) pathFor(statName string) string {
	return "/sys/class/net/" + repository.interfaceName + "/statistics/" + statName
}

func readIntFromFile(path string) int {
	var rawValue []byte
	var value int
	var err error

	rawValue, err = ioutil.ReadFile(path)

	if err != nil {
		return 0
	}

	value, err = strconv.Atoi(strings.Trim(string(rawValue), "\n"))

	if err != nil {
		return 0
	}

	return value
}

func (repository *bytesRepository) GetRx() int {
	return readIntFromFile(repository.pathFor("rx_bytes"))
}
