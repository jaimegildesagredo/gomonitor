package networks

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type InterfacesRepository interface {
	GetTxBytes(interfaceName string) int
	GetRxBytes(interfaceName string) int
	GetAllInterfaces() []string
	Exists(interfaceName string) bool
}

func NewInterfacesRepository() InterfacesRepository {
	repo := interfacesRepository{
		baseDir: "/sys/class/net",
	}
	return &repo
}

type interfacesRepository struct {
	baseDir string
}

func (repo *interfacesRepository) GetTxBytes(interfaceName string) int {
	return readIntFromFile(repo.pathFor(interfaceName, "tx_bytes"))
}

func (repo *interfacesRepository) pathFor(interfaceName string, statName string) string {
	return repo.baseDir + "/" + interfaceName + "/statistics/" + statName
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

func (repo *interfacesRepository) GetRxBytes(interfaceName string) int {
	return readIntFromFile(repo.pathFor(interfaceName, "rx_bytes"))
}

func (repo *interfacesRepository) GetAllInterfaces() []string {
	interfaces := []string{}

	contents, err := ioutil.ReadDir(repo.baseDir)

	if err != nil {
		log.Println("Error getting all network interfaces", err)
		return interfaces
	}

	for _, item := range contents {
		interfaces = append(interfaces, item.Name())
	}

	return interfaces
}

func (repo *interfacesRepository) Exists(interfaceName string) bool {
	for _, name := range repo.GetAllInterfaces() {
		if interfaceName == name {
			return true
		}
	}
	return false
}
