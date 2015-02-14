package networks

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type InterfacesRepository interface {
	GetTxBytes(interfaceName string) (int, error)
	GetRxBytes(interfaceName string) (int, error)
}

func NewInterfacesRepository() InterfacesRepository {
	repository := interfacesRepository{}
	return &repository
}

type interfacesRepository struct {
}

func (repository *interfacesRepository) GetTxBytes(interfaceName string) (int, error) {
	return readIntFromFile(repository.pathFor(interfaceName, "tx_bytes"))
}

func (repository *interfacesRepository) pathFor(interfaceName string, statName string) string {
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

func (repository *interfacesRepository) GetRxBytes(interfaceName string) (int, error) {
	return readIntFromFile(repository.pathFor(interfaceName, "rx_bytes"))
}
