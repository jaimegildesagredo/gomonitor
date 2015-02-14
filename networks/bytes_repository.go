package networks

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type BytesRepository interface {
	GetTx(string) (int, error)
	GetRx(string) (int, error)
}

func NewBytesRepository() BytesRepository {
	repository := bytesRepository{}
	return &repository
}

type bytesRepository struct {
}

func (repository *bytesRepository) GetTx(interfaceName string) (int, error) {
	return readIntFromFile(repository.pathFor(interfaceName, "tx_bytes"))
}

func (repository *bytesRepository) pathFor(interfaceName string, statName string) string {
	return "/sys/class/net/" + interfaceName + "/statistics/" + statName
}

func readIntFromFile(path string) (int, error) {
	var rawValue []byte
	var value int
	var err error

	rawValue, err = ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	value, err = strconv.Atoi(strings.Trim(string(rawValue), "\n"))

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (repository *bytesRepository) GetRx(interfaceName string) (int, error) {
	return readIntFromFile(repository.pathFor(interfaceName, "rx_bytes"))
}
